// Code generated by smithy-go-codegen DO NOT EDIT.

package lambda

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Returns your function's [recursive loop detection] configuration.
//
// [recursive loop detection]: https://docs.aws.amazon.com/lambda/latest/dg/invocation-recursion.html
func (c *Client) GetFunctionRecursionConfig(ctx context.Context, params *GetFunctionRecursionConfigInput, optFns ...func(*Options)) (*GetFunctionRecursionConfigOutput, error) {
	if params == nil {
		params = &GetFunctionRecursionConfigInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "GetFunctionRecursionConfig", params, optFns, c.addOperationGetFunctionRecursionConfigMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*GetFunctionRecursionConfigOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type GetFunctionRecursionConfigInput struct {

	//
	//
	// This member is required.
	FunctionName *string

	noSmithyDocumentSerde
}

type GetFunctionRecursionConfigOutput struct {

	// If your function's recursive loop detection configuration is Allow , Lambda
	// doesn't take any action when it detects your function being invoked as part of a
	// recursive loop.
	//
	// If your function's recursive loop detection configuration is Terminate , Lambda
	// stops your function being invoked and notifies you when it detects your function
	// being invoked as part of a recursive loop.
	//
	// By default, Lambda sets your function's configuration to Terminate . You can
	// update this configuration using the PutFunctionRecursionConfigaction.
	RecursiveLoop types.RecursiveLoop

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationGetFunctionRecursionConfigMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestjson1_serializeOpGetFunctionRecursionConfig{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpGetFunctionRecursionConfig{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "GetFunctionRecursionConfig"); err != nil {
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
	if err = addOpGetFunctionRecursionConfigValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opGetFunctionRecursionConfig(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opGetFunctionRecursionConfig(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "GetFunctionRecursionConfig",
	}
}