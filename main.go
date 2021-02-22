package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
)

// CupCake struct
type CupCake struct {
	Month        string `json:"month"`
	CupcakeCount int    `json:"cupCakeCount"`
}

type Item struct {
	Month        string
	CupcakeCount int
}

//=========================================

// Get all books
func getCsv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(books)
}

func uploadCsv(w http.ResponseWriter, r *http.Request) {
	//===================================
	// reading the data
	//===================================
	w.Header().Set("Content-Type", "application/json")
	csvFile, _ := os.Open("multiTimeline.csv")
	reader := csv.NewReader(csvFile)
	var cupcake []CupCake
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		no, _ := strconv.Atoi(line[1])
		cupcake = append(cupcake, CupCake{
			Month:        line[0],
			CupcakeCount: no,
		})
	}

	// cupcakeJSON, _ := json.Marshal(cupcake)
	fmt.Println(cupcake)

	//===================================
	// uploading the data
	//===================================

	// json.NewEncoder(w).Encode(book)
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewEnvCredentials(),
	})

	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// add items to the table
	// item := CupCake{Month: "June", CupcakeCount: "10"}
	item := Item{
		Month:        "2004-01",
		CupcakeCount: 20,
	}

	// marshall map the items
	av, err := dynamodbattribute.MarshalMap(item)
	fmt.Println(av)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Cupcake"),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(err.Error())
		// os.Exit(1)
	} else {
		json.NewEncoder(w).Encode("Successfully added the data...!")
	}
}

// Get all books
func updateCsv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(books)
}

func createmydbtable(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewEnvCredentials(),
	})

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Create table Movies
	tableName := "Cupcake"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("CupcakeCount"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Month"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("CupcakeCount"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Month"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(err.Error())
		// os.Exit(1)
	} else {
		res := "Created the table " + tableName + " ..!"
		fmt.Println(res)
		json.NewEncoder(w).Encode(res)
	}
}

//=========================================

// Main function
func main() {
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAISBWLMXFHGYPJXUA")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")

	// Init router
	r := mux.NewRouter()

	// Route handles for reading and aws export to dynamo db
	r.HandleFunc("/makeTable", createmydbtable).Methods("GET")
	r.HandleFunc("/uploadcsv", uploadCsv).Methods("GET")
	r.HandleFunc("/getcsv", getCsv).Methods("GET")
	r.HandleFunc("/updatecsv", updateCsv).Methods("POST")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}
