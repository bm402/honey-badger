package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// struct for the heatmap data response
type HeatmapDataResponse struct {
	HeatmapDataPoints []HeatmapDataPoint `json:"heatmap_data_points"`
}

type HeatmapDataPoint struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Count int     `json:"count"`
}

// retrieves heatmap data points from the aggregated logs table
func getHeatmapData() (HeatmapDataResponse, error) {

	projection := expression.NamesList(expression.Name("lat"),
		expression.Name("lon"), expression.Name("count"))

	expression, err := expression.NewBuilder().WithProjection(projection).Build()
	if err != nil {
		return HeatmapDataResponse{}, err
	}

	scanInput := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		ProjectionExpression:      expression.Projection(),
		TableName:                 aws.String(aggregatedLogsTableName),
	}

	scanOutput, err := dynamoClient.Scan(scanInput)
	if err != nil {
		return HeatmapDataResponse{}, err
	}

	var heatmapDataPoints []HeatmapDataPoint
	err = dynamodbattribute.UnmarshalListOfMaps(scanOutput.Items, &heatmapDataPoints)
	if err != nil {
		return HeatmapDataResponse{}, err
	}

	return HeatmapDataResponse{
		HeatmapDataPoints: heatmapDataPoints,
	}, nil
}
