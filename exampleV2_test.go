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

//nolint:all
package localstack_test

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/elgohr/go-localstack"
)

func ExampleInstance_EndpointV2() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	l, err := localstack.NewInstance()
	if err != nil {
		log.Fatalf("Could not connect to Docker %v", err)
	}
	if err := l.Start(); err != nil {
		log.Fatalf("Could not start localstack %v", err)
	}
	defer func() { // this should be t.Cleanup for better stability
		if err := l.Stop(); err != nil {
			log.Fatalf("Could not stop localstack %v", err)
		}
	}()

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

func ExampleInstance_withEndpointResolverV2() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	l, err := localstack.NewInstance()
	if err != nil {
		log.Fatalf("Could not connect to Docker %v", err)
	}
	if err := l.Start(); err != nil {
		log.Fatalf("Could not start localstack %v", err)
	}
	defer func() { // this should be t.Cleanup for better stability
		if err := l.Stop(); err != nil {
			log.Fatalf("Could not stop localstack %v", err)
		}
	}()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
	)
	if err != nil {
		log.Fatalf("Could not get config %v", err)
	}

	resolver := localstack.NewDynamoDbResolverV2(l)
	client := dynamodb.NewFromConfig(cfg, dynamodb.WithEndpointResolverV2(resolver))

	myTestWithV2Client(client)
}

func myTestWithV2(_ aws.Config) {}

func myTestWithV2Client(_ *dynamodb.Client) {}
