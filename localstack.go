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

package localstack

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamo_types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"

	"github.com/elgohr/go-localstack/internal"
)

// Instance manages the localstack
type Instance struct {
	cli              internal.DockerClient
	portMapping      map[Service]string
	containerId      string
	containerIdMutex sync.RWMutex
	version          string
	fixedPort        bool
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
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, fmt.Errorf("localstack: could not connect to docker: %w", err)
	}

	i := Instance{
		cli:         cli,
		version:     "latest",
		portMapping: map[Service]string{},
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
// Deprecated: Use StartWithContext instead.
func (i *Instance) Start() error {
	return i.start(context.Background())
}

// StartWithContext starts the localstack and ends it when the context is done.
// Use it to also start individual services, by default all are started.
func (i *Instance) StartWithContext(ctx context.Context, services ...Service) error {
	go func() {
		<-ctx.Done()
		if err := i.stop(); err != nil {
			log.Error(err)
		}
	}()
	return i.start(ctx, services...)
}

// Stop stops the localstack
// Deprecated: Use StartWithContext instead.
func (i *Instance) Stop() error {
	return i.stop()
}

// Endpoint returns the endpoint for the given service
// Endpoints are allocated dynamically (to avoid blocked ports), but are fix after starting the instance
func (i *Instance) Endpoint(service Service) string {
	if i.getContainerId() != "" {
		if i.fixedPort {
			return i.portMapping[FixedPort]
		}
		return i.portMapping[service]
	}
	return ""
}

// EndpointV2 returns the endpoint for the given service when used by aws-sdk-v2
// Endpoints are allocated dynamically (to avoid blocked ports), but are fix after starting the instance
func (i *Instance) EndpointV2(service Service) string {
	if i.getContainerId() != "" {
		if i.fixedPort {
			return "http://" + i.portMapping[FixedPort]
		}
		return "http://" + i.portMapping[service]
	}
	return ""
}

// Service represents an AWS service
type Service struct {
	Name string
	Port string
}

// Supported AWS/localstack services
var (
	FixedPort = Service{Name: "all", Port: "4566/tcp"}

	CloudFormation   = Service{Name: "cloudformation", Port: "4581/tcp"}
	CloudWatch       = Service{Name: "cloudwatch", Port: "4582/tcp"}
	CloudWatchLogs   = Service{Name: "cloudwatchlogs", Port: "4586/tcp"}
	CloudWatchEvents = Service{Name: "cloudwatchevents", Port: "4587/tcp"}
	DynamoDB         = Service{Name: "dynamoDB", Port: "4569/tcp"}
	DynamoDBStreams  = Service{Name: "dynamoDBStreams", Port: "4570/tcp"}
	EC2              = Service{Name: "ec2", Port: "4597/tcp"}
	ES               = Service{Name: "es", Port: "4578/tcp"}
	Firehose         = Service{Name: "firehose", Port: "4573/tcp"}
	IAM              = Service{Name: "iam", Port: "4593/tcp"}
	Kinesis          = Service{Name: "kinesis", Port: "4568/tcp"}
	Lambda           = Service{Name: "lambda", Port: "4574/tcp"}
	Redshift         = Service{Name: "redshift", Port: "4577/tcp"}
	Route53          = Service{Name: "route53", Port: "4580/tcp"}
	S3               = Service{Name: "s3", Port: "4572/tcp"}
	SecretsManager   = Service{Name: "secretsmanager", Port: "4584/tcp"}
	SES              = Service{Name: "ses", Port: "4579/tcp"}
	SNS              = Service{Name: "sns", Port: "4575/tcp"}
	SQS              = Service{Name: "sqs", Port: "4576/tcp"}
	SSM              = Service{Name: "ssm", Port: "4583/tcp"}
	STS              = Service{Name: "sts", Port: "4592/tcp"}
	StepFunctions    = Service{Name: "stepfunctions", Port: "4585/tcp"}
)

// AvailableServices provides a map of all services for faster searches
var AvailableServices = map[Service]struct{}{
	FixedPort:        {},
	CloudFormation:   {},
	CloudWatch:       {},
	CloudWatchLogs:   {},
	CloudWatchEvents: {},
	DynamoDB:         {},
	DynamoDBStreams:  {},
	EC2:              {},
	ES:               {},
	Firehose:         {},
	IAM:              {},
	Kinesis:          {},
	Lambda:           {},
	Redshift:         {},
	Route53:          {},
	S3:               {},
	SecretsManager:   {},
	SES:              {},
	SNS:              {},
	SQS:              {},
	SSM:              {},
	STS:              {},
	StepFunctions:    {},
}

func (i *Instance) start(ctx context.Context, services ...Service) error {
	if i.isAlreadyRunning() {
		log.Info("stopping an instance that is already running")
		if err := i.stop(); err != nil {
			return fmt.Errorf("localstack: can't stop an already running instance: %w", err)
		}
	}

	if err := i.startLocalstack(ctx, services...); err != nil {
		return err
	}

	log.Info("waiting for localstack to start...")
	return i.waitToBeAvailable(ctx)
}

func (i *Instance) startLocalstack(ctx context.Context, services ...Service) error {
	imageName := "localstack/localstack:" + i.version
	if !i.isDownloaded(ctx) {
		reader, err := i.cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			return fmt.Errorf("localstack: could not load image: %w", err)
		}
		defer func() {
			if err := reader.Close(); err != nil {
				log.Error(err)
			}
		}()

		// for reading the load output
		if _, err = io.Copy(log.StandardLogger().Out, reader); err != nil {
			return fmt.Errorf("localstack: %w", err)
		}
	}

	pm := nat.PortMap{}
	for service := range AvailableServices {
		pm[nat.Port(service.Port)] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: ""}}
	}

	environmentVariables := []string{}
	if len(services) > 0 {
		startServices := "SERVICES=dynamodb" // for waitToBeAvailable
		addedServices := 0
		for _, service := range services {
			if shouldBeAdded(service) {
				startServices += "," + service.Name
				addedServices++
			}
		}
		if addedServices > 0 {
			environmentVariables = append(environmentVariables, startServices)
		}
	}

	resp, err := i.cli.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
			Env:   environmentVariables,
		}, &container.HostConfig{
			PortBindings: pm,
			AutoRemove:   true,
		}, nil, nil, "")
	if err != nil {
		return fmt.Errorf("localstack: could not create container: %w", err)
	}

	i.setContainerId(resp.ID)

	log.Info("starting localstack")
	containerId := resp.ID
	if err := i.cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("localstack: could not start container: %w", err)
	}

	startedContainer, err := i.cli.ContainerInspect(ctx, containerId)
	if err != nil {
		return fmt.Errorf("localstack: could not get port from container: %w", err)
	}
	ports := startedContainer.NetworkSettings.Ports
	if i.fixedPort {
		i.portMapping[FixedPort] = "localhost:" + ports[nat.Port(FixedPort.Port)][0].HostPort
	} else {
		hasFilteredServices := len(services) > 0
		for service := range AvailableServices {
			if hasFilteredServices && containsService(services, service) {
				i.portMapping[service] = "localhost:" + ports[nat.Port(service.Port)][0].HostPort
			} else if !hasFilteredServices {
				i.portMapping[service] = "localhost:" + ports[nat.Port(service.Port)][0].HostPort
			}
		}
	}

	return nil
}

func (i *Instance) stop() error {
	containerId := i.getContainerId()
	if containerId == "" {
		return nil
	}
	timeout := time.Second
	if err := i.cli.ContainerStop(context.Background(), containerId, &timeout); err != nil {
		return err
	}
	i.setContainerId("")
	i.portMapping = map[Service]string{}
	return nil
}

func (i *Instance) isDownloaded(ctx context.Context) bool {
	list, err := i.cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		log.Error(err)
		return false
	}
	for _, image := range list {
		for _, tag := range image.RepoTags {
			if tag == "localstack/localstack:"+i.version {
				return true
			}
		}
	}
	return false
}

func (i *Instance) waitToBeAvailable(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := i.checkAvailable(ctx); err == nil {
				log.Info("localstack: finished waiting")
				return nil
			} else {
				log.Debug(err)
			}
		}
	}
}

func (i *Instance) checkAvailable(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(_, _ string, _ ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               i.EndpointV2(DynamoDB),
				SigningRegion:     "us-east-1",
				HostnameImmutable: true,
			}, nil
		})),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
	)
	if err != nil {
		return err
	}

	s := dynamodb.NewFromConfig(cfg)
	testTable := aws.String("bucket")
	if _, err = s.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []dynamo_types.AttributeDefinition{
			{AttributeName: aws.String("PK"), AttributeType: dynamo_types.ScalarAttributeTypeS},
		},
		KeySchema: []dynamo_types.KeySchemaElement{
			{AttributeName: aws.String("PK"), KeyType: dynamo_types.KeyTypeHash},
		},
		ProvisionedThroughput: &dynamo_types.ProvisionedThroughput{
			ReadCapacityUnits: aws.Int64(1), WriteCapacityUnits: aws.Int64(1),
		},
		TableName: testTable,
	}); err != nil {
		return err
	}

	_, err = s.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: testTable,
	})
	return err
}

func (i *Instance) isAlreadyRunning() bool {
	return i.getContainerId() != ""
}

func (i *Instance) getContainerId() string {
	i.containerIdMutex.RLock()
	defer i.containerIdMutex.RUnlock()
	return i.containerId
}

func (i *Instance) setContainerId(v string) {
	i.containerIdMutex.Lock()
	defer i.containerIdMutex.Unlock()
	i.containerId = v
}

func shouldBeAdded(service Service) bool {
	return service != DynamoDB && service != FixedPort && service != ES
}

func containsService(services []Service, service Service) bool {
	if service == DynamoDB {
		return true
	}
	for _, s := range services {
		if s == service {
			return true
		}
	}
	return false
}
