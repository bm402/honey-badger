package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// struct for log entry to dynamodb
type RawLogEntry struct {
	IngressPort string `json:"ingress_port"`
	Timestamp   string `json:"timestamp"`
	IpAddress   string `json:"ip_address"`
	Input       string `json:"input"`
}

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
		log.Fatal("Error creating aws session")
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

// listens for tcp connections on the given port
func listen(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Error creating listener on port", port)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Error accepting connection on port", port)
		}
		conn.SetDeadline(time.Now().Add(15 * time.Minute))
		go handle(conn, port)
	}
}

// serves a false command prompt on the given tcp connection and reads input
func handle(conn net.Conn, port string) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		_, err := conn.Write([]byte("$ "))
		if err != nil {
			break
		}

		rawInputLen, err := conn.Read(buf)
		if err != nil {
			break
		}
		input := string(buf[:rawInputLen])
		writeInputsToRawLogsTable(conn, port, input)
	}
}

// writes inputs from the tcp connection to dynamodb table
func writeInputsToRawLogsTable(conn net.Conn, port string, input string) {
	rawLogEntry := RawLogEntry{
		IngressPort: port,
		Timestamp:   time.Now().String(),
		IpAddress:   conn.RemoteAddr().String(),
		Input:       input,
	}

	dynamoDocument, err := dynamodbattribute.MarshalMap(rawLogEntry)
	if err != nil {
		log.Fatal("Error creating dynamodb document")
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
