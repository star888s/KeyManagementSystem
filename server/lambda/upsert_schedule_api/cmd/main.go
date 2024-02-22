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

    var body Body

    err := json.Unmarshal([]byte(request.Body), &body)
    if err != nil {
        fmt.Println("Could not decode body", err)
        return Response{
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
    }

    slog.Info("Received body: ", "%s",body)
    err = validateBody(body)
    if err != nil {
        return Response{
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
        
    }

    checkCondition, err := checkCondition(ctx, body)
    if err != nil {
        return Response{
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
    }

    if !checkCondition {
        return Response{
            Body:       "The schedule is already booked",
            StatusCode: 400,
        }, nil
    }

    err = upsertItem(ctx, body)
    if err != nil {
        return Response{
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
    }

    response := Response{
        StatusCode:      200,
        IsBase64Encoded: false,
        Body:            "Request processed successfully",
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }

    return response, nil

}


func main() {
    lambda.Start(HandleRequest)
}


func upsertItem(ctx context.Context, body Body) error {
    cfg, err := config.LoadDefaultConfig(ctx)

    svc := dynamodb.NewFromConfig(cfg)

    input := &dynamodb.PutItemInput{
        TableName: aws.String("ScheduleInfo"),
        Item: map[string]types.AttributeValue{
            "id":        &types.AttributeValueMemberS{Value: body.ID},
            "startTime": &types.AttributeValueMemberS{Value: body.StartTime},
            "endTime":   &types.AttributeValueMemberS{Value: body.EndTime},
            "name":      &types.AttributeValueMemberS{Value: body.Name},
            "scheduled": &types.AttributeValueMemberBOOL{Value: body.Scheduled},
            "memo":      &types.AttributeValueMemberS{Value: body.Memo},
        },
    }
    
    _, err = svc.PutItem(ctx, input)
    return err
}


func validateBody(body Body) error {
    if body.ID == "" {
        return errors.New("id is required")
    }
    if body.StartTime == "" {
        return errors.New("startTime is required")
    }
    if body.EndTime == "" {
        return errors.New("endTime is required")
    }
    startTime, err := time.Parse(time.RFC3339, body.StartTime)
    if err != nil {
        return errors.New("startTime must be a valid datetime string in RFC3339 format")
    }
    endTime, err := time.Parse(time.RFC3339, body.EndTime)
    if err != nil {
        return errors.New("endTime must be a valid datetime string in RFC3339 format")
    }
    FStartTime, _ := time.Parse("2006-01-02 15:04:05", startTime.Format("2006-01-02 15:04:05"))
    currentTime := time.Now().Add(3 * time.Minute)
    FCurrentTime, _ := time.Parse("2006-01-02 15:04:05", currentTime.Format("2006-01-02 15:04:05"))
    fmt.Printf("ct:%v", FStartTime)
    fmt.Printf("ct:%v", FCurrentTime)
    if FStartTime.Before(FCurrentTime) {
        return errors.New("startTime must be at least 3 minutes in the future")
    }
    if endTime.Before(startTime) {
        return errors.New("endTime must be after startTime")
    }
    if body.Name == "" {
        return errors.New("name is required")
    }
    return nil
}


func checkCondition(ctx context.Context,body Body) (bool, error) {
    cfg, err := config.LoadDefaultConfig(ctx)

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

    for _, item := range resp.Items {
        dbEndTime, ok := item["endTime"].(*types.AttributeValueMemberS)
        if !ok {
            continue
        }

        dbEndTimeParsed, err := time.Parse("2006-01-02T15:04:05Z", dbEndTime.Value)
        if err != nil {
            fmt.Printf("dbEndTime: %v", dbEndTime)
            return false, err
        }

        fmt.Printf("dbEndTimeParsed: %v", dbEndTimeParsed)

        bodyStartTimeParsed, err := time.Parse("2006-01-02T15:04:05Z", body.StartTime)
        if err != nil {
            fmt.Printf("bodyStartTime: %v", body.StartTime)
            return false, err
        }
        fmt.Printf("bodyStartTimeParsed: %v", bodyStartTimeParsed)

        if dbEndTimeParsed.After(bodyStartTimeParsed) {
            return false, nil
        }
    }

    return true, nil
}
