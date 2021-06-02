package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	apigwClient          *apigatewaymanagementapi.ApiGatewayManagementApi
	dynamoClient         dynamodbiface.DynamoDBAPI
	connectionsTableName string
)

// initialises connection to dynamodb
func init() {
	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("API_GATEWAY_ENDPOINT")

	dynamoSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatal("Error creating aws session for dynamodb: ", err.Error())
	}
	dynamoClient = dynamodb.New(dynamoSession)

	apigwSession, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(endpoint),
		Region:   aws.String(region),
	})
	if err != nil {
		log.Fatal("Error creating aws session for apigw: ", err.Error())
	}
	apigwClient = apigatewaymanagementapi.New(apigwSession)

	connectionsTableName = os.Getenv("CONNECTIONS_TABLE_NAME")
}

// lambda function handler
func handler(event events.DynamoDBEvent) error {
	return broadcastDynamoDBEvent(event)
}

// entry point
func main() {
	lambda.Start(handler)
}
