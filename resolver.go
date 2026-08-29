package localstack

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/batch"
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
	"net/http"
	"net/url"
)

// NewCloudformationResolverV2 resolves the services ResolverV2 endpoint
func NewCloudformationResolverV2(i *Instance) *CloudformationResolverV2 {
	return &CloudformationResolverV2{i: i}
}

type CloudformationResolverV2 struct{ i *Instance }

func (c *CloudformationResolverV2) ResolveEndpoint(_ context.Context, _ cloudformation.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(CloudFormation))
}

// NewCloudwatchResolverV2 resolves the services ResolverV2 endpoint
func NewCloudwatchResolverV2(i *Instance) *CloudwatchResolverV2 {
	return &CloudwatchResolverV2{i: i}
}

type CloudwatchResolverV2 struct{ i *Instance }

func (c *CloudwatchResolverV2) ResolveEndpoint(_ context.Context, _ cloudwatch.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(CloudWatch))
}

// NewBatchResolver resolves the services ResolverV2 endpoint
func NewBatchResolverV2(i *Instance) *BatchResolverV2 {
	return &BatchResolverV2{i: i}
}
type BatchResolverV2 struct{ i *Instance }

func (b *BatchResolverV2) ResolveEndpoint(_ context.Context, _ batch.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(b.i.EndpointV2(Batch))
}

// NewCloudwatchLogsResolverV2 resolves the services ResolverV2 endpoint
func NewCloudwatchLogsResolverV2(i *Instance) *CloudwatchLogsResolverV2 {
	return &CloudwatchLogsResolverV2{i: i}
}

type CloudwatchLogsResolverV2 struct{ i *Instance }

func (c *CloudwatchLogsResolverV2) ResolveEndpoint(_ context.Context, _ cloudwatchlogs.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(CloudWatchLogs))
}

// NewCloudwatchEventsResolverV2 resolves the services ResolverV2 endpoint
func NewCloudwatchEventsResolverV2(i *Instance) *CloudwatchEventsResolverV2 {
	return &CloudwatchEventsResolverV2{i: i}
}

type CloudwatchEventsResolverV2 struct{ i *Instance }

func (c *CloudwatchEventsResolverV2) ResolveEndpoint(_ context.Context, _ cloudwatchevents.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(CloudWatchEvents))
}

// NewDynamoDbResolverV2 resolves the services ResolverV2 endpoint
func NewDynamoDbResolverV2(i *Instance) *DynamoDbResolver {
	return &DynamoDbResolver{i: i}
}

type DynamoDbResolver struct{ i *Instance }

func (c *DynamoDbResolver) ResolveEndpoint(_ context.Context, _ dynamodb.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(DynamoDB))
}

// NewDynamoDbStreamsResolverV2 resolves the services ResolverV2 endpoint
func NewDynamoDbStreamsResolverV2(i *Instance) *DynamoDbStreamsResolverV2 {
	return &DynamoDbStreamsResolverV2{i: i}
}

type DynamoDbStreamsResolverV2 struct{ i *Instance }

func (c *DynamoDbStreamsResolverV2) ResolveEndpoint(_ context.Context, _ dynamodbstreams.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(DynamoDBStreams))
}

// NewEc2ResolverV2 resolves the services ResolverV2 endpoint
func NewEc2ResolverV2(i *Instance) *Ec2ResolverV2 {
	return &Ec2ResolverV2{i: i}
}

type Ec2ResolverV2 struct{ i *Instance }

func (c *Ec2ResolverV2) ResolveEndpoint(_ context.Context, _ ec2.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(EC2))
}

// NewElasticSearchResolverV2 resolves the services ResolverV2 endpoint
func NewElasticSearchResolverV2(i *Instance) *ElasticSearchResolverV2 {
	return &ElasticSearchResolverV2{i: i}
}

type ElasticSearchResolverV2 struct{ i *Instance }

func (c *ElasticSearchResolverV2) ResolveEndpoint(_ context.Context, _ elasticsearchservice.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(ES))
}

// NewFirehoseResolverV2 resolves the services ResolverV2 endpoint
func NewFirehoseResolverV2(i *Instance) *FirehoseResolverV2 {
	return &FirehoseResolverV2{i: i}
}

type FirehoseResolverV2 struct{ i *Instance }

func (c *FirehoseResolverV2) ResolveEndpoint(_ context.Context, _ firehose.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(Firehose))
}

// NewIamResolverV2 resolves the services ResolverV2 endpoint
func NewIamResolverV2(i *Instance) *IamResolverV2 {
	return &IamResolverV2{i: i}
}

type IamResolverV2 struct{ i *Instance }

func (c *IamResolverV2) ResolveEndpoint(_ context.Context, _ iam.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(IAM))
}

// NewKinesisResolverV2 resolves the services ResolverV2 endpoint
func NewKinesisResolverV2(i *Instance) *KinesisResolverV2 {
	return &KinesisResolverV2{i: i}
}

type KinesisResolverV2 struct{ i *Instance }

func (c *KinesisResolverV2) ResolveEndpoint(_ context.Context, _ kinesis.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(Kinesis))
}

// NewLambdaResolverV2 resolves the services ResolverV2 endpoint
func NewLambdaResolverV2(i *Instance) *LambdaResolverV2 {
	return &LambdaResolverV2{i: i}
}

type LambdaResolverV2 struct{ i *Instance }

func (c *LambdaResolverV2) ResolveEndpoint(_ context.Context, _ lambda.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(Lambda))
}

// NewRedshiftResolverV2 resolves the services ResolverV2 endpoint
func NewRedshiftResolverV2(i *Instance) *RedshiftResolverV2 {
	return &RedshiftResolverV2{i: i}
}

type RedshiftResolverV2 struct{ i *Instance }

func (c *RedshiftResolverV2) ResolveEndpoint(_ context.Context, _ redshift.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(Redshift))
}

// NewRoute53ResolverV2 resolves the services ResolverV2 endpoint
func NewRoute53ResolverV2(i *Instance) *Route53ResolverV2 {
	return &Route53ResolverV2{i: i}
}

type Route53ResolverV2 struct{ i *Instance }

func (c *Route53ResolverV2) ResolveEndpoint(_ context.Context, _ route53.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(Route53))
}

// NewS3ResolverV2 resolves the services ResolverV2 endpoint
func NewS3ResolverV2(i *Instance) *S3ResolverV2 {
	return &S3ResolverV2{i: i}
}

type S3ResolverV2 struct{ i *Instance }

func (c *S3ResolverV2) ResolveEndpoint(_ context.Context, _ s3.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(S3))
}

// NewSecretsManagerResolverV2 resolves the services ResolverV2 endpoint
func NewSecretsManagerResolverV2(i *Instance) *SecretsManagerResolverV2 {
	return &SecretsManagerResolverV2{i: i}
}

type SecretsManagerResolverV2 struct{ i *Instance }

func (c *SecretsManagerResolverV2) ResolveEndpoint(_ context.Context, _ secretsmanager.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(SecretsManager))
}

// NewSesResolverV2 resolves the services ResolverV2 endpoint
func NewSesResolverV2(i *Instance) *SesResolverV2 {
	return &SesResolverV2{i: i}
}

type SesResolverV2 struct{ i *Instance }

func (c *SesResolverV2) ResolveEndpoint(_ context.Context, _ ses.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(SES))
}

// NewSnsResolverV2 resolves the services ResolverV2 endpoint
func NewSnsResolverV2(i *Instance) *SnsResolverV2 {
	return &SnsResolverV2{i: i}
}

type SnsResolverV2 struct{ i *Instance }

func (c *SnsResolverV2) ResolveEndpoint(_ context.Context, _ sns.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(SNS))
}

// NewSqsResolverV2 resolves the services ResolverV2 endpoint
func NewSqsResolverV2(i *Instance) *SqsResolverV2 {
	return &SqsResolverV2{i: i}
}

type SqsResolverV2 struct{ i *Instance }

func (c *SqsResolverV2) ResolveEndpoint(_ context.Context, _ sqs.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(SQS))
}

// NewSsmResolverV2 resolves the services ResolverV2 endpoint
func NewSsmResolverV2(i *Instance) *SsmResolverV2 {
	return &SsmResolverV2{i: i}
}

type SsmResolverV2 struct{ i *Instance }

func (c *SsmResolverV2) ResolveEndpoint(_ context.Context, _ ssm.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(SSM))
}

// NewStsResolverV2 resolves the services ResolverV2 endpoint
func NewStsResolverV2(i *Instance) *StsResolverV2 {
	return &StsResolverV2{i: i}
}

type StsResolverV2 struct{ i *Instance }

func (c *StsResolverV2) ResolveEndpoint(_ context.Context, _ sts.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(STS))
}

// NewStepFunctionsResolverV2 resolves the services ResolverV2 endpoint
func NewStepFunctionsResolverV2(i *Instance) *StepFunctionsResolverV2 {
	return &StepFunctionsResolverV2{i: i}
}

type StepFunctionsResolverV2 struct{ i *Instance }

func (c *StepFunctionsResolverV2) ResolveEndpoint(_ context.Context, _ sfn.EndpointParameters) (smithyendpoints.Endpoint, error) {
	return resolveEndpoint(c.i.EndpointV2(StepFunctions))
}

func resolveEndpoint(endpoint string) (smithyendpoints.Endpoint, error) {
	uri, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return smithyendpoints.Endpoint{}, fmt.Errorf("failed to parse uri: %s", endpoint)
	}
	return smithyendpoints.Endpoint{
		URI:     *uri,
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
	}, nil
}
