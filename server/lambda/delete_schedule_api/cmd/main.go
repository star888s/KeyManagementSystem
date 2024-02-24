package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type Response events.APIGatewayProxyResponse

type Body struct {
    ID string `json:"id"`
    StartTime string `json:"startTime"`
}

type Bodys[]Body

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {

    var bodys Bodys

    corsHeaders := map[string]string{
        "Content-Type": "application/json",
        "Access-Control-Allow-Origin":  request.Headers["origin"],
        "Access-Control-Allow-Methods": "POST,OPTIONS",
        "Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
    }

    err := json.Unmarshal([]byte(request.Body), &bodys)
    if err != nil {
        fmt.Println("Could not decode body", err)
        return Response{
            Headers: corsHeaders,
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
    }

    slog.Info("Received body: ", "%s",bodys)
    err = validateBody(bodys)
    if err != nil {
        return Response{
            Headers: corsHeaders,
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
    }

    err = deleteItem(ctx, bodys)
    if err != nil {
        return Response{
            Headers: corsHeaders,
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
    }

    response := Response{
        StatusCode:      200,
        IsBase64Encoded: false,
        Body:            "Request processed successfully",
        Headers: corsHeaders,
    }

    return response, nil

}


func main() {
    lambda.Start(HandleRequest)
}


func deleteItem(ctx context.Context, bodys Bodys) error {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        slog.Error(err.Error())
        return err
    }

    client := dynamodb.NewFromConfig(cfg)

    for _, body := range bodys {

        input := &dynamodb.DeleteItemInput{
            Key: map[string]types.AttributeValue{
                "id": &types.AttributeValueMemberS{Value: body.ID},
                "startTime": &types.AttributeValueMemberS{Value: body.StartTime},
            },
            TableName: aws.String("ScheduleInfo"),
        }

        _, err = client.DeleteItem(context.TODO(), input)
        if err != nil {
            slog.Error("Got error calling DeleteItem",err)
            return err
        }
    }

    slog.Info("Deleted the item")
    return nil
}


func validateBody(bodys Bodys) error {

    for _, body := range bodys {

        if body.ID == "" {
            return errors.New("id is required")
        }
        if body.StartTime == "" {
            return errors.New("startTime is required")
        }
        startTime, err := time.Parse(time.RFC3339, body.StartTime)
        if err != nil {
            return errors.New("startTime must be a valid datetime string in RFC3339 format")
        }
        FStartTime, _ := time.Parse("2006-01-02 15:04:05", startTime.Format("2006-01-02 15:04:05"))
        currentTime := time.Now().Add(3 * time.Minute)
        FCurrentTime, _ := time.Parse("2006-01-02 15:04:05", currentTime.Format("2006-01-02 15:04:05"))
        if FStartTime.Before(FCurrentTime) {
            return errors.New("startTime must be at least 3 minutes in the future")
        }
}
 
    return nil
}

