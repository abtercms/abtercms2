package lhttp

import (
	"context"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware(next lmdrouter.Handler) lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
		path := req.HTTPMethod + " " + req.Path

		res, err = next(ctx, req)

		if err != nil {
			log.Error().
				Int("status", res.StatusCode).
				Str("error", err.Error()).
				Str("path", path).
				Msg(err.Error())

			return res, err
		}

		log.Info().
			Int("status", res.StatusCode).
			Str("path", path).
			Msg("success")

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
