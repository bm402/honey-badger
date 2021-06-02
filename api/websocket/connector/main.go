package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynamoClient         dynamodbiface.DynamoDBAPI
	connectionsTableName string
)

// initialises connection to dynamodb
func init() {
	region := os.Getenv("AWS_REGION")
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatal("Error creating aws session: ", err.Error())
	}

	dynamoClient = dynamodb.New(session)
	connectionsTableName = os.Getenv("CONNECTIONS_TABLE_NAME")
}

// lambda function handler
func handler(event events.APIGatewayWebsocketProxyRequest) (ConnectorResponse, error) {
	return handleNewConnectionWithResponse(event)
}

// entry point
func main() {
	lambda.Start(handler)
}
