package main

import (
	"log"
	"net"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// struct for log entry to dynamodb
type RawLogEntry struct {
	IngressPort string `json:"ingress_port"`
	Timestamp   int64  `json:"timestamp"`
	IpAddress   string `json:"ip_address"`
	Input       string `json:"input"`
}

// writes inputs from the tcp connection to dynamodb table
func writeInputsToRawLogsTable(conn net.Conn, port string, input string) {

	// timestamp in epoch millis
	timestamp := int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)

	rawLogEntry := RawLogEntry{
		IngressPort: port,
		Timestamp:   timestamp,
		IpAddress:   conn.RemoteAddr().String(),
		Input:       input,
	}

	dynamoDocument, err := dynamodbattribute.MarshalMap(rawLogEntry)
	if err != nil {
		log.Fatal("Error creating dynamodb document:", err.Error())
	}

	log.Println(dynamoDocument)

	putItemInput := &dynamodb.PutItemInput{
		Item:      dynamoDocument,
		TableName: aws.String(rawLogsTableName),
	}
	_, err = dynamoClient.PutItem(putItemInput)
	if err != nil {
		log.Fatal("Error writing log entry to dynamodb:", err.Error())
	}
}
