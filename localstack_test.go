package localstack_test

import (
	"github.com/elgohr/go-localstack"
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
