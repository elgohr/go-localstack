package internal

import (
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Pool
type Pool interface {
	RunWithOptions(opts *dockertest.RunOptions, hcOpts ...func(*docker.HostConfig)) (*dockertest.Resource, error)
	Purge(r *dockertest.Resource) error
	Retry(op func() error) error
}
