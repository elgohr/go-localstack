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
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/elgohr/go-localstack"
)

func TestMain(m *testing.M) {
	if err := clean(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(m.Run())
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

func clean() error {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := cli.ContainersPrune(ctx, filters.Args{}); err != nil {
		return err
	}
	if _, err := cli.NetworksPrune(ctx, filters.Args{}); err != nil {
		return err
	}
	if _, err := cli.VolumesPrune(ctx, filters.Args{}); err != nil {
		return err
	}
	if _, err := cli.BuildCachePrune(ctx, types.BuildCachePruneOptions{All: true}); err != nil {
		return err
	}
	if _, err := cli.ImagesPrune(ctx, filters.Args{}); err != nil {
		return err
	}
	return nil
}
