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
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	for _, scenario := range []struct {
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
		s := scenario
		t.Run(s.name, func(t *testing.T) {
			t.Parallel()
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
	ctx := t.Context()
	l, err := localstack.NewInstance(localstack.WithTimeout(time.Second))
	require.NoError(t, err)
	require.EqualError(t, l.StartWithContext(ctx), "localstack container has been stopped")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err)
	cli.NegotiateAPIVersion(ctx)

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	require.NoError(t, err)
	for _, c := range containers {
		if strings.Contains(c.Image, "go-localstack") {
			t.Fatalf("%s is still running but should be terminated", c.Image)
		}
	}
}

func TestWithTimeoutAfterStartup(t *testing.T) {
	ctx := t.Context()
	l, err := localstack.NewInstance(localstack.WithTimeout(20 * time.Second))
	require.NoError(t, err)

	require.NoError(t, l.StartWithContext(ctx))

	cli, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err)
	cli.NegotiateAPIVersion(ctx)
	require.Eventually(t, func() bool {
		containers, err := cli.ContainerList(ctx, container.ListOptions{})
		require.NoError(t, err)
		for _, c := range containers {
			if strings.Contains(c.Image, "go-localstack") {
				return false
			}
		}
		return true
	}, 5*time.Minute, 200*time.Millisecond, "image is still running but should be terminated")
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

			ctx := t.Context()

			require.NoError(t, l.StartWithContext(ctx))
			t.Cleanup(func() { require.NoError(t, l.Stop()) })

			cli, err := client.NewClientWithOpts()
			require.NoError(t, err)
			cli.NegotiateAPIVersion(ctx)

			containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
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
			t.Cleanup(func() {
				assert.NoError(t, l.Stop())
			})
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
			ctx := t.Context()
			l, err := localstack.NewInstance(s.input...)
			require.NoError(t, err)
			require.NoError(t, l.StartWithContext(ctx))
			s.expect(t, l)
		})
	}
}

func TestLocalStackWithIndividualServicesOnContext(t *testing.T) {
	cl := &http.Client{Timeout: time.Second}
	dialer := &net.Dialer{Timeout: time.Second}
	for service := range localstack.AvailableServices {
		t.Run(service.Name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			l, err := localstack.NewInstance()
			require.NoError(t, err)
			require.NoError(t, l.StartWithContext(ctx, service))
			for testService := range localstack.AvailableServices {
				conn, err := dialer.DialContext(t.Context(), "tcp", strings.TrimPrefix(l.EndpointV2(testService), "http://"))
				if testService == service || testService == localstack.DynamoDB {
					require.NoError(t, err, testService)
					require.NoError(t, conn.Close())
				}
			}
			cancel()

			// wait until service was shutdown
			require.Eventually(t, func() bool {
				address := l.EndpointV2(service)
				if address == "" {
					return true
				}
				req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, address, nil)
				res, err := cl.Do(req)
				defer func() {
					if res == nil || res.Body == nil {
						return
					}
					_ = res.Body.Close()
				}()
				return err != nil
			}, 5*time.Minute, 300*time.Millisecond)
		})
	}
}

func TestLocalStackWithIndividualServices(t *testing.T) {
	cl := &http.Client{Timeout: time.Second}
	dialer := &net.Dialer{Timeout: time.Second}
	for service := range localstack.AvailableServices {
		t.Run(service.Name, func(t *testing.T) {
			l, err := localstack.NewInstance()
			require.NoError(t, err)
			require.NoError(t, l.Start(service))
			for testService := range localstack.AvailableServices {
				conn, err := dialer.DialContext(t.Context(), "tcp", strings.TrimPrefix(l.EndpointV2(testService), "http://"))
				if testService == service || testService == localstack.DynamoDB {
					require.NoError(t, err, testService)
					require.NoError(t, conn.Close())
				}
			}
			assert.NoError(t, l.Stop())

			// wait until service was shutdown
			require.Eventually(t, func() bool {
				address := l.EndpointV2(service)
				if address == "" {
					return true
				}
				req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, address, nil)
				res, err := cl.Do(req)
				defer func() {
					if res == nil || res.Body == nil {
						return
					}
					_ = res.Body.Close()
				}()
				return err != nil
			}, 5*time.Minute, 300*time.Millisecond)
		})
	}
}

func TestInstanceStartedTwiceWithoutLeaking(t *testing.T) {
	dialer := &net.Dialer{Timeout: time.Second}
	l, err := localstack.NewInstance()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, l.Stop())
	})
	require.NoError(t, l.Start())
	firstInstance := l.Endpoint(localstack.S3)
	require.NoError(t, l.Start())
	_, err = dialer.DialContext(t.Context(), "tcp", firstInstance)
	require.Error(t, err, "should be teared down")
}

func TestContextInstanceStartedTwiceWithoutLeaking(t *testing.T) {
	dialer := &net.Dialer{Timeout: time.Second}
	l, err := localstack.NewInstance()
	require.NoError(t, err)
	require.NoError(t, l.Start())
	firstInstance := l.Endpoint(localstack.S3)
	require.NoError(t, l.StartWithContext(t.Context()))
	_, err = dialer.DialContext(t.Context(), "tcp", firstInstance)
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
	t.Setenv("DOCKER_URL", "what-is-this-thing:///var/run/not-a-valid-docker.sock")
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
	host := os.Getenv("DOCKER_HOST")
	if host == "" || strings.Contains(host, "podman.sock") {
		t.Skip()
	}
	for _, s := range []struct {
		name        string
		given       func(t *testing.T)
		expectOpt   func(t require.TestingT, opt localstack.InstanceOption, err error)
		expectStart func(t require.TestingT, err error)
	}{
		{
			name: "is ok with client from env",
			given: func(t *testing.T) {
				t.Setenv("DOCKER_API_VERSION", "0")
			},
			expectOpt: func(t require.TestingT, opt localstack.InstanceOption, err error) {
				require.NoError(t, err)
				require.NotNil(t, opt)
			},
			expectStart: func(t require.TestingT, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "client version 0 is too old.")
			},
		},
		{
			name: "publishes errors",
			given: func(t *testing.T) {
				t.Setenv("DOCKER_HOST", "localhost")
			},
			expectOpt: func(t require.TestingT, opt localstack.InstanceOption, err error) {
				require.EqualError(t, err, "localstack: could not connect to docker: unable to parse docker host `localhost`")
				require.Nil(t, opt)
			},
		},
	} {
		t.Run(s.name, func(t *testing.T) {
			s.given(t)
			opt, err := localstack.WithClientFromEnv()
			s.expectOpt(t, opt, err)
			if s.expectStart != nil {
				i, err := localstack.NewInstance(opt)
				require.NoError(t, err)
				s.expectStart(t, i.StartWithContext(t.Context()))
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

func atLeastOneContainerMatchesLabels(labels map[string]string, containers []container.Summary) bool {
	for _, c := range containers {
		if matchesLabels(labels, c) {
			return true
		}
	}
	return false
}

func matchesLabels(labels map[string]string, container container.Summary) bool {
	for k, v := range labels {
		val, exists := container.Labels[k]
		if !exists || v != val {
			return false
		}
	}
	return true
}

func clean() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timeout := int(time.Second.Seconds())
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	cli.NegotiateAPIVersion(ctx)

	list, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return err
	}

	for _, l := range list {
		if err := cli.ContainerStop(ctx, l.ID, container.StopOptions{Timeout: &timeout}); err != nil {
			log.Println(err)
		}
	}

	if _, err := cli.ContainersPrune(ctx, filters.Args{}); err != nil {
		log.Println(err)
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
