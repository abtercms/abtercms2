package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/abtercms/abtercms2/pkg/dynamo"
	"github.com/abtercms/abtercms2/websites/mocks"
)

func TestHandler_RetrieveCollection(t *testing.T) {
	t.Run("fail error in retrieving collection causes 500 internal server error", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites",
			HTTPMethod: http.MethodGet,
			QueryStringParameters: map[string]string{
				"exclusive_start_key": "qux",
			},
		}
		esk := dynamo.K1("qux")
		exclusiveStartKey := &esk
		scannedCountStub := int32(0)
		lastEvaluatedKeyStub := dynamo.Key{}

		// expectations
		expectedStatus := http.StatusInternalServerError

		// system under test
		sut, repoMock := createTestHandler()

		// mocks
		repoMock.On("List", ctx, limit, exclusiveStartKey, mock.Anything).
			Once().
			Return(lastEvaluatedKeyStub, scannedCountStub, assert.AnError)

		// execute
		res, err := sut.RetrieveCollection(ctx, requestStub)

		// asserts
		assert.Error(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("success w/o exclusive start key", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites",
			HTTPMethod: http.MethodGet,
			QueryStringParameters: map[string]string{
				"exclusive_start_key": "qux",
			},
		}
		esk := dynamo.K1("qux")
		exclusiveStartKey := &esk
		scannedCountStub := int32(30)
		lastEvaluatedKeyStub := dynamo.K1("foo")

		// expectations
		expectedStatus := http.StatusOK

		// system under test
		sut, repoMock := createTestHandler()

		// mocks
		repoMock.On("List", ctx, limit, exclusiveStartKey, mock.Anything).
			Once().
			Return(lastEvaluatedKeyStub, scannedCountStub, nil)

		// execute
		res, err := sut.RetrieveCollection(ctx, requestStub)

		// asserts
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("success /w exclusive start key", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites",
			HTTPMethod: http.MethodGet,
		}
		var exclusiveStartKey *dynamo.Key
		scannedCountStub := int32(30)
		lastEvaluatedKeyStub := dynamo.K1("foo")

		// expectations
		expectedStatus := http.StatusOK

		// system under test
		sut, repoMock := createTestHandler()

		// mocks
		repoMock.On("List", ctx, limit, exclusiveStartKey, mock.Anything).
			Once().
			Return(lastEvaluatedKeyStub, scannedCountStub, nil)

		// execute
		res, err := sut.RetrieveCollection(ctx, requestStub)

		// asserts
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})
}

func TestHandler_CreateEntity(t *testing.T) {
	t.Run("fail parsing payload causes 400 bad request", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites",
			HTTPMethod: http.MethodPost,
			Body:       `{"name:"bar"}`,
		}

		// expectations
		expectedStatus := http.StatusBadRequest

		// system under test
		sut, _ := createTestHandler()

		// execute
		res, err := sut.CreateEntity(ctx, requestStub)

		// asserts
		assert.Error(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("fail entity with existing id causes 400 bad request", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites",
			HTTPMethod: http.MethodPost,
			Body:       `{"pk":"foo","name":"bar"}`,
		}

		// expectations
		expectedStatus := http.StatusBadRequest

		// system under test
		sut, _ := createTestHandler()

		// execute
		res, err := sut.CreateEntity(ctx, requestStub)

		// asserts
		assert.Error(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("fail error in creating entity causes 500 internal server error", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites",
			HTTPMethod: http.MethodPost,
			Body:       `{"name":"bar"}`,
		}

		// expectations
		expectedStatus := http.StatusInternalServerError

		// system under test
		sut, repoMock := createTestHandler()

		// mocks
		repoMock.On("Create", ctx, mock.AnythingOfType("website")).
			Once().
			Return(assert.AnError)

		// execute
		res, err := sut.CreateEntity(ctx, requestStub)

		// asserts
		assert.Error(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites",
			HTTPMethod: http.MethodPost,
			Body:       `{"name":"bar"}`,
		}

		// expectations
		expectedStatus := http.StatusCreated

		// system under test
		sut, repoMock := createTestHandler()

		// mocks
		repoMock.On("Create", ctx, mock.AnythingOfType("website")).
			Once().
			Return(nil)

		// execute
		res, err := sut.CreateEntity(ctx, requestStub)

		// asserts
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})
}

func TestHandler_RetrieveEntity(t *testing.T) {
	t.Run("fail missing id in query parameters causes 400 bad request", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:           "/websites/foo",
			HTTPMethod:     http.MethodGet,
			PathParameters: map[string]string{},
		}

		// expectation
		expectedStatus := http.StatusBadRequest

		// system under test
		sut, _ := createTestHandler()

		// execute
		res, err := sut.RetrieveEntity(ctx, requestStub)

		// asserts
		assert.Error(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("fail error in retrieving entity causes 500 internal server error", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites/foo",
			HTTPMethod: http.MethodGet,
			PathParameters: map[string]string{
				"id": "foo",
			},
		}
		keyStub := dynamo.K1("foo")

		// expectation
		expectedStatus := http.StatusInternalServerError

		// system under test
		sut, repoMock := createTestHandler()

		// mocks
		repoMock.On("Get", ctx, keyStub, mock.Anything).
			Once().
			Return(assert.AnError)

		// execute
		res, err := sut.RetrieveEntity(ctx, requestStub)

		// asserts
		assert.Error(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("fail retrieving empty entity causes 404 not found", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites/foo",
			HTTPMethod: http.MethodGet,
			PathParameters: map[string]string{
				"id": "foo",
			},
		}
		keyStub := dynamo.K1("foo")

		// expectation
		expectedStatus := http.StatusNotFound

		// system under test
		sut, repoMock := createTestHandler()

		// mocks
		repoMock.On("Get", ctx, keyStub, mock.Anything).
			Once().
			Return(nil)

		// execute
		res, err := sut.RetrieveEntity(ctx, requestStub)

		// asserts
		assert.Error(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites/foo",
			HTTPMethod: http.MethodGet,
			PathParameters: map[string]string{
				"id": "foo",
			},
		}
		keyStub := dynamo.K1("foo")

		// expectation
		expectedStatus := http.StatusOK

		// system under test
		sut, repoMock := createTestHandler()

		// mocks
		websiteModifier := mock.MatchedBy(func(input *website) bool {
			input.ID = "foo"
			input.Name = "bar"

			return true
		})
		repoMock.On("Get", ctx, keyStub, websiteModifier).
			Once().
			Return(nil)

		// execute
		res, err := sut.RetrieveEntity(ctx, requestStub)

		// asserts
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})
}

func TestHandler_UpdateEntity(t *testing.T) {
	t.Run("fail error parsing request body causes 400 bad request", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites/foo",
			HTTPMethod: http.MethodPut,
			PathParameters: map[string]string{
				"id": "foo",
			},
			Body: `{"pk":"foo,"name":"bar"}`,
		}

		// expectations
		expectedStatus := http.StatusBadRequest

		// system under test
		sut, _ := createTestHandler()

		// execute
		res, err := sut.UpdateEntity(ctx, requestStub)

		// asserts
		assert.Error(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("fail id in request body not matching id in path causes 400 bad request", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites/foo",
			HTTPMethod: http.MethodPut,
			PathParameters: map[string]string{
				"id": "foo",
			},
			Body: `{"pk":"foo2","name":"bar"}`,
		}

		// expectations
		expectedStatus := http.StatusBadRequest

		// system under test
		sut, _ := createTestHandler()

		// execute
		res, err := sut.UpdateEntity(ctx, requestStub)

		// asserts
		assert.Error(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("fail empty id in path causes 400 bad request", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:           "/websites/foo",
			HTTPMethod:     http.MethodPut,
			PathParameters: map[string]string{},
			Body:           `{"pk":"foo","name":"bar"}`,
		}

		// expectations
		expectedStatus := http.StatusBadRequest

		// system under test
		sut, _ := createTestHandler()

		// execute
		res, err := sut.UpdateEntity(ctx, requestStub)

		// asserts
		assert.Error(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("fail error in updating entity causes 500 internal server error", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites/foo",
			HTTPMethod: http.MethodPut,
			PathParameters: map[string]string{
				"id": "foo",
			},
			Body: `{"pk":"foo","name":"bar"}`,
		}

		// expectations
		expectedStatus := http.StatusInternalServerError

		// system under test
		sut, repoMock := createTestHandler()

		// mocks
		repoMock.On("Update", ctx, mock.AnythingOfType("website")).
			Once().
			Return(assert.AnError)

		// execute
		res, err := sut.UpdateEntity(ctx, requestStub)

		// asserts
		assert.Error(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites/foo",
			HTTPMethod: http.MethodPut,
			PathParameters: map[string]string{
				"id": "foo",
			},
			Body: `{"pk":"foo","name":"bar"}`,
		}

		// expectations
		expectedStatus := http.StatusOK

		// system under test
		sut, repoMock := createTestHandler()

		// mocks
		repoMock.On("Update", ctx, mock.AnythingOfType("website")).
			Once().
			Return(nil)

		// execute
		res, err := sut.UpdateEntity(ctx, requestStub)

		// asserts
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})
}

func TestHandler_DeleteEntity(t *testing.T) {
	t.Run("fail error in deleting item causes 500 internal server error", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites/foo",
			HTTPMethod: http.MethodDelete,
			PathParameters: map[string]string{
				"id": "foo",
			},
		}
		keyStub := dynamo.K1("foo")

		// expectations
		expectedStatus := http.StatusInternalServerError

		// system under test
		sut, repoMock := createTestHandler()

		// mocks
		repoMock.On("Delete", ctx, keyStub).
			Once().
			Return(assert.AnError)

		// execute
		res, err := sut.DeleteEntity(ctx, requestStub)

		// asserts
		assert.Error(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// stubs
		ctx := context.Background()
		requestStub := events.APIGatewayProxyRequest{
			Path:       "/websites/foo",
			HTTPMethod: http.MethodDelete,
			PathParameters: map[string]string{
				"id": "foo",
			},
		}
		keyStub := dynamo.K1("foo")

		// expectations
		expectedStatus := http.StatusNoContent

		// system under test
		sut, repoMock := createTestHandler()

		// mocks
		repoMock.On("Delete", ctx, keyStub).
			Once().
			Return(nil)

		// execute
		res, err := sut.DeleteEntity(ctx, requestStub)

		// asserts
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, res.StatusCode)
	})
}

func createTestHandler() (*Handler, *mocks.Repo) {
	repoMock := &mocks.Repo{}

	sut := NewHandler(repoMock)

	return sut, repoMock
}
