package localstack_test

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	smithyauth "github.com/aws/smithy-go/auth"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/elgohr/go-localstack"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestEndpointResolversV2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l, err := localstack.NewInstance()
	require.NoError(t, err)
	require.NoError(t, l.Start())
	t.Cleanup(func() {
		require.NoError(t, l.Stop())
	})

	for service, resolver := range map[localstack.Service]func() (smithyendpoints.Endpoint, error){
		localstack.CloudFormation: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewCloudformationResolverV2(l).ResolveEndpoint(ctx, cloudformation.EndpointParameters{})
		},
		localstack.CloudWatch: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewCloudwatchResolverV2(l).ResolveEndpoint(ctx, cloudwatch.EndpointParameters{})
		},
		localstack.CloudWatchLogs: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewCloudwatchLogsResolverV2(l).ResolveEndpoint(ctx, cloudwatchlogs.EndpointParameters{})
		},
		localstack.CloudWatchEvents: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewCloudwatchEventsResolverV2(l).ResolveEndpoint(ctx, cloudwatchevents.EndpointParameters{})
		},
		localstack.DynamoDB: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewDynamoDbResolverV2(l).ResolveEndpoint(ctx, dynamodb.EndpointParameters{})
		},
		localstack.DynamoDBStreams: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewDynamoDbStreamsResolverV2(l).ResolveEndpoint(ctx, dynamodbstreams.EndpointParameters{})
		},
		localstack.EC2: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewEc2ResolverV2(l).ResolveEndpoint(ctx, ec2.EndpointParameters{})
		},
		localstack.ES: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewElasticSearchResolverV2(l).ResolveEndpoint(ctx, elasticsearchservice.EndpointParameters{})
		},
		localstack.Firehose: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewFirehoseResolverV2(l).ResolveEndpoint(ctx, firehose.EndpointParameters{})
		},
		localstack.IAM: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewIamResolverV2(l).ResolveEndpoint(ctx, iam.EndpointParameters{})
		},
		localstack.Kinesis: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewKinesisResolverV2(l).ResolveEndpoint(ctx, kinesis.EndpointParameters{})
		},
		localstack.Lambda: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewLambdaResolverV2(l).ResolveEndpoint(ctx, lambda.EndpointParameters{})
		},
		localstack.Redshift: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewRedshiftResolverV2(l).ResolveEndpoint(ctx, redshift.EndpointParameters{})
		},
		localstack.Route53: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewRoute53ResolverV2(l).ResolveEndpoint(ctx, route53.EndpointParameters{})
		},
		localstack.S3: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewS3ResolverV2(l).ResolveEndpoint(ctx, s3.EndpointParameters{})
		},
		localstack.SecretsManager: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewSecretsManagerResolverV2(l).ResolveEndpoint(ctx, secretsmanager.EndpointParameters{})
		},
		localstack.SES: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewSesResolverV2(l).ResolveEndpoint(ctx, ses.EndpointParameters{})
		},
		localstack.SNS: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewSnsResolverV2(l).ResolveEndpoint(ctx, sns.EndpointParameters{})
		},
		localstack.SQS: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewSqsResolverV2(l).ResolveEndpoint(ctx, sqs.EndpointParameters{})
		},
		localstack.SSM: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewSsmResolverV2(l).ResolveEndpoint(ctx, ssm.EndpointParameters{})
		},
		localstack.STS: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewStsResolverV2(l).ResolveEndpoint(ctx, sts.EndpointParameters{})
		},
		localstack.StepFunctions: func() (smithyendpoints.Endpoint, error) {
			return localstack.NewStepFunctionsResolverV2(l).ResolveEndpoint(ctx, sfn.EndpointParameters{})
		},
	} {
		t.Run(service.Name, func(t *testing.T) {
			u, err := url.ParseRequestURI(l.EndpointV2(service))
			require.NoError(t, err)
			endpoint, err := resolver()
			require.NoError(t, err)
			expected := smithyendpoints.Endpoint{
				URI:     *u,
				Headers: http.Header{},
				Properties: func() smithy.Properties {
					var out smithy.Properties
					smithyauth.SetAuthOptions(&out, []*smithyauth.Option{
						{
							SchemeID: "aws.auth#sigv4",
							SignerProperties: func() smithy.Properties {
								var sp smithy.Properties
								smithyhttp.SetSigV4SigningName(&sp, "dynamodb")
								smithyhttp.SetSigV4ASigningName(&sp, "dynamodb")
								smithyhttp.SetSigV4SigningRegion(&sp, "us-east-1")
								return sp
							}(),
						},
					})
					return out
				}(),
			}
			require.Equal(t, expected, endpoint)
		})
	}
}
