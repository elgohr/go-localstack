package localstack_test

import (
	"errors"
	"github.com/elgohr/go-localstack"
	golocalstackfakes "github.com/elgohr/go-localstack/go-localstackfakes"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"net"
	"strings"
	"testing"
)

func TestLocalStack(t *testing.T) {
	l, err := localstack.NewInstance()
	assert.NoError(t, err)
	assert.NoError(t, l.Start())
	defer l.Stop()

	assert.True(t, strings.HasPrefix(l.Endpoint(localstack.CloudFormation), "localhost:"), l.Endpoint(localstack.CloudFormation))
	assert.NotEmpty(t, l.Endpoint(localstack.CloudFormation))
	assert.NotEmpty(t, l.Endpoint(localstack.CloudWatch))
	assert.NotEmpty(t, l.Endpoint(localstack.CloudWatchLogs))
	assert.NotEmpty(t, l.Endpoint(localstack.CloudWatchEvents))
	assert.NotEmpty(t, l.Endpoint(localstack.DynamoDB))
	assert.NotEmpty(t, l.Endpoint(localstack.DynamoDBStreams))
	assert.NotEmpty(t, l.Endpoint(localstack.EC2))
	assert.NotEmpty(t, l.Endpoint(localstack.ES))
	assert.NotEmpty(t, l.Endpoint(localstack.Firehose))
	assert.NotEmpty(t, l.Endpoint(localstack.IAM))
	assert.NotEmpty(t, l.Endpoint(localstack.Kinesis))
	assert.NotEmpty(t, l.Endpoint(localstack.Lambda))
	assert.NotEmpty(t, l.Endpoint(localstack.Redshift))
	assert.NotEmpty(t, l.Endpoint(localstack.Route53))
	assert.NotEmpty(t, l.Endpoint(localstack.S3))
	assert.NotEmpty(t, l.Endpoint(localstack.SecretsManager))
	assert.NotEmpty(t, l.Endpoint(localstack.SES))
	assert.NotEmpty(t, l.Endpoint(localstack.SNS))
	assert.NotEmpty(t, l.Endpoint(localstack.SQS))
	assert.NotEmpty(t, l.Endpoint(localstack.SSM))
	assert.NotEmpty(t, l.Endpoint(localstack.STS))
	assert.NotEmpty(t, l.Endpoint(localstack.StepFunctions))
}

func TestInstanceStartedTwiceWithoutLeaking(t *testing.T) {
	l, err := localstack.NewInstance()
	assert.NoError(t, err)
	defer l.Stop()
	assert.NoError(t, l.Start())
	firstInstance := l.Endpoint(localstack.S3)
	assert.NoError(t, l.Start())
	_, err = net.Dial("tcp", firstInstance)
	assert.Error(t, err, "should be teared down")
}

func TestInstanceStopWithoutStarted(t *testing.T) {
	l, err := localstack.NewInstance()
	assert.NoError(t, err)
	assert.NoError(t, l.Stop())
}

func TestInstanceEndpointWithoutStarted(t *testing.T) {
	l, err := localstack.NewInstance()
	assert.NoError(t, err)
	assert.Empty(t, l.Endpoint(localstack.S3))
}

func TestInstance_Start_Fails(t *testing.T) {
	for _, tt := range [...]struct {
		name  string
		given func() *localstack.Instance
		then  func(err error)
	}{
		{
			name: "can't restart localstack when already running",
			given: func() *localstack.Instance {
				fakePool := &golocalstackfakes.FakePool{}
				fakePool.PurgeReturns(errors.New("can't start"))
				return &localstack.Instance{
					Pool:     fakePool,
					Resource: &dockertest.Resource{},
				}
			},
			then: func(err error) {
				assert.Equal(t, "localstack: can't stop an already running instance: can't start", err.Error())
			},
		},
		{
			name: "can't start container",
			given: func() *localstack.Instance {
				fakePool := &golocalstackfakes.FakePool{}
				fakePool.RunWithOptionsReturns(nil, errors.New("can't start container"))
				return &localstack.Instance{
					Pool: fakePool,
				}
			},
			then: func(err error) {
				assert.Equal(t, "localstack: could not start container: can't start container", err.Error())
			},
		},
		{
			name: "fails during waiting on startup",
			given: func() *localstack.Instance {
				fakePool := &golocalstackfakes.FakePool{}
				fakePool.RetryReturns(errors.New("can't wait"))
				return &localstack.Instance{
					Pool: fakePool,
				}
			},
			then: func(err error) {
				assert.Equal(t, "localstack: could not start environment: can't wait", err.Error())
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt.then(tt.given().Start())
		})
	}
}

func TestInstance_Stop_Fails(t *testing.T) {
	fakePool := &golocalstackfakes.FakePool{}
	fakePool.PurgeReturns(errors.New("can't stop"))
	i := &localstack.Instance{
		Pool:     fakePool,
		Resource: &dockertest.Resource{},
	}

	assert.Equal(t, "can't stop", i.Stop().Error())
}
