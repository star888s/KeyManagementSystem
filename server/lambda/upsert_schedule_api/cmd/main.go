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
    EndTime string `json:"endTime"`
    Name string `json:"name"`
    Scheduled bool `json:"scheduled"`
    Memo string `json:"memo"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {

    corsHeaders := map[string]string{
        "Content-Type": "application/json",
        "Access-Control-Allow-Origin":  request.Headers["origin"],
        "Access-Control-Allow-Methods": "POST,OPTIONS",
        "Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
    }

    var body Body

    err := json.Unmarshal([]byte(request.Body), &body)
    if err != nil {
        slog.Error("Could not decode body", err)
        return Response{
            Headers: corsHeaders,
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
    }

    slog.Info("Received body: ", "%s",body)
    err = validateBody(body)
    if err != nil {
        slog.Error("Could not validate body", err)
        return Response{
            Headers: corsHeaders,
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
        
    }

    checkCondition, err := checkCondition(ctx, body)
    if err != nil {
        slog.Error("Could not check condition", err)
        return Response{
            Headers: corsHeaders,
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
    }

    if !checkCondition {
        slog.Error("The schedule is already booked")
        return Response{
            Headers: corsHeaders,
            Body:       "The schedule is already booked",
            StatusCode: 400,
        }, nil
    }

    err = upsertItem(ctx, body)
    if err != nil {
        slog.Error("Could not upsert item", err)
        return Response{
            Headers: corsHeaders,
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
    }

    slog.Info("Request processed successfully")
    response := Response{
        Headers: corsHeaders,
        StatusCode:      200,
        IsBase64Encoded: false,
        Body:            "Request processed successfully",
    }

    return response, nil

}


func main() {
    lambda.Start(HandleRequest)
}


func upsertItem(ctx context.Context, body Body) error {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return err
    }

    svc := dynamodb.NewFromConfig(cfg)

    input := &dynamodb.PutItemInput{
        TableName: aws.String("ScheduleInfo"),
        Item: map[string]types.AttributeValue{
            "id":        &types.AttributeValueMemberS{Value: body.ID},
            "startTime": &types.AttributeValueMemberS{Value: body.StartTime},
            "endTime":   &types.AttributeValueMemberS{Value: body.EndTime},
            "name":      &types.AttributeValueMemberS{Value: body.Name},
            "scheduled": &types.AttributeValueMemberS{Value: "false"},
            "memo":      &types.AttributeValueMemberS{Value: body.Memo},
        },
    }
    
    _, err = svc.PutItem(ctx, input)
    return err
}


func validateBody(body Body) error {
    if body.ID == "" {
        slog.Error("id is required")
        return errors.New("id is required")
    }
    if body.StartTime == "" {
        slog.Error("startTime is required")
        return errors.New("startTime is required")
    }
    if body.EndTime == "" {
        slog.Error("endTime is required")
        return errors.New("endTime is required")
    }
    startTime, err := time.Parse(time.RFC3339, body.StartTime)
    if err != nil {
        slog.Error("startTime must be a valid datetime string in RFC3339 format")
        return errors.New("startTime must be a valid datetime string in RFC3339 format")
    }
    endTime, err := time.Parse(time.RFC3339, body.EndTime)
    if err != nil {
        slog.Error("endTime must be a valid datetime string in RFC3339 format")
        return errors.New("endTime must be a valid datetime string in RFC3339 format")
    }
    FStartTime, _ := time.Parse("2006-01-02 15:04:05", startTime.Format("2006-01-02 15:04:05"))
    currentTime := time.Now().Add(3 * time.Minute)
    FCurrentTime, _ := time.Parse("2006-01-02 15:04:05", currentTime.Format("2006-01-02 15:04:05"))
    if FStartTime.Before(FCurrentTime) {
        slog.Error("startTime must be at least 3 minutes in the future")
        return errors.New("startTime must be at least 3 minutes in the future")
    }
    if endTime.Before(startTime) {
        slog.Error("endTime must be after startTime")
        return errors.New("endTime must be after startTime")
    }
    if body.Name == "" {
        slog.Error("name is required")
        return errors.New("name is required")
    }
    return nil
}


func checkCondition(ctx context.Context,body Body) (bool, error) {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return false, err
    }
    slog.Info("checkCondition")

    svc := dynamodb.NewFromConfig(cfg)

    params := &dynamodb.QueryInput{
        TableName: aws.String("ScheduleInfo"),
        KeyConditionExpression: aws.String("id = :id"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":id": &types.AttributeValueMemberS{Value: body.ID},
        },
    }

    resp, err := svc.Query(ctx, params)
    if err != nil {
        return false, err
    }

    bodyStartTimeParsed, err := time.Parse("2006-01-02T15:04:05Z07:00", body.StartTime)
    if err != nil {
        fmt.Printf("bodyStartTime: %v", body.StartTime)
        return false, err
    }

    bodyEndTimeParsed, err := time.Parse("2006-01-02T15:04:05Z07:00", body.EndTime)
    if err != nil {
        fmt.Printf("bodyEndTime: %v", body.EndTime)
        return false, err
    }

    for _, item := range resp.Items {
        dbStartTime, ok := item["startTime"].(*types.AttributeValueMemberS)
        if !ok {
            continue
        }

        dbEndTime, ok := item["endTime"].(*types.AttributeValueMemberS)
        if !ok {
            continue
        }

        dbStartTimeParsed, err := time.Parse("2006-01-02T15:04:05Z07:00", dbStartTime.Value)
        if err != nil {
            fmt.Printf("dbStartTime: %v", dbStartTime)
            return false, err
        }

        dbEndTimeParsed, err := time.Parse("2006-01-02T15:04:05Z07:00", dbEndTime.Value)
        if err != nil {
            fmt.Printf("dbEndTime: %v", dbEndTime)
            return false, err
        }

        if (dbStartTimeParsed.Before(bodyEndTimeParsed) && dbEndTimeParsed.After(bodyEndTimeParsed)) ||
           (dbStartTimeParsed.Before(bodyStartTimeParsed) && dbEndTimeParsed.After(bodyStartTimeParsed)) ||
           (dbStartTimeParsed.After(bodyStartTimeParsed) && dbEndTimeParsed.Before(bodyEndTimeParsed)) {
            return false, nil
        }
    }

    return true, nil
}
