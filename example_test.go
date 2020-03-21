package localstack_test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/elgohr/go-localstack"
	"log"
)

func ExampleLocalstack() {
	l, err := localstack.NewInstance()
	if err != nil {
		log.Fatalf("Could not connect to Docker %v", err)
	}
	if err := l.Start(); err != nil {
		log.Fatalf("Could not start localstack %v", err)
	}

	session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("not", "empty", ""),
		DisableSSL:  aws.Bool(true),
		Region:      aws.String(endpoints.UsWest1RegionID),
		Endpoint:    aws.String(l.Endpoint(localstack.SQS)),
	})

	myTest()

	
	if err := l.Stop(); err != nil {
		log.Fatalf("Could not stop localstack %v", err)
	}
}

func myTest() {}
