package main

import (
	"github.com/aws/aws-lambda-go/events"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/LeoAdamek/monzo-auth-service/models"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"net/http"
	"os"
	"encoding/json"
	"net/url"
)

var db *dynamodb.DynamoDB
var tableName string

func init() {
	s := session.Must(session.NewSession())
	db = dynamodb.New(s)
	
	xray.Configure(xray.Config{
		LogLevel: "trace",
	})

	// TODO: Work out why the tracing doesn't work.
	//xray.AWS(db.Client)
	
	tableName = os.Getenv("TABLE_NAME")
}

func handler(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	r := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	
	t := models.NewToken()
	
	returnUrl := os.Getenv("RETURN_URL")
	redirectQuery := make(url.Values)
	
	redirectQuery.Set("client_id", os.Getenv("OAuthClientID"))
	redirectQuery.Set("response_type", "code")
	redirectQuery.Set("state", t.ID)
	redirectQuery.Set("redirect_uri", returnUrl)
	
	t.LoginURL = "https://auth.monzo.com/?" + redirectQuery.Encode()
	
	values, err := dynamodbattribute.MarshalMap(t)
	
	if err != nil {
		r.StatusCode = http.StatusInternalServerError
		r.Body = `{"error": "Unable to create token"}`
		return r, err
	}
	
	req := &dynamodb.PutItemInput{
		TableName: &tableName,
		Item: values,
	}
	
	_, err = db.PutItemWithContext(ctx, req)
	
	if err != nil {
		r.StatusCode = http.StatusInternalServerError
		r.Body = `{"error": "Unable to create token"}`
		return r, err
	}
	
	b, _ := json.Marshal(t)
	
	r.Body = string(b)
	
	return r, nil
}

func main() {
	lambda.Start(handler)
}
