package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

type Response events.APIGatewayProxyResponse

type Item struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
}

type Items []Item

type body struct {
    Items Items
    Message string
}

 
func HandleRequest(request events.APIGatewayProxyRequest) (Response, error) {

    corsHeaders := map[string]string{
        "Content-Type": "application/json",
        "Access-Control-Allow-Origin":  request.Headers["origin"],
        "Access-Control-Allow-Methods": "GET,OPTIONS",
        "Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
    }

    list,err := getAllSchedule()
    if err != nil {
        return Response{
            Headers: corsHeaders,
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
    }
    
    body := body{
        Items: list,
        Message: "Request processed successfully",
    }

    bodyJson, err := json.Marshal(body)
    if err != nil {
        return Response{
            Headers: corsHeaders,
            Body:       err.Error(),
            StatusCode: 400,
        }, nil
    }
    
    response := Response{
        Headers: corsHeaders,
        StatusCode:      200,
        Body:            string(bodyJson),
    }

    return response, nil

}


func main() {
    lambda.Start(HandleRequest)
}

func getAllSchedule() (Items ,error){
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        slog.Error(err.Error())
    }

    svc := dynamodb.NewFromConfig(cfg)

    params := &dynamodb.ScanInput{
        TableName: aws.String("KeyInfo"),
    }

    paginator := dynamodb.NewScanPaginator(svc, params)

    slog.Info("paginator: ","%s", paginator)

    var items Items

    for paginator.HasMorePages() {
        output, err := paginator.NextPage(context.TODO())
        if err != nil {
            return nil, err
        }
        slog.Info("output: ","%s", output)
        slog.Info("output: ","%s", output.Items)

        for _,i := range output.Items {
            slog.Info("item: ","%s", i)
            var item Item
            err := attributevalue.UnmarshalMap(i, &item)
            if err != nil {
                slog.Error("Failed to unmarshal:", err)
            }
            slog.Info("output: ","%s", item)
            items = append(items, item)
        }
        slog.Info("output: ","%s", items)
    }
    return items, nil
}
