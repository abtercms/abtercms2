package lhttp

import (
	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
)

// HandleError returns a problem response from an error
func HandleError(err error) (events.APIGatewayProxyResponse, error) {
	return lmdrouter.MarshalResponse(ToStatusCode(err), nil, ToProblem(err))
}
