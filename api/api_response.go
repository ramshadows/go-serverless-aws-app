package api

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ramshadows/Go-Serverless-App/utils"
)

func ApiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = status

	stringBody,err := json.Marshal(body)

	if err != nil {
		return utils.ServerError(err)
	}
	resp.Body = string(stringBody)
	return &resp, nil
}
