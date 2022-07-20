package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/abtercms/abtercms2/pkg/lhttp"
)

const (
	limit int32 = 25

	EnvAwsRegion                = "AWS_REGION"
	EnvTableName                = "TABLE_NAME"
	EnvAwsSamLocal              = "AWS_SAM_LOCAL"
	EnvAwsDynamoDBLocalEndpoint = "AWS_DYNAMODB_LOCAL_ENDPOINT"

	trueString = "true"

	errUnmarshallRequest          = "failed to unmarshalling request, err: %w"
	errInvalidIDDetail            = "value in path: \"%s\", in payload: \"%s\", err: %w"
	errPrimaryKeyNotAllowedDetail = "primary key: \"%s\", err: %w"
)

var (
	errPrimaryKeyNotAllowed = lhttp.NewProblem(http.StatusBadRequest, "primary key is not allowed when creating entity.")
	errInvalidID            = lhttp.NewProblem(http.StatusBadRequest, "received ids are invalid.")
)

func main() {
	var (
		awsRegion        = os.Getenv(EnvAwsRegion)
		tableName        = os.Getenv(EnvTableName)
		dynamoDBEndpoint = ""
	)

	if os.Getenv(EnvAwsSamLocal) == trueString {
		dynamoDBEndpoint = os.Getenv(EnvAwsDynamoDBLocalEndpoint)
	}

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = awsRegion

		return nil
	})
	if err != nil {
		log.Fatal().
			Err(err).
			Str(EnvAwsRegion, awsRegion).
			Str(EnvTableName, tableName).
			Msg("cannot establish connection with dynamodb")
	}

	lambda.Start(getRouter(getHandler(sdkConfig, tableName, dynamoDBEndpoint)).Handler)
}

type handler interface {
	RetrieveCollection(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	CreateEntity(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	RetrieveEntity(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	UpdateEntity(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	DeleteEntity(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func getRouter(h handler) *lmdrouter.Router {
	router := lmdrouter.NewRouter("/websites", lhttp.LoggerMiddleware)
	router.Route(http.MethodGet, "", h.RetrieveCollection)
	router.Route(http.MethodPost, "", h.CreateEntity)
	router.Route(http.MethodGet, "/:id", h.RetrieveEntity)
	router.Route(http.MethodPut, "/:id", h.UpdateEntity)
	router.Route(http.MethodDelete, "/:id", h.DeleteEntity)

	return router
}
