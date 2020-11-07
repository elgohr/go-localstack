# go-localstack
[![Actions Status](https://github.com/elgohr/go-localstack/workflows/Test/badge.svg)](https://github.com/elgohr/go-localstack/actions)
[![codecov](https://codecov.io/gh/elgohr/go-localstack/branch/master/graph/badge.svg)](https://codecov.io/gh/elgohr/go-localstack)
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

```go
func TestWithLocalStack(t *testing.T) {
	l, err := localstack.NewInstance()
	if err != nil {
		t.Fatalf("Could not connect to Docker %v", err)
	}
	if err := l.Start(); err != nil {
		t.Fatalf("Could not start localstack %v", err)
	}

	session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("not", "empty", ""),
		DisableSSL:  aws.Bool(true),
		Region:      aws.String(endpoints.UsWest1RegionID),
		Endpoint:    aws.String(l.Endpoint(localstack.SQS)),
	})

	myTest()

	if err := l.Stop(); err != nil {
		t.Fatalf("Could not stop localstack %v", err)
	}
}
```
