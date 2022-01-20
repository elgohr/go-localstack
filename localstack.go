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
	"io"
	"log"
	"os"
	"time"

	"github.com/elgohr/go-localstack/internal"
)

// Instance manages the localstack
type Instance struct {
	cli         internal.DockerClient
	containerId string
	portMapping map[Service]string
	version     string
	fixedPort   bool
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

// StartWithContext starts the localstack and ends it when the context is done
func (i *Instance) StartWithContext(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		if err := i.stop(); err != nil {
			log.Println(err)
		}
	}()
	return i.start(ctx)
}

// Stop stops the localstack
// Deprecated: Use StartWithContext instead.
func (i *Instance) Stop() error {
	return i.stop()
}

// Endpoint returns the endpoint for the given service
// Endpoints are allocated dynamically (to avoid blocked ports), but are fix after starting the instance
func (i *Instance) Endpoint(service Service) string {
	if i.containerId != "" {
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
	if i.containerId != "" {
		if i.fixedPort {
			return "http://" + i.portMapping[FixedPort]
		}
		return "http://" + i.portMapping[service]
	}
	return ""
}

// ContainerId returns the deployed container's ID
func (i *Instance) ContainerId() string {
	return i.containerId
}

// Service represents an AWS service
type Service string

// Supported AWS/localstack services
const (
	FixedPort = Service("4566/tcp")

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

func (i *Instance) start(ctx context.Context) error {
	if i.isAlreadyRunning() {
		log.Println("stopping an instance that is already running")
		if err := i.stop(); err != nil {
			return fmt.Errorf("localstack: can't stop an already running instance: %w", err)
		}
	}

	if err := i.startLocalstack(ctx); err != nil {
		return err
	}

	log.Println("waiting for localstack to start...")
	return i.waitToBeAvailable(ctx)
}

func (i *Instance) startLocalstack(ctx context.Context) error {
	imageName := "localstack/localstack:" + i.version
	if !i.isDownloaded(ctx) {
		reader, err := i.cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			return fmt.Errorf("localstack: could not load image: %w", err)
		}
		defer func() {
			if err := reader.Close(); err != nil {
				log.Println(err)
			}
		}()

		// for reading the load output
		if _, err = io.Copy(os.Stdout, reader); err != nil {
			return fmt.Errorf("localstack: %w", err)
		}
	}

	pm := nat.PortMap{}
	for service := range AvailableServices {
		pm[nat.Port(service)] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: ""}}
	}

	resp, err := i.cli.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
		}, &container.HostConfig{
			PortBindings: pm,
			AutoRemove:   true,
		}, nil, nil, "")
	if err != nil {
		return fmt.Errorf("localstack: could not create container: %w", err)
	}
	i.containerId = resp.ID

	log.Println("starting localstack")
	if err := i.cli.ContainerStart(ctx, i.containerId, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("localstack: could not start container: %w", err)
	}

	startedContainer, err := i.cli.ContainerInspect(ctx, i.containerId)
	if err != nil {
		return fmt.Errorf("localstack: could not get port from container: %w", err)
	}
	ports := startedContainer.NetworkSettings.Ports
	if i.fixedPort {
		i.portMapping[FixedPort] = "localhost:" + ports[nat.Port(FixedPort)][0].HostPort
	} else {
		for service := range AvailableServices {
			i.portMapping[service] = "localhost:" + ports[nat.Port(service)][0].HostPort
		}
	}

	return nil
}

func (i *Instance) stop() error {
	if i.containerId == "" {
		return nil
	}
	timeout := 5 * time.Second
	if err := i.cli.ContainerStop(context.Background(), i.containerId, &timeout); err != nil {
		return err
	}
	i.containerId = ""
	i.portMapping = map[Service]string{}
	return nil
}

func (i *Instance) isDownloaded(ctx context.Context) bool {
	list, err := i.cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		log.Println(err)
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
				log.Println("localstack: finished waiting")
				return nil
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
	return i.containerId != ""
}
