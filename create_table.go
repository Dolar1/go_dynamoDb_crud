// package main

// // snippet-start:[dynamodb.go.create_table.imports]
// import (
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/credentials"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/dynamodb"

// 	"fmt"
// 	"os"
// )

// // snippet-end:[dynamodb.go.create_table.imports]

// func main() {
// 	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAISBWLMXFHGYPJXUA")
// 	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")

// 	sess, _ := session.NewSession(&aws.Config{
// 		Region:      aws.String("us-east-2"),
// 		Credentials: credentials.NewEnvCredentials(),
// 	})

// 	// Create DynamoDB client
// 	svc := dynamodb.New(sess)

// 	// Create table Movies
// 	tableName := "Movies"

// 	input := &dynamodb.CreateTableInput{
// 		AttributeDefinitions: []*dynamodb.AttributeDefinition{
// 			{
// 				AttributeName: aws.String("Year"),
// 				AttributeType: aws.String("N"),
// 			},
// 			{
// 				AttributeName: aws.String("Title"),
// 				AttributeType: aws.String("S"),
// 			},
// 		},
// 		KeySchema: []*dynamodb.KeySchemaElement{
// 			{
// 				AttributeName: aws.String("Year"),
// 				KeyType:       aws.String("HASH"),
// 			},
// 			{
// 				AttributeName: aws.String("Title"),
// 				KeyType:       aws.String("RANGE"),
// 			},
// 		},
// 		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
// 			ReadCapacityUnits:  aws.Int64(10),
// 			WriteCapacityUnits: aws.Int64(10),
// 		},
// 		TableName: aws.String(tableName),
// 	}

// 	_, err := svc.CreateTable(input)
// 	if err != nil {
// 		fmt.Println("Got error calling CreateTable:")
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}

// 	fmt.Println("Created the table", tableName)
// 	// snippet-end:[dynamodb.go.create_table.call]
// }
