package main

import (
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var (
	dynamoClient     dynamodbiface.DynamoDBAPI
	rawLogsTableName string
)

// gets parameters and sets up dynamodb session
func init() {
	region := os.Getenv("AWS_REGION")

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatal("Error creating aws session:", err.Error())
	}

	ssmClient := ssm.New(session)
	getParameterOutput, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: aws.String("RawLogsTableName"),
	})
	if err != nil {
		log.Fatal("Error getting raw logs table name from ssm parameter store:", err.Error())
	}

	dynamoClient = dynamodb.New(session)
	rawLogsTableName = aws.StringValue(getParameterOutput.Parameter.Value)
}

func main() {
	port := flag.String("p", "8081", "Port to listen on")
	flag.Parse()

	listen(*port)
}
