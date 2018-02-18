package main

import (
	"context"
	"github.com/LeoAdamek/monzo-auth-service/models"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"net/url"
	"os"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/aws/aws-xray-sdk-go/xray"
)

type accessToken struct {
	Token string `json:"access_token"`
	UserID string `json:"user_id"`
	Lifetime int `json:"expires_in"`
}

func getToken(ctx context.Context, id string) (models.Token, error) {
	
	req := &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue {
			"id": {
				S: &id,
			},
		},
	}
	
	resp, err := db.GetItemWithContext(ctx, req)
	
	if err != nil {
		return models.Token{}, err
	}
	
	var t models.Token
	
	err = dynamodbattribute.UnmarshalMap(resp.Item, &t)
	
	return t, err
}


func exchangeCode(ctx context.Context, authCode string) (accessToken, error) {
	var t accessToken
	
	// AuthCode is set, exchange it for an access token and return it.
	baseUrl, _ := url.Parse("https://api.monzo.com/oauth2/token")
	query := make(url.Values)
	query.Set("grant_type","authorization_code")
	query.Set("client_id", os.Getenv("OAuthClientID"))
	query.Set("client_secret", oAuthClientSecret)
	query.Set("code", authCode)
	query.Set("redirect_uri", os.Getenv("RETURN_URL"))
	baseUrl.RawQuery = query.Encode()
	
	hc := xray.Client(http.DefaultClient)
	
	resp, err := hc.Get(baseUrl.String())
	
	if err != nil {
		return t, err
	}
	
	data, err := ioutil.ReadAll(resp.Body)
	
	if err != nil {
		return t, err
	}
	
	err = json.Unmarshal(data, &t)
	
	return t, err
}
