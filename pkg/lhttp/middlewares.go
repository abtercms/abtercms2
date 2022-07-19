package lhttp

import (
	"context"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware(next lmdrouter.Handler) lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
		code := http.StatusInternalServerError
		method := req.HTTPMethod
		path := req.Path

		res, err = next(ctx, req)
		if err == nil {
			code = res.StatusCode
		}

		log.Info().
			Int("Status", code).
			Err(err).
			Msg(method + " " + path)

		return res, err
	}
}

func AuthMiddleware(next lmdrouter.Handler) lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
		// TODO: Yes
		res, err = next(ctx, req)

		return res, err
	}
}
