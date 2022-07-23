package api

import (
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Error struct {
	Message string `json:"error_msg,omitempty"`
}

// Response returns a response to be returned by the API Gateway Request.
func Response(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: status,
		Body:       string(EncodeResponse(body)),
	}, nil
}

// StatusOK returns a response of an HTTP StatusOK with body.
func StatusOK(body interface{}) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: http.StatusOK,
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

// StatusUnhandledMethod returns a response of an HTTP StatusMethodNotAllowed and an error message of unhandled method.
func StatusUnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	body := string(EncodeResponse(Error{Message: errors.New("unhandled method").Error()}))

	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: http.StatusMethodNotAllowed,
		Body:       body,
	}, nil
}

// StatusUnhandledRequest returns a response of an HTTP StatusNotImplemented and an error message.
func StatusUnhandledRequest(err error) (*events.APIGatewayProxyResponse, error) {
	body := string(EncodeResponse(Error{Message: err.Error()}))

	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: http.StatusNotImplemented,
		Body:       body,
	}, nil
}
