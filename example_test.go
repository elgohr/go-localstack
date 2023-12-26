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
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"

	"github.com/elgohr/go-localstack"
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
