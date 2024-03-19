package main

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

func TestHandleRequest(t *testing.T) {


    // req := testcontainers.ContainerRequest{
    //     Image:        "localstack/localstack",
    //     ExposedPorts: []string{"4566/tcp"},
    //     WaitingFor:   wait.ForLog("Ready."),
    //     Env: map[string]string{
    //         "SERVICES": "dynamodb",
    //     },
    // }

    // localstack, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
    //     ContainerRequest: req,
    //     Started:          true,
    // })
    // if err != nil {
    //     fmt.Println(err)
    // }

    // defer localstack.Terminate(ctx)

    // endpoint, err := localstack.Endpoint(ctx, "")
    // if err != nil {
    //     fmt.Println(err)
    // }

    t.Setenv("AWS_ACCESS_KEY_ID", "test")
    t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
    t.Setenv("AWS_ENDPOINT_URL", "http://localhost:4566")
    
    ctx := context.Background()

    createScheduleInfo(ctx)
    createKeyInfo(ctx)

    case1ID := "testdummy1"
    case1StartTime := stringTime(time.Now().Add(30 * time.Minute))
    case1EndTime := stringTime(time.Now().Add(60 * time.Minute))

    case1Key := &dynamodb.PutItemInput{
        TableName: aws.String("KeyInfo"),
        Item: map[string]types.AttributeValue{
            "id":        &types.AttributeValueMemberS{Value: case1ID},
            "name": &types.AttributeValueMemberS{Value: "testdummyname1"},
            "secretKey":   &types.AttributeValueMemberS{Value: "testdummysecretkey1"},
            "uuid":      &types.AttributeValueMemberS{Value: "testdummyuuid1"},
        },
    }
    case1Schedule := &dynamodb.PutItemInput{
        TableName: aws.String("ScheduleInfo"),
        Item: map[string]types.AttributeValue{
            "id":        &types.AttributeValueMemberS{Value: case1ID},
            "startTime": &types.AttributeValueMemberS{Value: case1StartTime},
            "endTime":   &types.AttributeValueMemberS{Value: case1EndTime},
            "name":      &types.AttributeValueMemberS{Value: "testdummyname1"},
            "scheduled": &types.AttributeValueMemberS{Value: "false"},
            "repetition": &types.AttributeValueMemberS{Value: "false"},
            "memo":      &types.AttributeValueMemberS{Value: "dummy memo"},
        },
    }

    inputList := []dynamodb.PutItemInput{*case1Key, *case1Schedule}

    for _, input := range inputList {
        createItem(ctx, input)
    }
    
    testCases := []struct {
        name    string
        event   Event
        wantErr bool
    }{
        {
            name: "Test Case 1",
            event: Event{
                Invoked: Invoked{
                },
                Stream: Stream{
                    Records: []events.DynamoDBEventRecord{
                        {
                            EventID:   "testEventID1",
                            EventName: "INSERT",
                            Change: events.DynamoDBStreamRecord{
                                Keys: map[string]events.DynamoDBAttributeValue{
                                    "id": events.NewStringAttribute(case1ID),
                                    "startTime" :events.NewStringAttribute(case1StartTime),
                                },
                            },
                        },
                    },
                },
            },
            wantErr: false,
        },
        // 他のテストケースをここに追加します
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            _, err := HandleRequest(context.Background(), tc.event)
            if (err != nil) != tc.wantErr {
                t.Errorf("HandleRequest() error = %v, wantErr %v", err, tc.wantErr)
            }
        })
    }
}


func stringTime(t time.Time) string{
    timeString := t.Format("2006-01-02T15:04:05-07:00")
    return timeString
}

func createScheduleInfo(ctx context.Context) error {

    slog.Info("Creating ScheduleInfo table")

    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return fmt.Errorf("unable to load SDK config, %v", err)
    }

    client := dynamodb.NewFromConfig(cfg)

    input := &dynamodb.CreateTableInput{
        AttributeDefinitions: []types.AttributeDefinition{
            {
                AttributeName: aws.String("id"),
                AttributeType: types.ScalarAttributeTypeS,
            },
            {
                AttributeName: aws.String("startTime"),
                AttributeType: types.ScalarAttributeTypeS,
            },
        },
        KeySchema: []types.KeySchemaElement{
            {
                AttributeName: aws.String("id"),
                KeyType:       types.KeyTypeHash,
            },
            {
                AttributeName: aws.String("startTime"),
                KeyType:       types.KeyTypeRange,
            },
        },
        ProvisionedThroughput: &types.ProvisionedThroughput{
            ReadCapacityUnits:  aws.Int64(5),
            WriteCapacityUnits: aws.Int64(5),
        },
        TableName: aws.String("ScheduleInfo"),
    }

    _, err = client.CreateTable(ctx, input)
    if err != nil {
        return fmt.Errorf("failed to create table, %v", err)
    }

    slog.Info("create table")

    return nil
}


func createKeyInfo(ctx context.Context) error {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return fmt.Errorf("unable to load SDK config, %v", err)
    }

    client := dynamodb.NewFromConfig(cfg)

    input := &dynamodb.CreateTableInput{
        AttributeDefinitions: []types.AttributeDefinition{
            {
                AttributeName: aws.String("id"),
                AttributeType: types.ScalarAttributeTypeS,
            },
        },
        KeySchema: []types.KeySchemaElement{
            {
                AttributeName: aws.String("id"),
                KeyType:       types.KeyTypeHash,
            },
        },
        TableName: aws.String("KeyInfo"),
    }

    _, err = client.CreateTable(ctx, input)
    if err != nil {
        return fmt.Errorf("failed to create table, %v", err)
    }

    return nil
}

func createItem(ctx context.Context, input dynamodb.PutItemInput) error {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return err
    }

    svc := dynamodb.NewFromConfig(cfg)
    
    _, err = svc.PutItem(ctx, &input)
    return err
}
