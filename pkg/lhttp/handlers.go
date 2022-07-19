package lhttp

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

const (
	headerContentType  = "Content-Type"
	contentTypeProblem = "application/problem+json; charset=UTF-8"
)

// HandleError returns a problem response from an error
func HandleError(err error, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	p := ToProblem(err)

	if headers == nil {
		headers = make(map[string]string)
	}
	if _, ok := headers[headerContentType]; !ok {
		headers[headerContentType] = contentTypeProblem
	}

	body, _ := json.Marshal(p)

	return events.APIGatewayProxyResponse{
		StatusCode:      p.Status,
		IsBase64Encoded: false,
		Headers:         headers,
		Body:            string(body),
	}, err
}
