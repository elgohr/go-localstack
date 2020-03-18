package localstack

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ory/dockertest"
)

type LocalStack struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (l *LocalStack) Start() error {
	var err error
	l.pool, err = dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("localstack: could not connect to docker: %w", err)
	}
	l.resource, err = l.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "localstack/localstack",
		Tag:        "latest",
	})
	if err != nil {
		return fmt.Errorf("localstack: could not start resource: %w", err)
	}
	address := l.resource.GetHostPort("4576/tcp")

	if err := l.pool.Retry(func() error {
		var err error
		sess, err := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials("not", "empty", ""),
			DisableSSL:  aws.Bool(true),
			Region:      aws.String(endpoints.UsWest1RegionID),
			Endpoint:    aws.String(address),
		})
		if err != nil {
			fmt.Println("localstack: waiting on server to start...")
			return err
		}

		s := sqs.New(sess)

		createQueue, err := s.CreateQueue(&sqs.CreateQueueInput{
			QueueName: aws.String("queue"),
			Attributes: map[string]*string{
				"VisibilityTimeout": aws.String("1"),
			},
		})
		if err != nil {
			fmt.Println("localstack: waiting on server to initialize...")
			return err
		}

		if _, err := s.DeleteQueue(&sqs.DeleteQueueInput{
			QueueUrl: createQueue.QueueUrl,
		}); err != nil {
			return err
		}

		fmt.Println("localstack: finished waiting")
		return err
	}); err != nil {
		return fmt.Errorf("localstack: could start environment: %w", err)
	}
	return nil
}

func (l *LocalStack) Stop() error {
	return l.pool.Purge(l.resource)
}

func (l *LocalStack) Endpoint(service Service) string {
	return l.resource.GetHostPort(string(service))
}

type Service string

const (
	CloudFormation   = Service("4581/tcp")
	CloudWatch       = Service("4582/tcp")
	CloudWatchLogs   = Service("4586/tcp")
	CloudWatchEvents = Service("4587/tcp")
	DynamoDB         = Service("4569/tcp")
	DynamoDBStreams  = Service("4570/tcp")
	EC2              = Service("4597/tcp")
	ES               = Service("4578/tcp")
	Firehose         = Service("4573/tcp")
	IAM              = Service("4593/tcp")
	Kinesis          = Service("4568/tcp")
	KMS              = Service("4599/tcp")
	Lambda           = Service("4574/tcp")
	Redshift         = Service("4577/tcp")
	Route53          = Service("4580/tcp")
	S3               = Service("4572/tcp")
	SecretsManager   = Service("4584/tcp")
	SES              = Service("4579/tcp")
	SNS              = Service("4575/tcp")
	SQS              = Service("4576/tcp")
	SSM              = Service("4583/tcp")
	STS              = Service("4592/tcp")
	StepFunctions    = Service("4585/tcp")
)
