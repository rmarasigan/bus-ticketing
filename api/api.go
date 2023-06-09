package api

import (
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

const (
	CONTENT_TYPE = "application/json"
)

type Error struct {
	Message string `json:"error_msg,omitempty"`
}

// Response returns a response to be returned by the API Gateway Request.
func Response(status int, body interface{}) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE,
		},
		StatusCode: status,
		Body:       utility.EncodeJSON(body),
	}
}

// StatusOK returns a response of an HTTP StatusOK with body.
func StatusOK(body interface{}) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE,
		},
		StatusCode: http.StatusOK,
		Body:       utility.EncodeJSON(body),
	}, nil
}

// StatusOKWithoutBody returns a response of an HTTP StatusOK without body.
func StatusOKWithoutBody() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE,
		},
		StatusCode: http.StatusOK,
	}, nil
}

// StatusBadRequest returns a response of an HTTP StatusBadRequest and an error message.
func StatusBadRequest(err error) (*events.APIGatewayProxyResponse, error) {
	body := utility.EncodeJSON(Error{Message: err.Error()})

	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE,
		},
		StatusCode: http.StatusBadRequest,
		Body:       string(body),
	}, err
}

// StatusBadRequestWithBody returns a response of an HTTP StatusBadRequest with body.
func StatusBadRequestWithBody(body map[string]interface{}, err error) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE,
		},
		StatusCode: http.StatusBadRequest,
		Body:       utility.EncodeJSON(body),
	}, err
}

// StatusUnhandledMethod returns a response of an HTTP StatusMethodNotAllowed and an error message of unhandled method.
func StatusUnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE,
		},
		StatusCode: http.StatusMethodNotAllowed,
		Body:       utility.EncodeJSON(Error{Message: errors.New("unhandled method").Error()}),
	}, nil
}

// StatusInternalServerError returns a response of an HTTP StatusInternalServerError.
func StatusInternalServerError(err error) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE,
		},
		StatusCode: http.StatusInternalServerError,
	}, err
}

// StatusUnhandledRequest returns a response of an HTTP StatusNotImplemented and an error message.
func StatusUnhandledRequest(err error) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE,
		},
		StatusCode: http.StatusNotImplemented,
		Body:       utility.EncodeJSON(Error{Message: err.Error()}),
	}, err
}
