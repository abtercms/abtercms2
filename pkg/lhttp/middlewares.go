package lhttp

import (
	"context"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware(next lmdrouter.Handler) lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		path := req.HTTPMethod + " " + req.Path

		res, err := next(ctx, req)

		if err != nil {
			log.Error().
				Int("status", res.StatusCode).
				Str("error", err.Error()).
				Str("path", path).
				Send()

			return res, err
		}

		log.Info().
			Int("status", res.StatusCode).
			Str("path", path).
			Msg("success")

		return res, nil
	}
}
