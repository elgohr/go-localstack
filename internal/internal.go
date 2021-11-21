package internal

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
	"time"

	"github.com/Masterminds/semver/v3"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . DockerClient
// DockerClient represents a way to interact with the Docker deamon
type DockerClient interface {
	ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error)
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
	ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error
}

// MustParseConstraint panics if a semver constraint is invalid
func MustParseConstraint(constraint string) *semver.Constraints {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		panic(fmt.Errorf("localstack: invalid version constraint for port change: %w", err))
	}
	return c
}
