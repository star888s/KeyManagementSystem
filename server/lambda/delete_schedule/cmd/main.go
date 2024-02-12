package main

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
)

func main() {

	lambda.Start(Handler)

	}


type MyEvent struct {
	ID string `json:"id"`
	Action string `json:"action"`
}


func Handler(ctx context.Context, event MyEvent) (string, error) {

	expiredList ,err := getExpiredList()
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info(strings.Join(expiredList, ", "))
	err = deleteRules(ctx, expiredList)
	if err != nil {
		slog.Error(err.Error())
	}
	err = deleteDynamoSchedule()
	if err != nil {
		slog.Error(err.Error())
	}

	return "", nil
}


func getExpiredList() ([]string, error){
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
			slog.Error(err.Error())
	}

	client := eventbridge.NewFromConfig(cfg)

	now := time.Now().UTC()

	// 空のリストを作成する
	var expiredList []string

	var nextToken *string
	for {
			output, err := client.ListRules(context.Background(), &eventbridge.ListRulesInput{NextToken: nextToken})
			if err != nil {
					return nil, err
			}

			for _, rule := range output.Rules {
				if rule.ScheduleExpression != nil {
					schedule, err := cronToJST(*rule.ScheduleExpression)
					if err != nil {
						slog.Error("Error parsing cron expression for rule:", *rule.Name,err)
						continue
					}

					if schedule.Before(now) {
						expiredList = append(expiredList, *rule.Name)
						// slog.Info("Rule:", *rule.Name, "has a schedule before the current time")
					}
				}
			}

			if output.NextToken == nil {
					break
			}
			nextToken = output.NextToken
	}
	return expiredList, nil
	
}

func cronToJST(cron string) (time.Time, error) {
	// Remove "cron(" and ")" from the cron string
	cron = strings.TrimPrefix(cron, "cron(")
	cron = strings.TrimSuffix(cron, ")")

	// Split the cron string into parts
	parts := strings.Split(cron, " ")

	// Parse each part of the cron string
	minute, _ := strconv.Atoi(parts[0])
	hour, _ := strconv.Atoi(parts[1])
	day, _ := strconv.Atoi(parts[2])
	month, _ := strconv.Atoi(parts[3])
	year, _ := strconv.Atoi(parts[5])

	// Create a time in UTC using the parsed parts
	utc := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)

	// Convert the time to JST
	jst := utc.In(time.FixedZone("Asia/Tokyo", 9*60*60))

	return jst, nil
}

func deleteRules(ctx context.Context, ruleNames []string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
			return err
	}

	client := eventbridge.NewFromConfig(cfg)

	for _, ruleName := range ruleNames {
			// Get the targets for the rule
			targetsOutput, err := client.ListTargetsByRule(ctx, &eventbridge.ListTargetsByRuleInput{
					Rule: &ruleName,
			})
			if err != nil {
					return err
			}

			// Remove the targets from the rule
			if len(targetsOutput.Targets) > 0 {
					targetIds := make([]string, len(targetsOutput.Targets))
					for i, target := range targetsOutput.Targets {
							targetIds[i] = *target.Id
					}

					_, err = client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
							Rule: &ruleName,
							Ids:  targetIds,
					})
					if err != nil {
							return err
					}
			}

			// Delete the rule
			_, err = client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{
					Name: &ruleName,
			})
			if err != nil {
					return err
			}
	}

	return nil
}

type Item struct {
	ID        string `dynamodbav:"id"`
	StartTime string `dynamodbav:"startTime"`
}

func deleteDynamoSchedule() error{
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
			return err
	}

	client := dynamodb.NewFromConfig(cfg)

	// ScheduleInfoテーブルを指定します
	tableName := "ScheduleInfo"

	// テーブルから全ての項目を取得します
	scanInput := &dynamodb.ScanInput{
			TableName: &tableName,
	}
	scanOutput, err := client.Scan(ctx, scanInput)
	if err != nil {
			return err
	}

	// 現在の時間を取得します
	now := time.Now()

	// 各項目のIDとStartTimeを取得します
	for _, item := range scanOutput.Items {
			var i Item
			err = attributevalue.UnmarshalMap(item, &i)
			if err != nil {
					return err
			}

			// StartTimeを時間として解析します
			startTime, err := time.Parse(time.RFC3339, i.StartTime)
			if err != nil {
					return err
			}

			// StartTimeが現在時刻より過去であれば項目を削除します
			if startTime.Before(now) {
					deleteInput := &dynamodb.DeleteItemInput{
							TableName: &tableName,
							Key: map[string]types.AttributeValue{
									"id":        &types.AttributeValueMemberS{Value: i.ID},
									"startTime": &types.AttributeValueMemberS{Value: i.StartTime},
							},
					}
					_, err = client.DeleteItem(ctx, deleteInput)
					if err != nil {
							return err
					}
			}
	}
	return nil
}
