package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aria-music/aria-alexa/alexa"
)

func main() {
	lambda.Start(alexaHandler)
}

// alexaHandler is an entrypoint of AWS Lambda.
// This always returns err == nil. Func signature is just for Lambda.
func alexaHandler(ctx context.Context, request *alexa.Request) (*alexa.Response, error) {
	if resp := handleRequest(ctx, request); resp != nil {
		return resp, nil
	}
	return alexa.EmptyResponse, nil
}
