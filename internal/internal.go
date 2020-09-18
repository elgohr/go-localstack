package internal

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Pool

// Pool represents a dockertest.Pool
type Pool interface {
	RunWithOptions(opts *dockertest.RunOptions, hcOpts ...func(*docker.HostConfig)) (*dockertest.Resource, error)
	Purge(r *dockertest.Resource) error
	Retry(op func() error) error
}

// MustParseConstraint panics if a semver constraint is invalid
func MustParseConstraint(constraint string) *semver.Constraints {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		panic(fmt.Errorf("localstack: invalid version constraint for port change: %w", err))
	}
	return c
}
