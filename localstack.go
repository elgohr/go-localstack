package localstack

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/ory/dockertest/v3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/elgohr/go-localstack/internal"
)

// Instance manages the localstack
type Instance struct {
	pool      internal.Pool
	resource  *dockertest.Resource
	version   string
	fixedPort bool
}

// InstanceOption is an option that controls the behaviour of
// localstack.
type InstanceOption func(i *Instance)

// WithVersion configures the instance to use a specific version of
// localstack. Must be a valid version string or "latest".
func WithVersion(version string) InstanceOption {
	return func(i *Instance) {
		i.version = version
	}
}

// Semver constraint that tests it the version is affected by the port change.
var portChangeIntroduced = internal.MustParseConstraint(">= 0.11.5")

// NewInstance creates a new Instance
// Fails when Docker is not reachable
func NewInstance(opts ...InstanceOption) (*Instance, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("localstack: could not connect to docker: %w", err)
	}

	i := Instance{
		pool:    pool,
		version: "latest",
	}

	for _, opt := range opts {
		opt(&i)
	}

	if i.version == "latest" {
		i.fixedPort = true
	} else {
		version, err := semver.NewVersion(i.version)
		if err != nil {
			return nil, fmt.Errorf("localstack: invalid version %q specified: %w", i.version, err)
		}

		i.version = version.String()
		i.fixedPort = portChangeIntroduced.Check(version)
	}

	return &i, nil
}

// Start starts the localstack
func (i *Instance) Start() error {
	if i.isAlreadyRunning() {
		if err := i.tearDown(); err != nil {
			return err
		}
	}

	if err := i.startLocalstack(); err != nil {
		return err
	}

	if err := i.pool.Retry(func() error {
		return i.isAvailable()
	}); err != nil {
		return fmt.Errorf("localstack: could not start environment: %w", err)
	}

	return nil
}

// Stop stops the localstack
func (i Instance) Stop() error {
	if i.resource != nil {
		return i.pool.Purge(i.resource)
	}
	return nil
}

// Endpoint returns the endpoint for the given service
// Endpoints are allocated dynamically (to avoid blocked ports), but are fix after starting the instance
func (i Instance) Endpoint(service Service) string {
	if i.resource != nil {
		if i.fixedPort {
			return i.resource.GetHostPort("4566/tcp")
		}

		return i.resource.GetHostPort(string(service))
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

func (i Instance) isAvailable() error {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("not", "empty", ""),
		DisableSSL:  aws.Bool(true),
		Region:      aws.String(endpoints.UsWest1RegionID),
		Endpoint:    aws.String(i.Endpoint(SQS)),
	})
	if err != nil {
		fmt.Println("localstack: waiting on server to start...")
		return err
	}

	s := sqs.New(sess)
	createQueue, err := s.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String("test-Resource"),
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

func (i *Instance) startLocalstack() error {
	var err error
	i.resource, err = i.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "localstack/localstack",
		Tag:        i.version,
	})
	if err != nil {
		return fmt.Errorf("localstack: could not start container: %w", err)
	}
	return nil
}

func (i Instance) tearDown() error {
	if err := i.Stop(); err != nil {
		return fmt.Errorf("localstack: can't stop an already running instance: %w", err)
	}
	return nil
}

func (i Instance) isAlreadyRunning() bool {
	return i.pool != nil
}
