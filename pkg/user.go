package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ramshadows/Go-Serverless-App/utils"
)

// CreateUserRequest stores the create user request
type CreateUserRequest struct {
	UserId    string    `json:"userId,omitempty"`
	Username  string    `json:"username" binding:"required,alphanum"`
	Password  string    `json:"password" binding:"required,min=6"`
	FullName  string    `json:"full_name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// UserResponse stores the response that is returned to the user
type UserResponse struct {
	UserId    string    `json:"userId"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func CheckUserExist(tableName string) (*[]UserResponse, error) {
	// Creating session for client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	result, err := svc.Scan(input)
	if err != nil {
		return nil, errors.New("Failed to fetch users: " + err.Error())
	}
	item := new([]UserResponse)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)

	if err != nil {
		return nil, errors.New("Failed to Umarshall user list: " + err.Error())
	}
	return item, nil
}

func FetchUser(userId, tableName string) (*UserResponse, error) {
	// Creating session for client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Query DynamoDB to retrieve the created user
	// GetItem request
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(userId),
			},
		},
	})

	// Checking for errors, return error
	if err != nil {
		return nil, errors.New("Got error calling GetItem: " + err.Error())
	}

	// Create userRes of type UserResponse
	userRes := UserResponse{}

	// result is of type *dynamodb.GetItemOutput
	// result.Item is of type map[string]*dynamodb.AttributeValue
	// UnmarshallMap result.item into userRes
	err = dynamodbattribute.UnmarshalMap(result.Item, &userRes)

	if err != nil {
		return nil, errors.New("Failed to unmarshall result.item: " + err.Error())
	}

	fmt.Println("Returning user: ", &userRes)
	return &userRes, nil
}

func CreateUserReq(req events.APIGatewayProxyRequest, tableName string) (*UserResponse, error) {

	// Creating session for client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Unmarshal to Item to access object properties
	createUserReqString := req.Body
	createUserReqStruct := CreateUserRequest{}
	err := json.Unmarshal([]byte(createUserReqString), &createUserReqStruct)

	if err != nil {
		return nil, errors.New("Error Umarshalling user: " + err.Error())
	}
	// Generate a random int for user id
	userId := utils.RandomInt(1000, 10000)

	fmt.Println("Generated new User Id: ", userId)

	newUserId := strconv.Itoa(userId)

	// Validate email address
	if err := utils.ValidEmail(createUserReqStruct.Email); !err {
		return nil, errors.New("invalid email address: ")

	}

	fmt.Println("Validated Email: ", createUserReqStruct.Email)

	// compute an hashed password
	hashedPassword, err := utils.HashPassword(createUserReqStruct.Password)

	if err != nil {
		return nil, errors.New("Failed to encrypt password: " + err.Error())

	}

	fmt.Println("Hashed Password: ", hashedPassword)

	// Check empty value on Username
	if createUserReqStruct.Username == "" {
		return nil, errors.New("Username cannot be empty: " + err.Error())

	}

	fmt.Println("Validated Username: ", createUserReqStruct.Username)

	// Check empty value on FullName
	if createUserReqStruct.FullName == "" {
		return nil, errors.New("Full Name cannot be empty: " + err.Error())

	}

	fmt.Println("Validated FullName: ", createUserReqStruct.FullName)

	// Create new User of type CreateUser
	newUser := CreateUserRequest{
		UserId:    newUserId,
		Username:  createUserReqStruct.Username,
		FullName:  createUserReqStruct.FullName,
		Password:  hashedPassword,
		Email:     createUserReqStruct.Email,
		CreatedAt: time.Now(),
	}

	// Check if user exists
	userExist, err := CheckUserExist(tableName)

	if err != nil {
		return nil, errors.New("Got error checking if user already exist: " + err.Error())
	}

	for _, v := range *userExist {
		if v.Email == createUserReqStruct.Email {
			return nil, fmt.Errorf("user with email: %s already exists", v.Email)

		}

		if v.Username == createUserReqStruct.Username {
			return nil, fmt.Errorf("username already taken")

		}

	}

	// Marshal to dynamobb item
	createdUser, err := dynamodbattribute.MarshalMap(newUser)
	if err != nil {
		return nil, errors.New("Error marshalling new user: " + err.Error())
	}

	fmt.Printf("marshalled struct: %+v", createdUser)

	// Build put item input
	fmt.Printf("Putting new user: %v", createdUser)

	// Build put user input
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      createdUser,
	}

	_, err = svc.PutItem(input)

	// Checking for errors, return error
	if err != nil {
		return nil, errors.New("Got error calling PutItem: " + err.Error())
	}

	// Query DynamoDB to retrieve the created user
	// GetItem request
	userRes, err := FetchUser(newUserId, tableName)

	if err != nil {
		return nil, errors.New("Got error calling FetchUser(): " + err.Error())
	}

	return userRes, nil
}
