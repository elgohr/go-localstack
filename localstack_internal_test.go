// Copyright 2021 - Lars Gohr
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package localstack

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/elgohr/go-localstack/internal/internalfakes"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestInstance_Start_Fails(t *testing.T) {
	for _, tt := range [...]struct {
		when  string
		given func(f *internalfakes.FakeDockerClient) *Instance
		then  func(t *testing.T, err error, f *internalfakes.FakeDockerClient)
	}{
		{
			when: "can't restart localstack when already running",
			given: func(f *internalfakes.FakeDockerClient) *Instance {
				f.ContainerStopReturns(errors.New("can't stop"))
				return &Instance{
					cli:         f,
					log:         logrus.StandardLogger(),
					containerId: "running",
				}
			},
			then: func(t *testing.T, err error, f *internalfakes.FakeDockerClient) {
				require.EqualError(t, err, "localstack: can't stop an already running instance: can't stop")
				require.Equal(t, 0, f.ImageBuildCallCount())
				require.Equal(t, 0, f.ContainerCreateCallCount())
				require.Equal(t, 0, f.ContainerStartCallCount())
				require.Equal(t, 0, f.ContainerInspectCallCount())
			},
		},
		{
			when: "can't build image",
			given: func(f *internalfakes.FakeDockerClient) *Instance {
				f.ImageBuildReturns(types.ImageBuildResponse{}, errors.New("can't build"))
				return &Instance{
					cli: f,
					log: logrus.StandardLogger(),
				}
			},
			then: func(t *testing.T, err error, f *internalfakes.FakeDockerClient) {
				require.EqualError(t, err, "localstack: could not build image: can't build")
				require.Equal(t, 1, f.ImageBuildCallCount())
				require.Equal(t, 0, f.ContainerCreateCallCount())
				require.Equal(t, 0, f.ContainerStartCallCount())
				require.Equal(t, 0, f.ContainerInspectCallCount())
			},
		},
		{
			when: "can't close after building image",
			given: func(f *internalfakes.FakeDockerClient) *Instance {
				f.ImageBuildReturns(types.ImageBuildResponse{Body: ErrCloser(strings.NewReader(""), errors.New("can't close"))}, nil)
				f.ContainerCreateReturns(container.ContainerCreateCreatedBody{}, errors.New("can't create"))
				return &Instance{
					cli: f,
					log: logrus.StandardLogger(),
				}
			},
			then: func(t *testing.T, err error, f *internalfakes.FakeDockerClient) {
				require.EqualError(t, err, "localstack: could not create container: can't create")
			},
		},
		{
			when: "can't create container",
			given: func(f *internalfakes.FakeDockerClient) *Instance {
				f.ImageBuildReturns(types.ImageBuildResponse{Body: io.NopCloser(strings.NewReader(""))}, nil)
				f.ContainerCreateReturns(container.ContainerCreateCreatedBody{}, errors.New("can't create"))
				return &Instance{
					cli: f,
					log: logrus.StandardLogger(),
				}
			},
			then: func(t *testing.T, err error, f *internalfakes.FakeDockerClient) {
				require.EqualError(t, err, "localstack: could not create container: can't create")
				require.Equal(t, 1, f.ImageBuildCallCount())
				require.Equal(t, 1, f.ContainerCreateCallCount())
				ctx, config, hostConfig, networkingConfig, platform, containerName := f.ContainerCreateArgsForCall(0)
				require.NotNil(t, ctx)
				require.Equal(t, &container.Config{
					Image:        "go-localstack",
					Env:          []string{},
					Tty:          true,
					AttachStdout: true,
					AttachStderr: true,
				}, config)
				pm := nat.PortMap{}
				for service := range AvailableServices {
					pm[nat.Port(service.Port)] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: ""}}
				}
				require.Equal(t, &container.HostConfig{
					PortBindings: pm,
					AutoRemove:   true,
				}, hostConfig)
				require.Nil(t, networkingConfig)
				require.Nil(t, platform)
				require.Empty(t, containerName)
				require.Equal(t, 0, f.ContainerStartCallCount())
				require.Equal(t, 0, f.ContainerInspectCallCount())
			},
		},
		{
			when: "can't start container",
			given: func(f *internalfakes.FakeDockerClient) *Instance {
				f.ImageBuildReturns(types.ImageBuildResponse{Body: io.NopCloser(strings.NewReader(""))}, nil)
				f.ContainerStartReturns(errors.New("can't start"))
				return &Instance{
					cli: f,
					log: logrus.StandardLogger(),
				}
			},
			then: func(t *testing.T, err error, f *internalfakes.FakeDockerClient) {
				require.EqualError(t, err, "localstack: could not start container: can't start")
				require.Equal(t, 1, f.ImageBuildCallCount())
				require.Equal(t, 1, f.ContainerCreateCallCount())
				require.Equal(t, 1, f.ContainerStartCallCount())
				require.Equal(t, 0, f.ContainerInspectCallCount())
			},
		},
		{
			when: "container inspect doesn't contain ports",
			given: func(f *internalfakes.FakeDockerClient) *Instance {
				f.ImageBuildReturns(types.ImageBuildResponse{Body: io.NopCloser(strings.NewReader(""))}, nil)
				f.ContainerInspectReturns(types.ContainerJSON{NetworkSettings: &types.NetworkSettings{
					NetworkSettingsBase: types.NetworkSettingsBase{
						Ports: map[nat.Port][]nat.PortBinding{},
					},
				}}, nil)
				return &Instance{
					cli:       f,
					fixedPort: true,
					log:       logrus.StandardLogger(),
				}
			},
			then: func(t *testing.T, err error, f *internalfakes.FakeDockerClient) {
				require.EqualError(t, err, "localstack: could not get port from container")
				require.Equal(t, 1, f.ImageBuildCallCount())
				require.Equal(t, 1, f.ContainerCreateCallCount())
				require.Equal(t, 1, f.ContainerStartCallCount())
				require.Equal(t, 6, f.ContainerInspectCallCount())
			},
		},
	} {
		t.Run(tt.when, func(t *testing.T) {
			f := &internalfakes.FakeDockerClient{}
			f.ContainerLogsReturns(io.NopCloser(strings.NewReader("")), nil)
			tt.then(t, tt.given(f).Start(), f)
		})
	}
}

func TestInstance_StartWithContext_Fails_Stop_AfterTest(t *testing.T) {
	f := &internalfakes.FakeDockerClient{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	f.ContainerStopReturns(errors.New("can't stop"))
	i := &Instance{cli: f, log: logrus.StandardLogger(), containerId: "something"}
	require.EqualError(t, i.StartWithContext(ctx), "localstack: can't stop an already running instance: can't stop")
}

func TestInstance_Stop_Fails(t *testing.T) {
	f := &internalfakes.FakeDockerClient{}
	f.ContainerStopReturns(errors.New("can't stop"))
	i := &Instance{cli: f, log: logrus.StandardLogger(), containerId: "something"}
	require.EqualError(t, i.Stop(), "can't stop")
}

func TestInstance_checkAvailable_Session_Fails(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	require.NoError(t, os.Setenv("AWS_STS_REGIONAL_ENDPOINTS", "FAILURE"))
	defer func() {
		require.NoError(t, os.Unsetenv("AWS_STS_REGIONAL_ENDPOINTS"))
	}()
	i := &Instance{log: logrus.StandardLogger()}
	require.Error(t, i.checkAvailable(ctx))
}

func TestInstance_waitToBeAvailable_Context_Expired(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	i := &Instance{log: logrus.StandardLogger()}
	require.Error(t, i.waitToBeAvailable(ctx))
}

func ErrCloser(r io.Reader, err error) io.ReadCloser {
	return errCloser{Reader: r, Error: err}
}

type errCloser struct {
	io.Reader
	Error error
}

func (e errCloser) Close() error {
	return e.Error
}
