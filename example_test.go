package localstack_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/elgohr/go-localstack"
	"testing"
)

func TestExampleLocalstack(t *testing.T) {
	l := &localstack.Instance{}
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

func myTest() {}
