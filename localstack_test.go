package localstack_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/elgohr/go-localstack"
)

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
			require.NoError(t, l.Start())
			defer l.Stop()
			s.expect(t, l)
		})
	}
}

func TestInstanceStartedTwiceWithoutLeaking(t *testing.T) {
	l, err := localstack.NewInstance()
	require.NoError(t, err)
	defer l.Stop()
	require.NoError(t, l.Start())
	firstInstance := l.Endpoint(localstack.S3)
	require.NoError(t, l.Start())
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
	defer os.Setenv("DOCKER_URL", urlIfSet)

	os.Setenv("DOCKER_URL", "what-is-this-thing:///var/run/not-a-valid-docker.sock")

	_, err := localstack.NewInstance()
	require.Error(t, err)
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
	for _, service := range services {
		endpoints[l.Endpoint(service)] = struct{}{}
	}
	require.Equal(t, 1, len(endpoints), endpoints)
}

func havingIndividualEndpoints(t *testing.T, l *localstack.Instance) {
	endpoints := map[string]struct{}{}
	for _, service := range services {
		endpoint := l.Endpoint(service)
		checkAddress(t, endpoint)

		_, exists := endpoints[endpoint]
		require.False(t, exists, fmt.Sprintf("%s duplicated in %v", endpoint, endpoints))

		endpoints[endpoint] = struct{}{}
	}
	require.Equal(t, len(services), len(endpoints))
}

func checkAddress(t *testing.T, val string) {
	require.True(t, strings.HasPrefix(val, "localhost:"), val)
	require.NotEmpty(t, val[10:])
}

var services = []localstack.Service{
	localstack.CloudFormation,
	localstack.CloudWatch,
	localstack.CloudWatchLogs,
	localstack.CloudWatchEvents,
	localstack.DynamoDB,
	localstack.DynamoDBStreams,
	localstack.EC2,
	localstack.ES,
	localstack.Firehose,
	localstack.IAM,
	localstack.Kinesis,
	localstack.Lambda,
	localstack.Redshift,
	localstack.Route53,
	localstack.S3,
	localstack.SecretsManager,
	localstack.SES,
	localstack.SNS,
	localstack.SQS,
	localstack.SSM,
	localstack.STS,
	localstack.StepFunctions,
}
