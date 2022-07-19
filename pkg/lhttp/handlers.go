package lhttp

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

const (
	headerContentType  = "Content-Type"
	contentTypeProblem = "application/problem+json; charset=UTF-8"
)

// HandleError returns a problem response from an error.
func HandleError(err error, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	problem := ToProblem(err)

	if headers == nil {
		headers = make(map[string]string)
	}

	if _, ok := headers[headerContentType]; !ok {
		headers[headerContentType] = contentTypeProblem
	}

	body, err2 := json.Marshal(problem)
	if err2 != nil {
		err = fmt.Errorf("additional error in marshaling problem. marshaling err: %s, original err: %w", err2.Error(), err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      problem.Status,
		IsBase64Encoded: false,
		Headers:         headers,
		Body:            string(body),
	}, err
}
