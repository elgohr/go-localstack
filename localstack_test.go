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

package localstack_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/elgohr/go-localstack"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	if err := clean(); err != nil {
		log.Fatalln(err)
	}
	run := m.Run()
	if err := clean(); err != nil {
		log.Fatalln(err)
	}
	os.Exit(run)
}

func TestWithLogger(t *testing.T) {
	for _, s := range []struct {
		name   string
		level  log.Level
		expect func(t require.TestingT, object interface{}, msgAndArgs ...interface{})
	}{
		{
			name:   "with debug log level",
			level:  log.DebugLevel,
			expect: require.NotEmpty,
		},
		{
			name:   "with fatal log level",
			level:  log.FatalLevel,
			expect: require.Empty,
		},
	} {
		t.Run(s.name, func(t *testing.T) {
			buf := &concurrentWriter{buf: &bytes.Buffer{}}
			logger := log.New()
			logger.SetLevel(s.level)
			logger.SetOutput(buf)
			l, err := localstack.NewInstance(localstack.WithLogger(logger))
			require.NoError(t, err)
			require.NoError(t, l.Start())
			require.NoError(t, l.Stop())
			s.expect(t, buf.Bytes())
		})
	}
}

func TestWithTimeoutOnStartup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l, err := localstack.NewInstance(localstack.WithTimeout(time.Second))
	require.NoError(t, err)
	require.EqualError(t, l.StartWithContext(ctx), "localstack container has been stopped")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err)
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	require.NoError(t, err)
	for _, c := range containers {
		if c.Image == "go-localstack" {
			t.Fatal("image is still running but should be terminated")
		}
	}
}

func TestWithTimeoutAfterStartup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l, err := localstack.NewInstance(localstack.WithTimeout(20 * time.Second))
	require.NoError(t, err)
	timer := time.NewTimer(25 * time.Second)
	defer timer.Stop()
	require.NoError(t, l.StartWithContext(ctx))

	<-timer.C
	cli, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err)
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	require.NoError(t, err)
	for _, c := range containers {
		if c.Image == "go-localstack" {
			t.Fatal("image is still running but should be terminated")
		}
	}
}

func TestWithLabels(t *testing.T) {
	for _, s := range []struct {
		name   string
		labels map[string]string
	}{
		{
			name: "with multiple labels",
			labels: map[string]string{
				"label1": "aaa111",
				"label2": "bbb222",
				"label3": "ccc333",
			},
		},
		{
			name: "with nil label map",
		},
	} {
		t.Run(s.name, func(t *testing.T) {
			l, err := localstack.NewInstance(localstack.WithLabels(s.labels))
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			require.NoError(t, l.StartWithContext(ctx))
			defer func() { require.NoError(t, l.Stop()) }()

			cli, err := client.NewClientWithOpts()
			require.NoError(t, err)

			containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true, Quiet: true})
			require.NoError(t, err)

			require.True(t, atLeastOneContainerMatchesLabels(s.labels, containers))
		})
	}
}

func TestLocalStack(t *testing.T) {
	for _, s := range []struct {
		name   string
		input  []localstack.InstanceOption
		expect func(t *testing.T, l *localstack.Instance)
	}{
		{
			name:   "with version before breaking change",
			input:  []localstack.InstanceOption{localstack.WithVersion("0.11.4")},
			expect: havingIndividualEndpoints,
		},
		{
			name:   "with nil",
			input:  nil,
			expect: havingOneEndpoint,
		},
		{
			name:   "with empty",
			input:  []localstack.InstanceOption{},
			expect: havingOneEndpoint,
		},
		{
			name:   "with breaking change version",
			input:  []localstack.InstanceOption{localstack.WithVersion("0.11.5")},
			expect: havingOneEndpoint,
		},
		{
			name:   "with version after breaking change",
			input:  []localstack.InstanceOption{localstack.WithVersion("latest")},
			expect: havingOneEndpoint,
		},
	} {
		t.Run(s.name, func(t *testing.T) {
			l, err := localstack.NewInstance(s.input...)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, l.Stop())
			}()
			require.NoError(t, l.Start())
			s.expect(t, l)
		})
	}
}

func TestLocalStackWithContext(t *testing.T) {
	for _, s := range []struct {
		name   string
		input  []localstack.InstanceOption
		expect func(t *testing.T, l *localstack.Instance)
	}{
		{
			name:   "with version before breaking change",
			input:  []localstack.InstanceOption{localstack.WithVersion("0.11.4")},
			expect: havingIndividualEndpoints,
		},
		{
			name:   "with nil",
			input:  nil,
			expect: havingOneEndpoint,
		},
		{
			name:   "with empty",
			input:  []localstack.InstanceOption{},
			expect: havingOneEndpoint,
		},
		{
			name:   "with breaking change version",
			input:  []localstack.InstanceOption{localstack.WithVersion("0.11.5")},
			expect: havingOneEndpoint,
		},
		{
			name:   "with version after breaking change",
			input:  []localstack.InstanceOption{localstack.WithVersion("latest")},
			expect: havingOneEndpoint,
		},
	} {
		t.Run(s.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			l, err := localstack.NewInstance(s.input...)
			require.NoError(t, err)
			require.NoError(t, l.StartWithContext(ctx))
			s.expect(t, l)
		})
	}
}

func TestLocalStackWithIndividualServicesOnContext(t *testing.T) {
	cl := http.Client{Timeout: time.Second}
	for service := range localstack.AvailableServices {
		t.Run(service.Name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			l, err := localstack.NewInstance([]localstack.InstanceOption{localstack.WithVersion("0.11.4")}...)
			require.NoError(t, err)
			require.NoError(t, l.StartWithContext(ctx, service))
			for testService := range localstack.AvailableServices {
				conn, err := net.DialTimeout("tcp", strings.TrimPrefix(l.EndpointV2(testService), "http://"), time.Second)
				if testService == service || testService == localstack.DynamoDB {
					require.NoError(t, err, testService)
					require.NoError(t, conn.Close())
				} else if testService != localstack.FixedPort {
					require.Error(t, err, testService)
				}
			}
			cancel()
			require.Eventually(t, func() bool {
				_, err := cl.Get(l.EndpointV2(service))
				return err != nil
			}, time.Minute, time.Second)
		})
	}
}

func TestInstanceStartedTwiceWithoutLeaking(t *testing.T) {
	l, err := localstack.NewInstance()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, l.Stop())
	}()
	require.NoError(t, l.Start())
	firstInstance := l.Endpoint(localstack.S3)
	require.NoError(t, l.Start())
	_, err = net.Dial("tcp", firstInstance)
	require.Error(t, err, "should be teared down")
}

func TestContextInstanceStartedTwiceWithoutLeaking(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	l, err := localstack.NewInstance()
	require.NoError(t, err)
	require.NoError(t, l.Start())
	firstInstance := l.Endpoint(localstack.S3)
	require.NoError(t, l.StartWithContext(ctx))
	_, err = net.Dial("tcp", firstInstance)
	require.Error(t, err, "should be teared down")
}

func TestInstanceWithVersions(t *testing.T) {
	for _, s := range []struct {
		version string
		expect  func(t require.TestingT, err error, msgAndArgs ...interface{})
	}{
		{version: "0.11.5", expect: require.NoError},
		{version: "0.11.3", expect: require.NoError},
		{version: "latest", expect: require.NoError},
		{version: "bad.version.34", expect: require.Error},
	} {
		t.Run(s.version, func(t *testing.T) {
			_, err := localstack.NewInstance(localstack.WithVersion(s.version))
			s.expect(t, err)
		})
	}
}

func TestInstanceWithBadDockerEnvironment(t *testing.T) {
	urlIfSet := os.Getenv("DOCKER_URL")
	defer func() {
		require.NoError(t, os.Setenv("DOCKER_URL", urlIfSet))
	}()

	require.NoError(t, os.Setenv("DOCKER_URL", "what-is-this-thing:///var/run/not-a-valid-docker.sock"))

	_, err := localstack.NewInstance()
	require.NoError(t, err)
}

func TestInstanceStopWithoutStarted(t *testing.T) {
	l, err := localstack.NewInstance()
	require.NoError(t, err)
	require.NoError(t, l.Stop())
}

func TestInstanceEndpointWithoutStarted(t *testing.T) {
	l, err := localstack.NewInstance()
	require.NoError(t, err)
	require.Empty(t, l.Endpoint(localstack.S3))
}

func TestWithClientFromEnv(t *testing.T) {
	for _, s := range []struct {
		name        string
		given       func(t *testing.T)
		expectOpt   func(t require.TestingT, opt localstack.InstanceOption, err error)
		expectStart func(t require.TestingT, err error)
	}{
		{
			name: "is ok with client from env",
			given: func(t *testing.T) {
				require.NoError(t, os.Setenv("DOCKER_API_VERSION", "0"))
			},
			expectOpt: func(t require.TestingT, opt localstack.InstanceOption, err error) {
				require.NoError(t, err)
				require.NotNil(t, opt)
			},
			expectStart: func(t require.TestingT, err error) {
				require.Error(t, err)
				require.True(t, strings.HasPrefix(err.Error(), "localstack: could not build image: Error response from daemon: client version 0 is too old."), err)
			},
		},
		{
			name: "publishes errors",
			given: func(t *testing.T) {
				require.NoError(t, os.Setenv("DOCKER_HOST", "localhost"))
			},
			expectOpt: func(t require.TestingT, opt localstack.InstanceOption, err error) {
				require.EqualError(t, err, "localstack: could not connect to docker: unable to parse docker host `localhost`")
				require.Nil(t, opt)
			},
		},
	} {
		t.Run(s.name, func(t *testing.T) {
			defer func() {
				require.NoError(t, os.Unsetenv("DOCKER_HOST"))
				require.NoError(t, os.Unsetenv("DOCKER_API_VERSION"))
			}()
			s.given(t)
			opt, err := localstack.WithClientFromEnv()
			s.expectOpt(t, opt, err)
			if s.expectStart != nil {
				i, err := localstack.NewInstance(opt)
				require.NoError(t, err)

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				s.expectStart(t, i.StartWithContext(ctx))
			}
		})
	}
}

func havingOneEndpoint(t *testing.T, l *localstack.Instance) {
	endpoints := map[string]struct{}{}
	for service := range localstack.AvailableServices {
		endpoints[l.Endpoint(service)] = struct{}{}
	}
	require.Equal(t, 1, len(endpoints), endpoints)
}

func havingIndividualEndpoints(t *testing.T, l *localstack.Instance) {
	endpoints := map[string]struct{}{}
	for service := range localstack.AvailableServices {
		endpoint := l.Endpoint(service)
		checkAddress(t, endpoint)

		_, exists := endpoints[endpoint]
		require.False(t, exists, fmt.Sprintf("%s duplicated in %v", endpoint, endpoints))

		endpoints[endpoint] = struct{}{}
	}
	require.Equal(t, len(localstack.AvailableServices), len(endpoints))
}

func checkAddress(t *testing.T, val string) {
	require.True(t, strings.HasPrefix(val, "localhost:"), val)
	require.NotEmpty(t, val[10:])
}

func atLeastOneContainerMatchesLabels(labels map[string]string, containers []types.Container) bool {
	for _, container := range containers {
		if matchesLabels(labels, container) {
			return true
		}
	}
	return false
}

func matchesLabels(labels map[string]string, container types.Container) bool {
	for k, v := range labels {
		val, exists := container.Labels[k]
		if !exists || v != val {
			return false
		}
	}
	return true
}

func clean() error {
	timeout := time.Second
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if list, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true}); err == nil {
		for _, l := range list {
			if err := cli.ContainerStop(ctx, l.ID, &timeout); err != nil {
				log.Println(err)
			}
		}
	} else {
		return err
	}
	if _, err := cli.ContainersPrune(ctx, filters.Args{}); err != nil {
		return err
	}
	return nil
}

type concurrentWriter struct {
	buf *bytes.Buffer
	mu  sync.RWMutex
}

func (c *concurrentWriter) Write(p []byte) (n int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buf.Write(p)
}

func (c *concurrentWriter) Bytes() []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.buf.Bytes()
}
