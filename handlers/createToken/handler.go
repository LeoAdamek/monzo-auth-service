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
	"github.com/aws/aws-sdk-go/service/kms"
	"encoding/base64"
	"net/url"
	"github.com/averagesecurityguy/random"
)

var db *dynamodb.DynamoDB
var kmsClient *kms.KMS
var tableName string

// OAuth client secret, will be decrypted and stored on first load.
var clientSecret string

func init() {
	s := session.Must(session.NewSession())
	db = dynamodb.New(s)
	kmsClient = kms.New(session.Must(session.NewSession()))
	
	xray.Configure(xray.Config{
		LogLevel: "trace",
	})
	// xray.AWS(db.Client)
	
	tableName = os.Getenv("TABLE_NAME")
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	
	r := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	
	if clientSecret == "" {
		
		
		encryptedSecret := os.Getenv("OAuthClientSecret")
		
		blob, err := base64.StdEncoding.DecodeString(encryptedSecret)
		
		if err != nil {
			r.StatusCode = http.StatusInternalServerError
			r.Body = `{"error": "Invalid OAuthClientSecret value"}`
			return r, nil
		}
		
		req := &kms.DecryptInput{
			CiphertextBlob: blob,
		}

		res, err := kmsClient.DecryptWithContext(ctx, req)
		
		if err != nil {
			r.StatusCode = http.StatusInternalServerError
			r.Body = `{"error": "Unable to load client config"}`
			return r, err
		}
		
		clientSecret = string(res.Plaintext)
	}
	
	stateToken, err := random.AlphaNum(32)
	
	if err != nil {
		r.StatusCode = http.StatusInternalServerError
		r.Body = `{"error": "Unable to generate secure token"}`
		return r, nil
	}
	t := models.NewToken()
	
	redirectQuery := make(url.Values)
	
	redirectQuery.Set("client_id", os.Getenv("OAuthClientID"))
	redirectQuery.Set("response_type", "code")
	redirectQuery.Set("state", stateToken)
	
	returnUrl, err := url.Parse(os.Getenv("RETURN_URL"))
	returnQuery := returnUrl.Query()
	returnQuery.Set("sid", t.ID)
	returnUrl.RawQuery = returnQuery.Encode()
	
	redirectQuery.Set("redirect_uri", returnUrl.String())
	
	
	t.StateToken = stateToken
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
