package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
)

type Event struct {
    Invoked
    Stream
}

type Invoked struct {
    ID string `json:"ID"`
}

type Stream events.DynamoDBEvent

func HandleRequest(ctx context.Context, event Event) (string, error) {

    slog.Info("start lambda function")

    slog.Info("event: ","%s",event)

    //文字列を複数持つ配列を定義する
    var pkList []string

    switch {
    case !reflect.DeepEqual(event.Invoked, Invoked{}):

        slog.Info("invoked: ","%s",event)

        pkList = append(pkList, event.ID)

    case !reflect.DeepEqual(event.Stream, events.DynamoDBEvent{}):
        slog.Info("stream: ","%s",event)

        for _, record := range event.Records {
            // if record.EventName == "REMOVE"{
            //     //削除された場合は処理をスキップする
            //     slog.Info("skip remove event")
            //     continue
            // }

            keys := record.Change.Keys
            slog.Info("keys: ","%s", keys)
            key := keys["id"].String()

            //pkListにkeyを追加するすでにある場合は追加しない
            if !contains(pkList, key) {
                pkList = append(pkList, key)
            }
        }

    default:
        slog.Info("event: ","%s",event)
        slog.Info("event: ","%s",reflect.TypeOf(event))
        slog.Error("unknown event")
        return "", fmt.Errorf("unknown event")
    }

    slog.Info("pkList: ","%s", pkList)

    if len(pkList) == 0 {
        slog.Info("No items found")
        return "No items found", nil
    }

    //pkLisのサイズ分だけループして、getScheduleを実行する
    for _, pk := range pkList {
        
        scheduleItems,err := getSchedule(pk)
        if err != nil {
            slog.Error(err.Error())
            return "", err
        }
        if len(scheduleItems) != 0 {
            schedule := scheduleItems[0]
            //scheduleの中身を確認する
            slog.Info("schedule: ","%s", schedule)
            slog.Info("Schedulevalues: ","%s", schedule["startTime"].(*types.AttributeValueMemberS).Value)
            scheduled := schedule["scheduled"].(*types.AttributeValueMemberS).Value
            slog.Info("scheduled: ","%b", scheduled)
            
        }

        newestSchedule,err := getNewestSchedule(pk)
        if err != nil {
            slog.Error(err.Error())
            return "", err
        }
        
        //newestScheduleの中身を確認する
        slog.Info("newestSchedule: ","%s", newestSchedule)
        slog.Info("newestSchedulevalues: ","%s", newestSchedule["startTime"].(*types.AttributeValueMemberS).Value)
        newestScheduleFlg := newestSchedule["scheduled"].(*types.AttributeValueMemberS).Value
        slog.Info("newestScheduleFlg: ","%b", newestScheduleFlg)

        id := newestSchedule["id"].(*types.AttributeValueMemberS).Value
        name := newestSchedule["name"].(*types.AttributeValueMemberS).Value
        startTime := newestSchedule["startTime"].(*types.AttributeValueMemberS).Value
        endTime := newestSchedule["endTime"].(*types.AttributeValueMemberS).Value
        
        //scheduleとnewestScheduleの内容が異なる==最新のものがtrueになっていない==イベントを上書きする必要がある
        if len(scheduleItems) == 1 && !reflect.DeepEqual(scheduleItems[0], newestSchedule){
            createRule(id, "start", name, startTime)
            createRule(id, "end", name, endTime)
            updateFlagTrue(pk, startTime)

            oldPk := scheduleItems[0]["id"].(*types.AttributeValueMemberS).Value
            oldSK :=scheduleItems[0]["startTime"].(*types.AttributeValueMemberS).Value
            updateSchedule(oldPk, oldSK)

        }else if len(scheduleItems) == 0 {
            createRule(id, "start", name, startTime)
            createRule(id, "end", name, endTime)
            updateFlagTrue(pk, startTime)

        }else{
            slog.Error("No need update schedule")
        }
    }

    return "Successfully create schedule", nil
}

func main() {
    lambda.Start(HandleRequest)
}

func getSchedule(pk string) ([]map[string]types.AttributeValue,error){

    slog.Info("start getSchedule")

    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        slog.Error(err.Error())
    }

    svc := dynamodb.NewFromConfig(cfg)

    now := time.Now().Format(time.RFC3339)
    
    slog.Info("now: ","%s", now)

    params := &dynamodb.QueryInput{
        TableName: aws.String("ScheduleInfo"),
        IndexName: aws.String("id-scheduled-index"),
        KeyConditionExpression: aws.String("id = :id AND scheduled = :scheduled"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":id":        &types.AttributeValueMemberS{Value: pk},
            ":scheduled": &types.AttributeValueMemberS{Value: "true"},
        },
        ScanIndexForward: aws.Bool(true),
        Limit:            aws.Int32(1),
    }


    resp, err := svc.Query(context.TODO(), params)
    if err != nil {
        slog.Error(err.Error())
        return nil, err
    }

    slog.Info("resp: ","%s", resp)

    return resp.Items, nil
    }

func getNewestSchedule(pk string) (map[string]types.AttributeValue,error){

    slog.Info("start getNewestSchedule")

    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        slog.Error(err.Error())
    }
    
    svc := dynamodb.NewFromConfig(cfg)
    
    now := time.Now().Format(time.RFC3339)
        
    slog.Info("now: ","%s", now)
    
    params := &dynamodb.QueryInput{
        TableName: aws.String("ScheduleInfo"),
        KeyConditionExpression: aws.String("id = :id AND startTime >= :now"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":id": &types.AttributeValueMemberS{Value: pk},
            ":now": &types.AttributeValueMemberS{Value: now},
        },
        ScanIndexForward: aws.Bool(true),
        Limit:            aws.Int32(1), 
    }
    
    
    resp, err := svc.Query(context.TODO(), params)
    if err != nil {
        slog.Error(err.Error())
        return nil, err
    }
    
    if len(resp.Items) > 0 {
        itemString := fmt.Sprintf("%v", resp.Items[0])
        slog.Info(itemString)
        return resp.Items[0], nil
    } else {
        slog.Error("No items found")
        return nil, fmt.Errorf("No items found")
    }
}

func updateSchedule(pk string, sk string) error {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        slog.Error(err.Error())
        return err
    }

    svc := dynamodb.NewFromConfig(cfg)

    params := &dynamodb.UpdateItemInput{
        TableName: aws.String("ScheduleInfo"),
        Key: map[string]types.AttributeValue{
            "id":        &types.AttributeValueMemberS{Value: pk},
            "startTime": &types.AttributeValueMemberS{Value: sk},
        },
        ExpressionAttributeNames: map[string]string{
            "#S": "scheduled",
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":s": &types.AttributeValueMemberS{Value: "false"},
        },
        UpdateExpression: aws.String("SET #S = :s"),
        ReturnValues:     types.ReturnValueUpdatedNew,
    }

    _, err = svc.UpdateItem(context.TODO(), params)
    if err != nil {
        slog.Error(err.Error())
        return err
    }

    slog.Info("Successfully updated schedule")
    return nil
}

// eventbridge schedulerに変更する
func createRule(id string ,flg string, name string, time string) {

    slog.Info("start createRule")

    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        slog.Error(err.Error())
    }

    client := eventbridge.NewFromConfig(cfg)

    convertedCron, err := convertISO8601ToCron(time)
    if err != nil {
        slog.Error(err.Error())
        return
    }
    // 文字列"cron"にconvertCeronを埋めこみ、変数にする
    cron := fmt.Sprintf("cron(%s)", convertedCron)
    
    ruleName := flg + "-" + name

    // Delete existing rule with the same name
    deleteRuleInput := &eventbridge.DeleteRuleInput{
        Name: aws.String(ruleName),
    }
    _, err = client.DeleteRule(context.TODO(), deleteRuleInput)
    if err != nil {
        slog.Error(err.Error())
        // Continue even if the rule does not exist
    }

    ruleInput := &eventbridge.PutRuleInput{
        Name:               aws.String(ruleName),
        ScheduleExpression: aws.String(cron),
        State:              eventTypes.RuleStateEnabled,
    }

    _, err = client.PutRule(context.TODO(), ruleInput)
    if err != nil {
        slog.Error(err.Error())
        return
    }

    //環境変数ARNを取得する
    arn,_ := os.LookupEnv("ARN")

    //UUIDを生成する
    uuid, _ := uuid.NewRandom()

    var action string

    if flg == "end" {
        action = "close"
    } else {
        action = "open"
    }

    event := map[string]interface{}{
        "id":     id,
        "action": action,
        "sk":     time,
    }
    slog.Info("event: ","%s", event)

    jsonEvent, err := json.Marshal(event)
    if err != nil {
        slog.Error(err.Error())
        return
    }

    //targetがある場合はすべて削除する
    deleteTargetsInput := &eventbridge.ListTargetsByRuleInput{
        Rule: aws.String(ruleName),
    }

    targetsOutput, err := client.ListTargetsByRule(context.TODO(), deleteTargetsInput)
    if err != nil {
        slog.Error(err.Error())
        return
    }

    if len(targetsOutput.Targets) > 0 {
        var targetIds []string
        for _, target := range targetsOutput.Targets {
            targetIds = append(targetIds, *target.Id)
        }
        
        removeTargetsInput := &eventbridge.RemoveTargetsInput{
            Rule: aws.String(ruleName),
            Ids: targetIds,
        }
        
        _, err = client.RemoveTargets(context.TODO(), removeTargetsInput)     
        
        if err != nil {
            slog.Error(err.Error())
            return
        }
    }


    targets := []eventTypes.Target{
        {
            Arn:   aws.String(arn),
            Id:    aws.String(uuid.String()),
            Input: aws.String(string(jsonEvent)),
        },
    }

    targetsInput := &eventbridge.PutTargetsInput{
        Rule:    aws.String(ruleName),
        Targets: targets,
    }

    _, err = client.PutTargets(context.TODO(), targetsInput)
    if err != nil {
        slog.Error(err.Error())
        return
    }

    slog.Info("Successfully created one-time scheduled rule")
}

func convertISO8601ToCron(t string) (string, error) {

    slog.Info("start convert cron")

    parsedTime, err := time.Parse(time.RFC3339, t)
    if err != nil {
        return "", err
    }

    parsedTime = parsedTime.Add(-9 * time.Hour)

    cron := fmt.Sprintf("%d %d %d %d ? %d", parsedTime.Minute(), parsedTime.Hour(), parsedTime.Day(), parsedTime.Month(), parsedTime.Year())

    return cron, nil
}


func contains(s []string, str string) bool {
    for _, v := range s {
        if v == str {
            return true
        }
    }
    return false
}


func updateFlagTrue(pk string, sk string) error {

    slog.Info("start updateFlagTrue")

    ctx := context.TODO()
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        slog.Error(err.Error())
    }

    client := dynamodb.NewFromConfig(cfg)

    input := &dynamodb.UpdateItemInput{
        ExpressionAttributeNames: map[string]string{
            "#F": "scheduled",
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":f": &types.AttributeValueMemberS{Value: "true"},
        },
        TableName: aws.String("ScheduleInfo"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: pk},
            "startTime": &types.AttributeValueMemberS{Value: sk},
        },
        UpdateExpression: aws.String("set #F = :f"),
    }

    _, err = client.UpdateItem(ctx, input)

    if err != nil {
        slog.Error(err.Error())
    }

    return err
}
