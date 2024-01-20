package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/chmike/cmac-go"
)

func main() {

	lambda.Start(Handler)
	// info := &roomInfo{}
	// err := info.GetInfo("test")
	// if err != nil{
	// 	slog.Error("err: %s",err)
	// 	slog.Error("Failed to retrieve data from dynamo")
	// 	// return "Key was not toggled" ,errors.New("Failed to retrieve data from dynamo")
	}


type MyEvent struct {
	ID string `json:"id"`
	Action bool `json:"action"`
}

type roomInfo struct {
	id string
	name string
	uuid string
	secretKey string
}


// dynamo access
type Dynamo interface {
	GetInfo(id string) ()
}

func (i roomInfo)GetInfo(id string) (roomInfo, error){

	region := "ap-northeast-1"

	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = region
		return nil
	})
	if err != nil {
		slog.Error(err.Error())
		slog.Error("err: %s",err)
		slog.Error("Failed to create config")
	}

	table := aws.String("KeyInfo")

	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: id},
}

	svc := dynamodb.NewFromConfig(cfg)
	res, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
			TableName: table,
			Key: key,
	})

	if err != nil {
		slog.Error(err.Error())
		slog.Error("Failed to GetItem")
	}

	// fmt.Println(reflect.TypeOf(res))
	// fmt.Println(res)

	fmt.Println(reflect.TypeOf(res.Item))
	fmt.Println(res.Item)

// 	if val, ok := res.Item[""]; ok {
// 		fmt.Printf("foo exists. The value is %#v", val)
// }
	name := res.Item["name"].(*types.AttributeValueMemberS).Value
	uuid := res.Item["uuid"].(*types.AttributeValueMemberS).Value
	secretKey := res.Item["secretKey"].(*types.AttributeValueMemberS).Value

	i.id = id
	i.name = name
	i.uuid = uuid
	i.secretKey = secretKey
	fmt.Println(i)
	
	return i ,nil
}


// KeyAccess
type Sesame interface {
	toggleKey(info *roomInfo,Action bool, apiKey string)(string, error)
}

func (roomInfo)toggleKey(info *roomInfo,Action bool, apiKey string)(string, error){
	// fernet暗号化とかは後日検討
	// 88/82/83 = toggle/lock/unlock
	var cmd int
	var status string

	if Action == true {
		status = "closed"
		cmd = 83
	} else {
		status = "opened"
		cmd = 82
	}

	// env
	history := "regular exe"
	base64History := base64.StdEncoding.EncodeToString([]byte(history))

	slog.Info("create history","history", base64History)

	i := int32(time.Now().Unix())
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(i))
	byteKey, err := hex.DecodeString(info.secretKey)
	if err != nil {
		slog.Error("err: %s",err)
		slog.Error("Failed to decode string")
		return "", errors.New("Failed to decode string")
	}
	cm, err := cmac.New(aes.NewCipher, byteKey)
	if err != nil {
		slog.Error("err: %s",err)
		slog.Error("Failed to create cmac")
		return "", errors.New("Failed to create cmac")
	}

	cm.Write(buf[1:4])
	m := cm.Sum(nil)
	signature := hex.EncodeToString(m)

	slog.Info("create signature", "signature: %s", signature)

	//　env
	url := fmt.Sprintf("https://app.candyhouse.co/api/sesame2/%s/cmd", info.uuid)

	body := map[string]interface{}{
		"cmd":     cmd,
		"history": base64History,
		"sign":    signature,
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		slog.Error("err: %s",err)
		slog.Error("Failed to create request")
		return "", errors.New("Failed to create request")
	}

	headers := map[string]string{
		"x-api-key": apiKey,
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		slog.Error("err: %s",err)
		slog.Error("Failed to http request")
		return "", errors.New("Failed to create request")
	}
	defer res.Body.Close()

	slog.Info("Succeed request", "status: %s",res.Status)

	responseText := new(bytes.Buffer)
	responseText.ReadFrom(res.Body)
	slog.Info("Succeed request", "response: %s",responseText.String())

	return fmt.Sprintf("key was %s" , status),nil
}


func Handler(ctx context.Context, event MyEvent) (string,error) {
	// idを元にdynamodbに問い合わせる
	info := &roomInfo{}

	i,err := info.GetInfo(event.ID)
	if err != nil{
		slog.Error("err: %s",err)
		slog.Error("Failed to retrieve data from dynamo")
		return "Key was not toggled" ,errors.New("Failed to retrieve data from dynamo")
	}
	fmt.Println(i)
	fmt.Println(&i)

	apiKey,err2 := os.LookupEnv("APIKEY")
	if err2!= true {
		slog.Error("error","status",err2)
		slog.Error("Api key is not defined")
		return "Key was not toggled" ,errors.New("Api key is not defined")
	}

	status,err := info.toggleKey(&i,event.Action,apiKey)
	if err != nil {
		slog.Error("status: %b",err)
		slog.Error("Failed to toggle key")
		return "Key was not toggled" ,errors.New("Failed to toggle key")
	}

	slog.Info(status)

	return status, nil
}
