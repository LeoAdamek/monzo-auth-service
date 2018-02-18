package main

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"context"
	"github.com/aws/aws-sdk-go/aws"
)


func setAuthCode(ctx context.Context, tokenID string, authCode string) error {
	req := &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": { S: &tokenID },
		},
		UpdateExpression: aws.String("SET auth_code = :authcode"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":authcode": {S: &authCode},
		},
	}
	
	_, err := db.UpdateItemWithContext(ctx, req)
	
	return err
}

