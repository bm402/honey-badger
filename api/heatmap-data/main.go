package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynamoClient            dynamodbiface.DynamoDBAPI
	aggregatedLogsTableName string
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
	aggregatedLogsTableName = os.Getenv("AGGREGATED_LOGS_TABLE_NAME")
}

// lambda function handler
func handler() (HeatmapDataResponse, error) {
	return getHeatmapDataResponse()
}

// entry point
func main() {
	lambda.Start(handler)
}
