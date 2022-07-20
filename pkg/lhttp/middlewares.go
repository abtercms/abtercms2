package lhttp

import (
	"context"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware(next lmdrouter.Handler) lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		res, err := next(ctx, req)

		path := req.HTTPMethod + " " + req.Path

		if err != nil {
			log.Error().
				Int("status", ToProblem(err).Status).
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
