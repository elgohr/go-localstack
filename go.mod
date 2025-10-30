module github.com/elgohr/go-localstack

go 1.24.0

require (
	github.com/Masterminds/semver/v3 v3.4.0
	github.com/aws/aws-sdk-go v1.55.5
	github.com/aws/aws-sdk-go-v2 v1.30.5
	github.com/aws/aws-sdk-go-v2/config v1.27.33
	github.com/aws/aws-sdk-go-v2/credentials v1.17.32
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.53.7
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.40.7
	github.com/aws/aws-sdk-go-v2/service/cloudwatchevents v1.25.7
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.39.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.34.9
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.22.7
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.177.3
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.30.7
	github.com/aws/aws-sdk-go-v2/service/firehose v1.32.4
	github.com/aws/aws-sdk-go-v2/service/iam v1.35.2
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.29.7
	github.com/aws/aws-sdk-go-v2/service/lambda v1.58.3
	github.com/aws/aws-sdk-go-v2/service/redshift v1.46.8
	github.com/aws/aws-sdk-go-v2/service/route53 v1.43.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.61.2
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.32.8
	github.com/aws/aws-sdk-go-v2/service/ses v1.26.2
	github.com/aws/aws-sdk-go-v2/service/sfn v1.31.2
	github.com/aws/aws-sdk-go-v2/service/sns v1.31.7
	github.com/aws/aws-sdk-go-v2/service/sqs v1.34.8
	github.com/aws/aws-sdk-go-v2/service/ssm v1.52.8
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.7
	github.com/aws/smithy-go v1.20.4
	github.com/docker/docker v28.5.1+incompatible
	github.com/docker/go-connections v0.6.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.12.0
	github.com/opencontainers/image-spec v1.1.1
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.4 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.7 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/sys/atomicwriter v0.1.0 // indirect
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.54.0 // indirect
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.22.0 // indirect
	go.opentelemetry.io/otel/metric v1.29.0 // indirect
	go.opentelemetry.io/otel/sdk v1.22.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/tools v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.0.3 // indirect
)

tool github.com/maxbrunsfeld/counterfeiter/v6
