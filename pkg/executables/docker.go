package executables

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/eks-anywhere/pkg/logger"
)

const dockerPath = "docker"

type Docker struct {
	executable Executable
}

func NewDocker(executable Executable) *Docker {
	return &Docker{executable: executable}
}

func (d *Docker) GetDockerLBPort(ctx context.Context, clusterName string) (port string, err error) {
	clusterLBName := fmt.Sprintf("%s-lb", clusterName)
	if stdout, err := d.executable.Execute(ctx, "port", clusterLBName, "6443/tcp"); err != nil {
		return "", err
	} else {
		return strings.Split(stdout.String(), ":")[1], nil
	}
}

func (d *Docker) DockerPullImage(ctx context.Context, image string) error {
	logger.V(2).Info("Pulling docker image", "image", image)
	if _, err := d.executable.Execute(ctx, "pull", image); err != nil {
		return err
	} else {
		return nil
	}
}

func (d *Docker) SetUpCLITools(ctx context.Context, image string) error {
	logger.V(1).Info("Setting up cli docker dependencies")
	if err := d.DockerPullImage(ctx, image); err != nil {
		return err
	} else {
		return nil
	}
}

func (d *Docker) Version(ctx context.Context) (int, error) {
	cmdOutput, err := d.executable.Execute(ctx, "version", "--format", "{{.Client.Version}}")
	if err != nil {
		return 0, fmt.Errorf("please check if docker is installed and running %v", err)
	}
	dockerVersion := strings.TrimSpace(cmdOutput.String())
	versionSplit := strings.Split(dockerVersion, ".")
	installedMajorVersion := versionSplit[0]
	installedMajorVersionInt, err := strconv.Atoi(installedMajorVersion)
	if err != nil {
		return 0, err
	}
	return installedMajorVersionInt, nil
}

func (d *Docker) AllocatedMemory(ctx context.Context) (uint64, error) {
	cmdOutput, err := d.executable.Execute(ctx, "info", "--format", "'{{json .MemTotal}}'")
	if err != nil {
		return 0, fmt.Errorf("please check if docker is installed and running %v", err)
	}
	totalMemory := cmdOutput.String()
	totalMemory = totalMemory[1 : len(totalMemory)-2]
	return strconv.ParseUint(totalMemory, 10, 64)
}