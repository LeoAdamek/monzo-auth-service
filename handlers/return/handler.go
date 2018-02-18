package main

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"github.com/aws/aws-lambda-go/lambda"
)

var db *dynamodb.DynamoDB

var tableName string

func init() {
	s := session.Must(session.NewSession())
	db = dynamodb.New(s)
	tableName = os.Getenv("TABLE_NAME")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}

	tokenId := req.QueryStringParameters["state"]
	
	if tokenId == "" {
		r.Body = templateFail
		r.StatusCode = http.StatusNotFound
		return r, nil
	}
	
	// Store the authorization code in DB — We will exchange it for an access token when the
	// Original app requests it.
	
	authCode := req.QueryStringParameters["code"]
	err := setAuthCode(ctx, tokenId, authCode)
	
	if err != nil {
		r.StatusCode = http.StatusNotFound
		
		r.Body = templateFail
		return r, nil
	}
	
	r.Body = templateSuccess
	return r, nil
}

func main() {
	lambda.Start(handler)
}
