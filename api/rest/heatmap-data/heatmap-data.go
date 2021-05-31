package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// struct for the api gateway http response
type HeatmapDataResponse struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

// struct for the heatmap data points
type HeatmapDataPoints struct {
	HeatmapDataPoints []HeatmapDataPoint `json:"heatmap_data_points"`
}

// struct for a single heatmap data point
type HeatmapDataPoint struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Count int     `json:"count"`
}

// packages heatmap data points into a http response format for the api gateway
func getHeatmapDataResponse() (HeatmapDataResponse, error) {

	// http headers including cors
	headers := map[string]string{
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "OPTIONS,GET",
		"Content-Type":                 "application/json",
	}

	// function for creating an error response
	createErrorResponse := func(err error) HeatmapDataResponse {
		return HeatmapDataResponse{
			StatusCode:      400,
			Headers:         headers,
			Body:            "{\"error\":\"" + err.Error() + "\"}",
			IsBase64Encoded: false,
		}
	}

	heatmapDataPoints, err := getHeatmapDataPoints()
	if err != nil {
		return createErrorResponse(err), nil
	}

	body, err := json.Marshal(heatmapDataPoints)
	if err != nil {
		return createErrorResponse(err), nil
	}

	return HeatmapDataResponse{
		StatusCode:      200,
		Headers:         headers,
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}

// retrieves heatmap data points from the aggregated logs table
func getHeatmapDataPoints() (HeatmapDataPoints, error) {

	projection := expression.NamesList(expression.Name("lat"),
		expression.Name("lon"), expression.Name("count"))

	expression, err := expression.NewBuilder().WithProjection(projection).Build()
	if err != nil {
		return HeatmapDataPoints{}, err
	}

	scanInput := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		ProjectionExpression:      expression.Projection(),
		TableName:                 aws.String(aggregatedLogsTableName),
	}

	scanOutput, err := dynamoClient.Scan(scanInput)
	if err != nil {
		return HeatmapDataPoints{}, err
	}

	var heatmapDataPoints []HeatmapDataPoint
	err = dynamodbattribute.UnmarshalListOfMaps(scanOutput.Items, &heatmapDataPoints)
	if err != nil {
		return HeatmapDataPoints{}, err
	}

	return HeatmapDataPoints{
		HeatmapDataPoints: heatmapDataPoints,
	}, nil
}
