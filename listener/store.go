package main

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// struct for log entry to dynamodb
type RawLogEntry struct {
	IngressPort string  `json:"ingress_port"`
	Timestamp   int64   `json:"timestamp"`
	IpAddress   string  `json:"ip_address"`
	Location    string  `json:"location"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Input       string  `json:"input"`
}

// writes inputs from the tcp connection to dynamodb table
func writeInputsToRawLogsTable(conn net.Conn, port string, input string) {

	// timestamp in epoch millis
	timestamp := int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)

	remoteAddress := conn.RemoteAddr().String()
	ipAddress := strings.Split(remoteAddress, ":")[0]

	ipLocationData := getIpLocationData(ipAddress)

	rawLogEntry := RawLogEntry{
		IngressPort: port,
		Timestamp:   timestamp,
		IpAddress:   ipAddress,
		Location:    ipLocationData.Location,
		Lat:         ipLocationData.Lat,
		Lon:         ipLocationData.Lon,
		Input:       input,
	}

	err := writeRawLogEntryToDb(rawLogEntry)
	if err != nil {
		log.Println("Error writing log entry to dynamodb: ", err.Error())
		log.Println("Retrying with sanitised input")

		rawLogEntry.Input = "<binary>"
		err = writeRawLogEntryToDb(rawLogEntry)
		if err != nil {
			log.Fatal("Error writing log entry to dynamodb: ", err.Error())
		}
	}
}

// sends data to dynamodb
func writeRawLogEntryToDb(rawLogEntry RawLogEntry) error {
	dynamoDocument, err := dynamodbattribute.MarshalMap(rawLogEntry)
	if err != nil {
		log.Fatal("Error creating dynamodb document: ", err.Error())
	}

	putItemInput := &dynamodb.PutItemInput{
		Item:      dynamoDocument,
		TableName: aws.String(rawLogsTableName),
	}
	_, err = dynamoClient.PutItem(putItemInput)
	return err
}
