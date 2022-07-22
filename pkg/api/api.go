package api

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Error struct {
	Message string `json:"error_msg,omitempty"`
}

// Response returns a response to be returned by the API Gateway Reequest.
func Response(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: status,
		Body:       string(EncodeResponse(body)),
	}, nil
}

// StatusBadRequest returns a response of an HTTP StatusBadRequest and an error message.
func StatusBadRequest(err error) (*events.APIGatewayProxyResponse, error) {
	body := EncodeResponse(Error{Message: err.Error()})

	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: http.StatusBadRequest,
		Body:       string(body),
	}, nil
}
