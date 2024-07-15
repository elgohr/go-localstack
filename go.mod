module github.com/elgohr/go-localstack

go 1.21

toolchain go1.22.5

require (
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/aws/aws-sdk-go v1.54.19
	github.com/aws/aws-sdk-go-v2 v1.30.3
	github.com/aws/aws-sdk-go-v2/config v1.27.26
	github.com/aws/aws-sdk-go-v2/credentials v1.17.26
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.53.3
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.40.3
	github.com/aws/aws-sdk-go-v2/service/cloudwatchevents v1.25.3
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.37.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.34.3
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.22.3
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.170.0
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.30.3
	github.com/aws/aws-sdk-go-v2/service/firehose v1.31.3
	github.com/aws/aws-sdk-go-v2/service/iam v1.34.3
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.29.3
	github.com/aws/aws-sdk-go-v2/service/lambda v1.56.3
	github.com/aws/aws-sdk-go-v2/service/redshift v1.46.4
	github.com/aws/aws-sdk-go-v2/service/route53 v1.42.3
	github.com/aws/aws-sdk-go-v2/service/s3 v1.58.2
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.32.3
	github.com/aws/aws-sdk-go-v2/service/ses v1.25.2
	github.com/aws/aws-sdk-go-v2/service/sfn v1.29.3
	github.com/aws/aws-sdk-go-v2/service/sns v1.31.3
	github.com/aws/aws-sdk-go-v2/service/sqs v1.34.3
	github.com/aws/aws-sdk-go-v2/service/ssm v1.52.3
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.3
	github.com/aws/smithy-go v1.20.3
	github.com/docker/docker v27.0.3+incompatible
	github.com/docker/go-connections v0.5.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.8.1
	github.com/moby/moby v27.0.3+incompatible
	github.com/opencontainers/image-spec v1.1.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.3 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.4 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.53.0 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.22.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk v1.22.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/tools v0.23.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.0.3 // indirect
)
