package executables

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/eks-anywhere/pkg/logger"
)

const redactMask = "*****"

var redactedEnvKeys = []string{vSphereUsernameKey, vSpherePasswordKey}

type executable struct {
	cli string
}

type linuxDockerExecutable struct {
	cli      string
	image    string
	mountDir string
}

type Executable interface {
	Execute(ctx context.Context, args ...string) (stdout bytes.Buffer, err error)
	ExecuteWithEnv(ctx context.Context, envs map[string]string, args ...string) (stdout bytes.Buffer, err error)
	ExecuteWithStdin(ctx context.Context, in []byte, args ...string) (stdout bytes.Buffer, err error)
}

// this should only be called through the executables.builder
func NewExecutable(cli string) Executable {
	return &executable{
		cli: cli,
	}
}

// This currently returns a linuxDockerExecutable, but if we support other types of docker executables we can change
// the name of this constructor
func NewDockerExecutable(cli, image string, mountDir string) Executable {
	return &linuxDockerExecutable{
		cli:      cli,
		image:    image,
		mountDir: mountDir,
	}
}

func (e *linuxDockerExecutable) workingDirectory() (string, error) {
	path, err := filepath.Abs(e.mountDir)
	if err != nil {
		return "", fmt.Errorf("error getting abs path for working dir: %v", err)
	}

	return path, nil
}

func (e *linuxDockerExecutable) Execute(ctx context.Context, args ...string) (bytes.Buffer, error) {
	var stdout bytes.Buffer
	if command, err := e.buildCommand(map[string]string{}, e.cli, args...); err != nil {
		return stdout, err
	} else {
		return execute(ctx, "docker", nil, command...)
	}
}

func (e *linuxDockerExecutable) ExecuteWithStdin(ctx context.Context, in []byte, args ...string) (bytes.Buffer, error) {
	var stdout bytes.Buffer
	if command, err := e.buildCommand(map[string]string{}, e.cli, args...); err != nil {
		return stdout, err
	} else {
		return execute(ctx, "docker", in, command...)
	}
}

func (e *linuxDockerExecutable) buildCommand(envs map[string]string, cli string, args ...string) ([]string, error) {
	directory, err := e.workingDirectory()
	if err != nil {
		return nil, err
	}

	var envVars []string
	for k, v := range envs {
		envVars = append(envVars, "-e", fmt.Sprintf("%s=%s", k, v))
	}
	dockerCommands := []string{
		"run", "-i", "--network", "host", "-v", fmt.Sprintf("%s:%s", directory, directory), "-w",
		directory, "-v", "/var/run/docker.sock:/var/run/docker.sock", "-v", "/usr/bin/docker:/usr/bin/docker",
	}
	dockerCommands = append(dockerCommands, envVars...)
	dockerCommands = append(dockerCommands, "--entrypoint", cli, e.image)
	dockerCommands = append(dockerCommands, args...)
	return dockerCommands, nil
}

func (e *linuxDockerExecutable) ExecuteWithEnv(ctx context.Context, envs map[string]string, args ...string) (bytes.Buffer, error) {
	var stdout bytes.Buffer
	if command, err := e.buildCommand(envs, e.cli, args...); err != nil {
		return stdout, err
	} else {
		return execute(ctx, "docker", nil, command...)
	}
}

func (e *executable) Execute(ctx context.Context, args ...string) (bytes.Buffer, error) {
	return execute(ctx, e.cli, nil, args...)
}

func (e *executable) ExecuteWithStdin(ctx context.Context, in []byte, args ...string) (bytes.Buffer, error) {
	return execute(ctx, e.cli, in, args...)
}

func (e *executable) ExecuteWithEnv(ctx context.Context, envs map[string]string, args ...string) (stdout bytes.Buffer, err error) {
	for k, v := range envs {
		os.Setenv(k, v)
	}
	return e.Execute(ctx, args...)
}

func redactCreds(cmd string) string {
	redactedEnvs := []string{}
	for _, redactedEnvKey := range redactedEnvKeys {
		if env, found := os.LookupEnv(redactedEnvKey); found {
			redactedEnvs = append(redactedEnvs, env)
		}
	}

	for _, redactedEnv := range redactedEnvs {
		cmd = strings.ReplaceAll(cmd, redactedEnv, redactMask)
	}
	return cmd
}

func execute(ctx context.Context, cli string, in []byte, args ...string) (bytes.Buffer, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, cli, args...)
	logger.V(6).Info("Executing command", "cmd", redactCreds(cmd.String()))
	cmd.Stdout = &stdout
	if logger.MaxLogging() {
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = &stderr
	}
	if len(in) != 0 {
		cmd.Stdin = bytes.NewReader(in)
	}

	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return stdout, errors.New(stderr.String())
		} else {
			if !logger.MaxLogging() {
				logger.V(8).Info(cli, "stdout", stdout.String())
				logger.V(8).Info(cli, "stderr", stderr.String())
			}
			return stdout, errors.New(fmt.Sprint(err))
		}
	}
	if !logger.MaxLogging() {
		logger.V(8).Info(cli, "stdout", stdout.String())
		logger.V(8).Info(cli, "stderr", stderr.String())
	}
	return stdout, nil
}