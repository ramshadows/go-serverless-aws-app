package main

/*
The reason why you need to entry points one local and one for the AWS lambda
is because in order to receive the event request from AWS API Gateway you are
going to need to add an API proxy.
unfortunately testing lambdamain.go is not really possible, so you need to have
another entry point for local development so you can test it using regular HTTP calls.
This is a very simple and direct way to test your lambdas locally.
*/

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ramshadows/Go-Serverless-App/api"
)

func main() {

	//main function forwards the request to the CreateUser Handler
	lambda.Start(api.CreateUser)
}
