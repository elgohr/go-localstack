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
	"archive/tar"
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamotypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"sync"
	"time"

	"github.com/elgohr/go-localstack/internal"
)

// Instance manages the localstack
type Instance struct {
	cli internal.DockerClient
	log *logrus.Logger

	portMapping      map[Service]string
	portMappingMutex sync.RWMutex

	containerId      string
	containerIdMutex sync.RWMutex

	labels    map[string]string
	version   string
	fixedPort bool
	timeout   time.Duration
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

// WithLogger configures the instance to use the specified logger.
func WithLogger(logger *logrus.Logger) InstanceOption {
	return func(i *Instance) {
		i.log = logger
	}
}

// WithLabels configures the labels that will be applied on the instance.
func WithLabels(labels map[string]string) InstanceOption {
	return func(i *Instance) {
		i.labels = labels
	}
}

// WithTimeout configures the timeout for terminating the localstack instance.
// This was invented to prevent orphaned containers after panics.
// The default timeout is set to 5 minutes.
func WithTimeout(timeout time.Duration) InstanceOption {
	return func(i *Instance) {
		i.timeout = timeout
	}
}

// WithClientFromEnv configures the instance to use a client that respects environment variables.
func WithClientFromEnv() (InstanceOption, error) {
	return WithClientFromEnvCtx(context.Background())
}

// WithClientFromEnvCtx like WithClientFromEnv but with context
func WithClientFromEnvCtx(ctx context.Context) (InstanceOption, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("localstack: could not connect to docker: %w", err)
	}
	cli.NegotiateAPIVersion(ctx)
	return func(i *Instance) {
		i.cli = cli
	}, nil
}

// Semver constraint that tests it the version is affected by the port change.
var portChangeIntroduced = internal.MustParseConstraint(">= 0.11.5")

// NewInstance creates a new Instance
// Fails when Docker is not reachable
func NewInstance(opts ...InstanceOption) (*Instance, error) {
	return NewInstanceCtx(context.Background(), opts...)
}

// NewInstanceCtx is NewInstance, but with Context
func NewInstanceCtx(ctx context.Context, opts ...InstanceOption) (*Instance, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, fmt.Errorf("localstack: could not connect to docker: %w", err)
	}
	cli.NegotiateAPIVersion(ctx)

	i := Instance{
		cli:         cli,
		log:         logrus.StandardLogger(),
		version:     "latest",
		portMapping: map[Service]string{},
		timeout:     5 * time.Minute,
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
func (i *Instance) Start(services ...Service) error {
	return i.start(context.Background(), services...)
}

// StartWithContext starts the localstack and ends it when the context is done.
// Deprecated: Use Start/Stop instead, as shutdown is not reliable
func (i *Instance) StartWithContext(ctx context.Context, services ...Service) error {
	go func() {
		<-ctx.Done()
		if err := i.stop(); err != nil {
			i.log.Error(err)
		}
	}()
	return i.start(ctx, services...)
}

// Stop stops the localstack
func (i *Instance) Stop() error {
	return i.stop()
}

// Endpoint returns the endpoint for the given service
// Endpoints are allocated dynamically (to avoid blocked ports), but are fix after starting the instance
func (i *Instance) Endpoint(service Service) string {
	if i.getContainerId() != "" {
		if i.fixedPort {
			return i.getPortMapping(FixedPort)
		}
		return i.getPortMapping(service)
	}
	return ""
}

// EndpointV2 returns the endpoint for the given service when used by aws-sdk-v2
// Endpoints are allocated dynamically (to avoid blocked ports), but are fix after starting the instance
func (i *Instance) EndpointV2(service Service) string {
	if i.getContainerId() != "" {
		if i.fixedPort {
			return "http://" + i.getPortMapping(FixedPort)
		}
		return "http://" + i.getPortMapping(service)
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
		i.log.Info("stopping an instance that is already running")
		if err := i.stop(); err != nil {
			return fmt.Errorf("localstack: can't stop an already running instance: %w", err)
		}
	}
	return i.startContainer(ctx, services, 0)
}

func (i *Instance) startContainer(ctx context.Context, services []Service, try int) error {
	if err := i.startLocalstack(ctx, services...); err != nil {
		return err
	}

	i.log.Info("waiting for localstack to start...")
	err := i.waitToBeAvailable(ctx)
	if errors.As(err, &containerMissing{}) {
		if try > 3 {
			return err
		}
		i.log.Debugln("missing container retrying")
		return i.startContainer(ctx, services, try+1)
	}
	return err
}

const imageName = "go-localstack"

func (i *Instance) startLocalstack(ctx context.Context, services ...Service) error {
	if err := i.buildLocalImage(ctx); err != nil {
		return fmt.Errorf("localstack: could not build image: %w", err)
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
			Image:        imageName,
			Env:          environmentVariables,
			Labels:       i.labels,
			Tty:          true,
			AttachStdout: true,
			AttachStderr: true,
		}, &container.HostConfig{
			PortBindings: pm,
			AutoRemove:   true,
		}, nil, nil, "")
	if err != nil {
		return fmt.Errorf("localstack: could not create container: %w", err)
	}

	containerId := resp.ID
	i.setContainerId(containerId)

	i.log.Info("starting localstack")
	if err := i.cli.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
		return fmt.Errorf("localstack: could not start container: %w", err)
	}

	if i.log.Level == logrus.DebugLevel {
		go i.writeContainerLogToLogger(ctx, containerId)
	}

	return i.mapPorts(ctx, services, containerId, 0)
}

//go:embed Dockerfile
var dockerTemplate string

func (i *Instance) buildLocalImage(ctx context.Context) error {
	buf := &bytes.Buffer{}
	tw := tar.NewWriter(buf)
	defer logClose(tw)

	dockerFile := "Dockerfile"
	dockerFileContent := []byte(fmt.Sprintf(dockerTemplate, i.version, int(i.timeout.Seconds())))
	if err := tw.WriteHeader(&tar.Header{
		Name: dockerFile,
		Size: int64(len(dockerFileContent)),
	}); err != nil {
		return err
	}

	if _, err := tw.Write(dockerFileContent); err != nil {
		return err
	}

	dockerFileTarReader := bytes.NewReader(buf.Bytes())
	imageBuildResponse, err := i.cli.ImageBuild(ctx, dockerFileTarReader, build.ImageBuildOptions{
		Tags:           []string{imageName},
		Dockerfile:     dockerFile,
		SuppressOutput: true,
		Remove:         true,
		ForceRemove:    true,
	})
	if err != nil {
		return err
	}
	defer logClose(imageBuildResponse.Body)

	_, err = io.Copy(io.Discard, imageBuildResponse.Body)
	return err
}

func (i *Instance) mapPorts(ctx context.Context, services []Service, containerId string, try int) error {
	if try > 10 {
		return errors.New("localstack: could not get port from container")
	}
	startedContainer, err := i.cli.ContainerInspect(ctx, containerId)
	if err != nil {
		return fmt.Errorf("localstack: could not inspect container: %w", err)
	}
	ports := startedContainer.NetworkSettings.Ports
	if i.fixedPort {
		bindings := ports[nat.Port(FixedPort.Port)]
		if len(bindings) == 0 {
			time.Sleep(300 * time.Millisecond)
			return i.mapPorts(ctx, services, containerId, try+1)
		}
		i.savePortMappings(map[Service]string{
			FixedPort: "localhost:" + bindings[0].HostPort,
		})
	} else {
		hasFilteredServices := len(services) > 0
		newMapping := make(map[Service]string, len(AvailableServices))
		for service := range AvailableServices {
			bindings := ports[nat.Port(service.Port)]
			if len(bindings) == 0 {
				time.Sleep(300 * time.Millisecond)
				return i.mapPorts(ctx, services, containerId, try+1)
			}
			if hasFilteredServices && containsService(services, service) {
				newMapping[service] = "localhost:" + bindings[0].HostPort
			} else if !hasFilteredServices {
				newMapping[service] = "localhost:" + bindings[0].HostPort
			}
		}
		i.savePortMappings(newMapping)
	}
	return nil
}

func (i *Instance) stop() error {
	i.containerIdMutex.Lock()
	defer i.containerIdMutex.Unlock()
	if i.containerId == "" {
		return nil
	}
	if err := i.cli.ContainerStop(context.Background(), i.containerId, container.StopOptions{
		Signal: "SIGKILL",
	}); err != nil {
		return err
	}
	i.containerId = ""
	i.resetPortMapping()
	return nil
}

func (i *Instance) waitToBeAvailable(ctx context.Context) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := i.isRunning(ctx); err != nil {
				return containerMissing{err: err}
			}
			if err := i.checkAvailable(ctx); err != nil {
				i.log.Debug(err)
			} else {
				i.log.Info("localstack: finished waiting")
				return nil
			}
		}
	}
}

func (i *Instance) isRunning(ctx context.Context) error {
	_, err := i.cli.ContainerInspect(ctx, i.getContainerId())
	if err != nil {
		i.log.Debug(err)
		return errors.New("localstack container has been stopped")
	}
	return nil
}

func (i *Instance) checkAvailable(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("local"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
		config.WithRetryer(func() aws.Retryer {
			return aws.NopRetryer{}
		}),
	)
	if err != nil {
		return err
	}

	resolver := NewDynamoDbResolverV2(i)
	s := dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolverV2(resolver))
	testTable := aws.String("bucket")

	if _, err = s.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []dynamotypes.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: dynamotypes.ScalarAttributeTypeS},
		},
		KeySchema: []dynamotypes.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: dynamotypes.KeyTypeHash},
		},
		ProvisionedThroughput: &dynamotypes.ProvisionedThroughput{
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

func (i *Instance) setContainerId(containerId string) {
	i.containerIdMutex.Lock()
	defer i.containerIdMutex.Unlock()
	i.containerId = containerId
}

func (i *Instance) getContainerId() string {
	i.containerIdMutex.RLock()
	defer i.containerIdMutex.RUnlock()
	return i.containerId
}

func (i *Instance) resetPortMapping() {
	i.savePortMappings(map[Service]string{})
}

func (i *Instance) savePortMappings(newMapping map[Service]string) {
	i.portMappingMutex.Lock()
	defer i.portMappingMutex.Unlock()
	i.portMapping = newMapping
}

func (i *Instance) getPortMapping(service Service) string {
	i.portMappingMutex.RLock()
	defer i.portMappingMutex.RUnlock()
	return i.portMapping[service]
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

func (i *Instance) writeContainerLogToLogger(ctx context.Context, containerId string) {
	reader, err := i.cli.ContainerLogs(ctx, containerId, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: true,
	})
	if err != nil {
		i.log.Error(err)
		return
	}
	defer logClose(reader)

	w := i.log.Writer()
	defer logClose(w)

	if _, err := io.Copy(w, reader); err != nil {
		if err := w.CloseWithError(err); err != nil {
			i.log.Println(err)
		}
	}
}

func logClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Println(err)
	}
}

type containerMissing struct {
	err error
}

func (c containerMissing) Error() string {
	return c.err.Error()
}
