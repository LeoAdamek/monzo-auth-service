package main

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kms"
	"os"
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"context"
	"net/http"
	"time"
	"encoding/json"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var db *dynamodb.DynamoDB

var tableName, oAuthClientSecret string

func init() {
	s := session.Must(session.NewSession())
	db = dynamodb.New(s)
	tableName = os.Getenv("TABLE_NAME")
	
	k := kms.New(s)
	
	xray.Configure(xray.Config{})
	
	// TODO: Figure out why the tracing isn't working
	// xray.AWS(k.Client)
	// xray.AWS(db.Client)
	
	
	encryptedSecret := os.Getenv("OAuthClientSecret")
	
	blob, err := base64.StdEncoding.DecodeString(encryptedSecret)
	
	if err != nil {
		panic("Invalid B64 value for OAuthClientSecret")
	}
	
	resp, err := k.DecryptWithContext(context.Background(), &kms.DecryptInput{CiphertextBlob: blob})
	
	if err != nil {
		panic("Could not decrypt OAuthClientSecret")
	}
	
	oAuthClientSecret = string(resp.Plaintext)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	
	r := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string {
			"Content-Type": "application/json",
		},
	}
	
	id := req.QueryStringParameters["id"]
	
	if id == "" {
		r.StatusCode = http.StatusBadRequest
		r.Body = `{"error": "No token ID supplied"}`
		return r, nil
	}
	
	token, err := getToken(ctx, id)

	if err != nil {
		r.StatusCode = http.StatusNotFound
		r.Body = `{"error": "No such token ID"}`
		return r, nil
	}
	
	// AuthCode isn't yet set/available.
	// Return "No Content" to let the caller know they should try again.
	if token.AuthCode == "" {
		r.StatusCode = http.StatusNoContent
		r.Body = ""
		return r, nil
	}

	t, err := exchangeCode(ctx, token.AuthCode)
	
	if err != nil {
		r.StatusCode = http.StatusBadGateway
		return r, nil
	}
	
	expiresAt := time.Now().Add(time.Duration(t.Lifetime) * time.Second)
	
	response := map[string]string {
		"token": t.Token,
		"user_id": t.UserID,
		"expires_at": expiresAt.Format(time.RFC3339),
	}
	
	body, _ := json.Marshal(response)
	
	r.Body = string(body)
	
	
	return r, nil
	
}

func main() {
	lambda.Start(handler)
}