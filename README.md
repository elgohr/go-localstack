# go-localstack

[![Test](https://github.com/elgohr/go-localstack/workflows/Test/badge.svg)](https://github.com/elgohr/go-localstack/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/elgohr/go-localstack/branch/main/graph/badge.svg)](https://codecov.io/gh/elgohr/go-localstack)
[![CodeQL](https://github.com/elgohr/go-localstack/workflows/CodeQL/badge.svg)](https://github.com/elgohr/go-localstack/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/elgohr/go-localstack)](https://goreportcard.com/report/github.com/elgohr/go-localstack)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/elgohr/go-localstack)](https://pkg.go.dev/github.com/elgohr/go-localstack)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gojp/goreportcard/blob/master/LICENSE)

Go Wrapper for using [localstack](https://github.com/localstack/localstack) in go testing

# Installation

Please make sure that you have Docker installed.

```bash
go get github.com/elgohr/go-localstack
```

# Usage

With SDK V2
```go
func ExampleLocalstackWithContextSdkV2() {
    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()
    
    l, err := localstack.NewInstance()
    if err != nil {
        log.Fatalf("Could not connect to Docker %v", err)
    }
    if err := l.StartWithContext(ctx); err != nil {
        log.Fatalf("Could not start localstack %v", err)
    }
    
    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion("us-east-1"),
        config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(_, _ string, _ ...interface{}) (aws.Endpoint, error) {
            return aws.Endpoint{
			    PartitionID:       "aws", 
			    URL:               l.EndpointV2(localstack.SQS), 
			    SigningRegion:     "us-east-1", 
			    HostnameImmutable: true,
		    }, nil
        })),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
    )
    if err != nil {
        log.Fatalf("Could not get config %v", err)
    }
    
    myTestWithV2(cfg)
}
```

With SDK V1
```go
func TestWithLocalStack(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()

    l, err := localstack.NewInstance()
    if err != nil {
        log.Fatalf("Could not connect to Docker %v", err)
    }
    if err := l.StartWithContext(ctx); err != nil {
        log.Fatalf("Could not start localstack %v", err)
    }

    myTestWith(&aws.Config{
        Credentials: credentials.NewStaticCredentials("not", "empty", ""),
        DisableSSL:  aws.Bool(true),
        Region:      aws.String(endpoints.UsWest1RegionID),
        Endpoint:    aws.String(l.Endpoint(localstack.SQS)),
    })
}
```

