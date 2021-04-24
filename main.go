package main

import (
	// "crypto/rand"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
)

// CupCake struct
type CupCake struct {
	PrimaryId    string
	Month        string `json:"month"`
	CupcakeCount int    `json:"cupCakeCount"`
	Time         string
}

type Item struct {
	PrimaryId    string
	Month        string
	CupcakeCount int
	Time         string
}

//=========================================

// Get all books
func getCsv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(books)
}

func createmydbtable(w http.ResponseWriter, r *http.Request) {

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewEnvCredentials(),
	})
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Create table Movies
	tableName := "Cupcakes"

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

	//===================================
	// uploading the data
	//===================================

	// json.NewEncoder(w).Encode(book)
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewEnvCredentials(),
	})

	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// looping over data to write into the dynamodb
	for i := 0; i < len(cupcake); i++ {
		item := Item{
			PrimaryId:    strconv.Itoa(i + 1),
			Month:        cupcake[i].Month,
			CupcakeCount: cupcake[i].CupcakeCount,
			Time:         time.Now().String(),
		}

		// marshall map the items
		av, _ := dynamodbattribute.MarshalMap(item)
		fmt.Println(i)
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("Cupcakes"),
		}

		_, err = svc.PutItem(input)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}

	// sending the response from here...
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
func GetAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("site hit to get data")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewEnvCredentials(),
	})

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// getting al the data
	// Build the query input parameters
	params := &dynamodb.ScanInput{
		TableName: aws.String("Cupcakes"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}

	// sending the response from here...
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(err.Error())
		// os.Exit(1)
	} else {
		json.NewEncoder(w).Encode(result)
	}
}

// Get all books
func updateRandom100(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("update api hit")

	//generate a random list of intizers b/w 1-200
	var random_arr []int

	for len(random_arr) < 100 {
		var r = rand.Intn(200)
		random_arr = append(random_arr, r)
	}

	//look for primary key and update them with the time in a new row

	// write a call back to run this task every 5 mints
	// untill the request is manullay closed

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewEnvCredentials(),
	})

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	//update code
	tableName := "Cupcakes"

	for i := 0; i < len(random_arr); i++ {

		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#T": aws.String("Time"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":Time": {
					S: aws.String(time.Now().String()),
				},
			},
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"PrimaryId": {
					S: aws.String(strconv.Itoa(random_arr[i])),
				},
			},
			// ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("SET #T = :Time"),
		}

		_, err := svc.UpdateItem(input)

		fmt.Print(i, " => ", random_arr[i], "\n")

		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}

	// sending the response from here...
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		json.NewEncoder(w).Encode(err.Error())
		// os.Exit(1)
	} else {
		json.NewEncoder(w).Encode("Successfully updated")
	}
}

//=========================================

// Main function
func main() {
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIA2GJWK274ILBJFKWH")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")

	// Init router
	r := mux.NewRouter()

	// Route handles for reading and aws export to dynamo db
	r.HandleFunc("/makeTable", createmydbtable).Methods("GET")
	r.HandleFunc("/uploadcsv", uploadCsv).Methods("GET")
	r.HandleFunc("/getcsv", getCsv).Methods("GET")
	r.HandleFunc("/getall", GetAll).Methods("GET")
	r.HandleFunc("/updateRandom100", updateRandom100).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}
