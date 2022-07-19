//go:generate mockery-latest --all --exported --case underscore
package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/abtercms/abtercms2/websites/mocks"
)

func TestRouter(t *testing.T) {
	log.Logger = zerolog.Nop()

	t.Run("retrieve collection", func(t *testing.T) {
		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites",
			HTTPMethod: http.MethodGet,
		}
		expectedStatus := http.StatusAccepted
		responseStub := events.APIGatewayProxyResponse{
			StatusCode: expectedStatus,
		}

		// mocks
		h := &mocks.Handler{}
		h.On("RetrieveCollection", ctx, requestStub).
			Once().
			Return(responseStub, nil)

		// system under test
		r := getRouter(h)

		// execute
		res, err := r.Handler(ctx, requestStub)

		// asserts
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("create entity", func(t *testing.T) {
		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites",
			HTTPMethod: http.MethodPost,
		}
		expectedStatus := http.StatusAccepted
		responseStub := events.APIGatewayProxyResponse{
			StatusCode: expectedStatus,
		}

		// mocks
		h := &mocks.Handler{}
		h.On("CreateEntity", ctx, requestStub).
			Once().
			Return(responseStub, nil)

		// system under test
		r := getRouter(h)

		// execute
		res, err := r.Handler(ctx, requestStub)

		// asserts
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("retrieve entity", func(t *testing.T) {
		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:           "/websites/abc",
			HTTPMethod:     http.MethodGet,
			PathParameters: map[string]string{"id": "abc"},
		}
		expectedStatus := http.StatusAccepted
		responseStub := events.APIGatewayProxyResponse{
			StatusCode: expectedStatus,
		}

		// mocks
		h := &mocks.Handler{}
		h.On("RetrieveEntity", ctx, requestStub).
			Once().
			Return(responseStub, nil)

		// system under test
		r := getRouter(h)

		// execute
		res, err := r.Handler(ctx, requestStub)

		// asserts
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("update entity", func(t *testing.T) {
		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:           "/websites/abc",
			HTTPMethod:     http.MethodPut,
			PathParameters: map[string]string{"id": "abc"},
		}
		expectedStatus := http.StatusAccepted
		responseStub := events.APIGatewayProxyResponse{
			StatusCode: expectedStatus,
		}

		// mocks
		h := &mocks.Handler{}
		h.On("UpdateEntity", ctx, requestStub).
			Once().
			Return(responseStub, nil)

		// system under test
		r := getRouter(h)

		// execute
		res, err := r.Handler(ctx, requestStub)

		// asserts
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("delete entity", func(t *testing.T) {
		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:           "/websites/abc",
			HTTPMethod:     http.MethodDelete,
			PathParameters: map[string]string{"id": "abc"},
		}
		expectedStatus := http.StatusAccepted
		responseStub := events.APIGatewayProxyResponse{
			StatusCode: expectedStatus,
		}

		// mocks
		h := &mocks.Handler{}
		h.On("DeleteEntity", ctx, requestStub).
			Once().
			Return(responseStub, nil)

		// system under test
		r := getRouter(h)

		// execute
		res, err := r.Handler(ctx, requestStub)

		// asserts
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})
}
