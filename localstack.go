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

// Instance manages the localstack
type Instance struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

// Start starts the localstack
func (l *Instance) Start() error {
	if isAlreadyRunning(l) {
		if err := tearDown(l); err != nil {
			return err
		}
	}

	if err := startLocalstack(l); err != nil {
		return err
	}

	if err := l.pool.Retry(func() error {
		return isAvailable(l)
	}); err != nil {
		return fmt.Errorf("localstack: could start environment: %w", err)
	}

	return nil
}

// Stop stops the localstack
func (l *Instance) Stop() error {
	if l.pool != nil && l.resource != nil {
		return l.pool.Purge(l.resource)
	}
	return nil
}

// Endpoint returns the endpoint for the given service
// Endpoints are allocated dynamically (to avoid blocked ports), but are fix after starting the instance
func (l *Instance) Endpoint(service Service) string {
	if l.resource != nil {
		return l.resource.GetHostPort(string(service))
	}
	return ""
}

// Service represents an AWS service
type Service string

// Supported AWS/localstack services
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

func isAvailable(l *Instance) error {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("not", "empty", ""),
		DisableSSL:  aws.Bool(true),
		Region:      aws.String(endpoints.UsWest1RegionID),
		Endpoint:    aws.String(l.Endpoint(SQS)),
	})
	if err != nil {
		fmt.Println("localstack: waiting on server to start...")
		return err
	}

	s := sqs.New(sess)
	createQueue, err := s.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String("test-resource"),
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
	return nil
}

func startLocalstack(l *Instance) error {
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
	return nil
}

func tearDown(l *Instance) error {
	if err := l.Stop(); err != nil {
		return fmt.Errorf("localstack: can't stop an already running instance: %w", err)
	}
	return nil
}

func isAlreadyRunning(l *Instance) bool {
	return l.pool != nil
}
