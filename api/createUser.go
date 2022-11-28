package api

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ramshadows/Go-Serverless-App/pkg"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

var (
	tableName = os.Getenv("DYNAMODB_USERS_TABLE")
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func CreateUser(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	result, err := pkg.CreateUserReq(req, tableName)
	if err != nil {
		return ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return ApiResponse(http.StatusCreated, result)

}
