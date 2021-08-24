package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"

	"encoding/json"
	"fmt"
	"os"
)

type Content struct {
	Id     string `json:"id,omitempty"`
	UserID string `json:"user_id"`
	Title  string `json:"title"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Creating session for client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// New uuid for content id
	contentUuid := uuid.New().String()

	fmt.Println("Generated new content uuid:", contentUuid)
	// Unmarshal to Content to access object properties
	contentString := request.Body
	contentStruct := Content{}
	json.Unmarshal([]byte(contentString), &contentStruct)

	if contentStruct.Title == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400}, nil
	}

	// Create new content of type content
	content := Content{
		Id:     contentUuid,
		UserID: contentStruct.UserID,
		Title:  contentStruct.Title,
	}

	fmt.Println("content:", content)

	// Marshal to dynamobb content
	av, err := dynamodbattribute.MarshalMap(content)
	if err != nil {
		fmt.Println("Error marshalling content: ", err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	tableName := os.Getenv("DYNAMODB_TABLE")

	// Build put content input
	fmt.Println("Putting content: %v", av)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	// PutItem request
	_, err = svc.PutItem(input)

	// Checking for errors, return error
	if err != nil {
		fmt.Println("Got error calling PutItem: ", err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	// Marshal content to return
	itemMarshalled, err := json.Marshal(content)

	fmt.Println("Returning content: ", string(itemMarshalled))

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(itemMarshalled), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
