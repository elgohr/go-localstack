// Code generated by smithy-go-codegen DO NOT EDIT.

package kinesis

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/ptr"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Writes a single data record into an Amazon Kinesis data stream. Call PutRecord
// to send data into the stream for real-time ingestion and subsequent processing,
// one record at a time. Each shard can support writes up to 1,000 records per
// second, up to a maximum data write total of 1 MiB per second.
//
// When invoking this API, you must use either the StreamARN or the StreamName
// parameter, or both. It is recommended that you use the StreamARN input
// parameter when you invoke this API.
//
// You must specify the name of the stream that captures, stores, and transports
// the data; a partition key; and the data blob itself.
//
// The data blob can be any type of data; for example, a segment from a log file,
// geographic/location data, website clickstream data, and so on.
//
// The partition key is used by Kinesis Data Streams to distribute data across
// shards. Kinesis Data Streams segregates the data records that belong to a stream
// into multiple shards, using the partition key associated with each data record
// to determine the shard to which a given data record belongs.
//
// Partition keys are Unicode strings, with a maximum length limit of 256
// characters for each key. An MD5 hash function is used to map partition keys to
// 128-bit integer values and to map associated data records to shards using the
// hash key ranges of the shards. You can override hashing the partition key to
// determine the shard by explicitly specifying a hash value using the
// ExplicitHashKey parameter. For more information, see [Adding Data to a Stream] in the Amazon Kinesis
// Data Streams Developer Guide.
//
// PutRecord returns the shard ID of where the data record was placed and the
// sequence number that was assigned to the data record.
//
// Sequence numbers increase over time and are specific to a shard within a
// stream, not across all shards within a stream. To guarantee strictly increasing
// ordering, write serially to a shard and use the SequenceNumberForOrdering
// parameter. For more information, see [Adding Data to a Stream]in the Amazon Kinesis Data Streams
// Developer Guide.
//
// After you write a record to a stream, you cannot modify that record or its
// order within the stream.
//
// If a PutRecord request cannot be processed because of insufficient provisioned
// throughput on the shard involved in the request, PutRecord throws
// ProvisionedThroughputExceededException .
//
// By default, data records are accessible for 24 hours from the time that they
// are added to a stream. You can use IncreaseStreamRetentionPeriodor DecreaseStreamRetentionPeriod to modify this retention period.
//
// [Adding Data to a Stream]: https://docs.aws.amazon.com/kinesis/latest/dev/developing-producers-with-sdk.html#kinesis-using-sdk-java-add-data-to-stream
func (c *Client) PutRecord(ctx context.Context, params *PutRecordInput, optFns ...func(*Options)) (*PutRecordOutput, error) {
	if params == nil {
		params = &PutRecordInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "PutRecord", params, optFns, c.addOperationPutRecordMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*PutRecordOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// Represents the input for PutRecord .
type PutRecordInput struct {

	// The data blob to put into the record, which is base64-encoded when the blob is
	// serialized. When the data blob (the payload before base64-encoding) is added to
	// the partition key size, the total size must not exceed the maximum record size
	// (1 MiB).
	//
	// This member is required.
	Data []byte

	// Determines which shard in the stream the data record is assigned to. Partition
	// keys are Unicode strings with a maximum length limit of 256 characters for each
	// key. Amazon Kinesis Data Streams uses the partition key as input to a hash
	// function that maps the partition key and associated data to a specific shard.
	// Specifically, an MD5 hash function is used to map partition keys to 128-bit
	// integer values and to map associated data records to shards. As a result of this
	// hashing mechanism, all data records with the same partition key map to the same
	// shard within the stream.
	//
	// This member is required.
	PartitionKey *string

	// The hash value used to explicitly determine the shard the data record is
	// assigned to by overriding the partition key hash.
	ExplicitHashKey *string

	// Guarantees strictly increasing sequence numbers, for puts from the same client
	// and to the same partition key. Usage: set the SequenceNumberForOrdering of
	// record n to the sequence number of record n-1 (as returned in the result when
	// putting record n-1). If this parameter is not set, records are coarsely ordered
	// based on arrival time.
	SequenceNumberForOrdering *string

	// The ARN of the stream.
	StreamARN *string

	// The name of the stream to put the data record into.
	StreamName *string

	noSmithyDocumentSerde
}

func (in *PutRecordInput) bindEndpointParams(p *EndpointParameters) {
	p.StreamARN = in.StreamARN
	p.OperationType = ptr.String("data")
}

// Represents the output for PutRecord .
type PutRecordOutput struct {

	// The sequence number identifier that was assigned to the put data record. The
	// sequence number for the record is unique across all records in the stream. A
	// sequence number is the identifier associated with every record put into the
	// stream.
	//
	// This member is required.
	SequenceNumber *string

	// The shard ID of the shard where the data record was placed.
	//
	// This member is required.
	ShardId *string

	// The encryption type to use on the record. This parameter can be one of the
	// following values:
	//
	//   - NONE : Do not encrypt the records in the stream.
	//
	//   - KMS : Use server-side encryption on the records in the stream using a
	//   customer-managed Amazon Web Services KMS key.
	EncryptionType types.EncryptionType

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationPutRecordMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpPutRecord{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpPutRecord{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "PutRecord"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addOpPutRecordValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opPutRecord(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opPutRecord(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "PutRecord",
	}
}