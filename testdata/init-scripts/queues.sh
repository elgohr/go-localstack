#!bin/bash

echo "########### Creating profile ###########"

aws configure set aws_access_key_id default_access_key --profile=localstack
aws configure set aws_secret_access_key default_secret_key --profile=localstack
aws configure set region eu-west-1 --profile=localstack

echo "########### Listing profile ###########"
aws configure list --profile=localstack


echo "########### Creating SQS ###########"
aws sqs create-queue --endpoint-url=http://localhost:4566 --queue-name=local_queue_my --profile=localstack
aws sqs create-queue --endpoint-url=http://localhost:4566 --queue-name=local_queue_my_1 --profile=localstack

echo "########### Bootstrap Complete ###########"