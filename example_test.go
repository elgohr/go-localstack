package localstack_test

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/elgohr/go-localstack"
	"log"
	"time"
)

func ExampleLocalstackWithContext() {
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

func ExampleLocalstack() {
	l, err := localstack.NewInstance()
	if err != nil {
		log.Fatalf("Could not connect to Docker %v", err)
	}
	if err := l.Start(); err != nil {
		log.Fatalf("Could not start localstack %v", err)
	}

	myTestWith(&aws.Config{
		Credentials: credentials.NewStaticCredentials("not", "empty", ""),
		DisableSSL:  aws.Bool(true),
		Region:      aws.String(endpoints.UsWest1RegionID),
		Endpoint:    aws.String(l.Endpoint(localstack.SQS)),
	})

	if err := l.Stop(); err != nil {
		log.Fatalf("Could not stop localstack %v", err)
	}
}

func myTestWith(_ *aws.Config) {}
