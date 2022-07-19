package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/abtercms/abtercms2/pkg/lhttp"
)

const (
	limit int32 = 25

	EnvAwsRegion = "AWS_REGION"
	EnvTableName = "TABLE_NAME"

	errUnmarshallRequest    = "failed to unmarshalling request, err: %w"
	errInvalidID            = "primary key is required"
	errPrimaryKeyNotAllowed = "primary key is not allowed when creating entity: %s"
)

var (
	awsRegion string
	tableName string
	sdkConfig aws.Config
)

func init() {
	awsRegion = os.Getenv(EnvAwsRegion)
	tableName = os.Getenv(EnvTableName)

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var err error
	sdkConfig, err = config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
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
}

func main() {
	lambda.Start(getHandler(sdkConfig).TestAny)
	//lambda.Start(getRouter(getHandler(sdkConfig)).Handler)
}

type handler interface {
	RetrieveCollection(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	CreateEntity(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	RetrieveEntity(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	UpdateEntity(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	DeleteEntity(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	TestAny(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func getRouter(h handler) *lmdrouter.Router {
	router := lmdrouter.NewRouter("/websites", lhttp.LoggerMiddleware, lhttp.AuthMiddleware)
	router.Route(http.MethodGet, "", h.RetrieveCollection)
	router.Route(http.MethodPost, "", h.CreateEntity)
	router.Route(http.MethodGet, "/:id", h.RetrieveEntity)
	router.Route(http.MethodPut, "/:id", h.UpdateEntity)
	router.Route(http.MethodDelete, "/:id", h.DeleteEntity)
	router.Route(http.MethodConnect, "/:id", h.TestAny)

	return router
}

type listParams struct {
	ExclusiveStartKey string `lambda:"query.exclusive_start_key"` // a query parameter named "exclusive_start_key"
}

type entityParams struct {
	ID string `lambda:"path.id"` // a path parameter declared as :id
}

type website struct {
	ID   string `json:"pk"`
	Name string `json:"name"`
}
