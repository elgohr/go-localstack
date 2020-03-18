# go-localstack
[![Actions Status](https://github.com/elgohr/go-localstack/workflows/Test/badge.svg)](https://github.com/elgohr/go-localstack/actions)

Golang Wrapper for using [localstack](https://github.com/localstack/localstack) in go testing

# Installation
```bash
go get github.com/elgohr/go-localstack
```

# Usage

```go
func TestWithLocalStack(t *testing.T) {
    l := &localstack.Instance{}
    if err := l.Start(); err != nil {
        t.Fatalf("Could not start localstack %v", err)
    }
    sqsSess, err := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials("not", "empty", ""),
        DisableSSL:  aws.Bool(true),
        Region:      aws.String(endpoints.UsWest1RegionID),
        Endpoint:    aws.String(l.Endpoint(localstack.SQS)),
    })
    
    myTestWith(sqsSess)

    if err := l.Stop(); err != nil {
        t.Fatalf("Could not stop localstack %v", err)
    }
}
```
