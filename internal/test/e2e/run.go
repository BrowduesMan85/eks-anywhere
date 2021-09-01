package e2e

import (
	"fmt"
	"strings"
	"sync"

	"github.com/aws/eks-anywhere/internal/pkg/ec2"
	"github.com/aws/eks-anywhere/internal/pkg/ssm"
	"github.com/aws/eks-anywhere/pkg/logger"
)

type ParallelRunConf struct {
	MaxIntances         int
	AmiId               string
	InstanceProfileName string
	StorageBucket       string
	JobId               string
	SubnetId            string
	Regex               string
}

type instanceTestsResults struct {
	conf instanceRunConf
	err  error
}

func RunTestsInParallel(conf ParallelRunConf) error {
	testsList, err := listTests(conf.Regex)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	resultCh := make(chan instanceTestsResults)

	instancesConf := splitTests(testsList, conf)
	logTestGroups(instancesConf)
	for _, instanceConf := range instancesConf {
		go func(c instanceRunConf) {
			defer wg.Done()
			r := instanceTestsResults{conf: c}

			if err := RunTests(c); err != nil {
				r.err = err
			}

			resultCh <- r
		}(instanceConf)
		wg.Add(1)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	failedInstances := 0
	for r := range resultCh {
		if r.err != nil {
			logger.Error(r.err, "An e2e instance run has failed", "jobId", r.conf.jobId, "tests", r.conf.regex)
			failedInstances += 1
		} else {
			logger.Info("Ec2 instance tests completed succesfully", "jobId", r.conf.jobId, "tests", r.conf.regex)
		}
	}

	if failedInstances > 0 {
		return fmt.Errorf("%d/%d e2e instances failed", failedInstances, len(instancesConf))
	}

	return nil
}

type instanceRunConf struct {
	amiId, instanceProfileName, storageBucket, jobId, subnetId, regex string
}

func RunTests(conf instanceRunConf) error {
	session, err := newSession(conf.amiId, conf.instanceProfileName, conf.storageBucket, conf.jobId, conf.subnetId)
	if err != nil {
		return err
	}

	err = session.setup(conf.regex)
	if err != nil {
		return err
	}

	err = session.runTests(conf.regex)
	if err != nil {
		return err
	}

	return nil
}

func (e *E2ESession) runTests(regex string) error {
	logger.V(1).Info("Running e2e tests", "regex", regex)
	command := "./bin/e2e.test -test.v"
	if regex != "" {
		command = fmt.Sprintf("%s -test.run %s", command, regex)
	}

	command = e.commandWithEnvVars(command)

	err := ssm.Run(
		e.session,
		e.instanceId,
		command,
	)
	if err != nil {
		return fmt.Errorf("error running e2e tests: %v", err)
	}

	key := "Integration-Test-Done"
	value := "TRUE"
	err = ec2.TagInstance(e.session, e.instanceId, key, value)
	if err != nil {
		return fmt.Errorf("error tagging instance for e2e success: %v", err)
	}

	return nil
}

func (e *E2ESession) commandWithEnvVars(command string) string {
	fullCommand := make([]string, 0, len(e.testEnvVars)+1)

	for k, v := range e.testEnvVars {
		fullCommand = append(fullCommand, fmt.Sprintf("export %s=%s", k, v))
	}
	fullCommand = append(fullCommand, command)

	return strings.Join(fullCommand, "; ")
}

func splitTests(testsList []string, conf ParallelRunConf) []instanceRunConf {
	testPerInstance := len(testsList) / conf.MaxIntances
	if testPerInstance == 0 {
		testPerInstance = 1
	}

	runConfs := make([]instanceRunConf, 0, conf.MaxIntances)

	testsInCurrentInstance := make([]string, 0, testPerInstance)
	for _, testName := range testsList {
		testsInCurrentInstance = append(testsInCurrentInstance, testName)
		if len(testsInCurrentInstance) == testPerInstance {
			runConfs = append(runConfs, instanceRunConf{
				amiId:               conf.AmiId,
				instanceProfileName: conf.InstanceProfileName,
				storageBucket:       conf.StorageBucket,
				jobId:               fmt.Sprintf("%s-%d", conf.JobId, len(runConfs)),
				subnetId:            conf.SubnetId,
				regex:               strings.Join(testsInCurrentInstance, "|"),
			})

			testsInCurrentInstance = make([]string, 0, testPerInstance)
		}
	}

	return runConfs
}

func logTestGroups(instancesConf []instanceRunConf) {
	testGroups := make([]string, 0, len(instancesConf))
	for _, i := range instancesConf {
		testGroups = append(testGroups, i.regex)
	}
	logger.V(1).Info("Running tests in parallel", "testsGroups", testGroups)
}