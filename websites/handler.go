package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/abtercms/abtercms2/pkg/dynamo"
	"github.com/abtercms/abtercms2/pkg/id"
	"github.com/abtercms/abtercms2/pkg/lhttp"
)

type listParams struct {
	ExclusiveStartKey string `lambda:"query.exclusive_start_key"` // a query parameter named "exclusive_start_key"
}

type listResponse struct {
	Items            interface{} `json:"items"`
	LastEvaluatedKey dynamo.Key  `json:"last_evaluated_key"`
	ScannedCount     int32       `json:"scanned_count"`
}

type entityParams struct {
	ID string `lambda:"path.id"` // a path parameter declared as :id
}

type website struct {
	ID   string `json:"pk"`
	Name string `json:"name"`
}

type repo interface {
	ListTables(context.Context, int32) ([]string, error)
	Get(context.Context, string, dynamo.Key, interface{}) error
	List(context.Context, string, int32, dynamo.Key, interface{}) (dynamo.Key, int32, error)
	Create(context.Context, string, interface{}) error
	Update(context.Context, string, interface{}) error
	Delete(context.Context, string, dynamo.Key) error
}

type Handler struct {
	repo repo
}

func getHandler(sdkConfig aws.Config, dynamoDbEndpoint string) *Handler {
	return &Handler{
		repo: dynamo.NewRepo(sdkConfig, dynamoDbEndpoint),
	}
}

func (h *Handler) ListTables(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	collection, err := h.repo.ListTables(ctx, limit)
	if err != nil {
		return lhttp.HandleError(err, nil)
	}

	return lmdrouter.MarshalResponse(http.StatusOK, nil, listResponse{Items: collection})
}

func (h *Handler) RetrieveCollection(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	var params listParams
	err = lmdrouter.UnmarshalRequest(req, false, &params)
	if err != nil {
		return lhttp.HandleError(fmt.Errorf(errUnmarshallRequest, err), nil)
	}

	var collection []website
	lastEvaluatedKey, scannedCount, err := h.repo.List(ctx, tableName, limit, dynamo.K1(params.ExclusiveStartKey), &collection)
	if err != nil {
		return lhttp.HandleError(err, nil)
	}

	return lmdrouter.MarshalResponse(http.StatusOK, nil, listResponse{Items: collection, LastEvaluatedKey: lastEvaluatedKey, ScannedCount: scannedCount})
}

func (h *Handler) CreateEntity(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	var entity website
	err = lmdrouter.UnmarshalRequest(req, true, &entity)
	if err != nil {
		return lhttp.HandleError(fmt.Errorf(errUnmarshallRequest, err), nil)
	}
	if entity.ID != "" {
		return lhttp.HandleError(fmt.Errorf(errPrimaryKeyNotAllowed, entity.ID), nil)
	}

	entity.ID, err = id.NewGenerator().NewString()
	if err != nil {
		return lhttp.HandleError(err, nil)
	}

	err = h.repo.Create(ctx, tableName, entity)
	if err != nil {
		return lhttp.HandleError(err, nil)
	}

	return lmdrouter.MarshalResponse(http.StatusCreated, nil, entity)
}

func (h *Handler) RetrieveEntity(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	var params entityParams
	err = lmdrouter.UnmarshalRequest(req, false, &params)
	if err != nil {
		return lhttp.HandleError(fmt.Errorf(errUnmarshallRequest, err), nil)
	}
	if params.ID == "" {
		return lhttp.HandleError(errors.New(errInvalidID), nil)
	}

	var entity website
	err = h.repo.Get(ctx, tableName, dynamo.K1(params.ID), &entity)
	if err != nil {
		return lhttp.HandleError(err, nil)
	}

	return lmdrouter.MarshalResponse(http.StatusOK, nil, entity)
}

func (h *Handler) UpdateEntity(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	var entity website
	err = lmdrouter.UnmarshalRequest(req, true, &entity)
	if err != nil {
		return lhttp.HandleError(fmt.Errorf(errUnmarshallRequest, err), nil)
	}

	var params entityParams
	err = lmdrouter.UnmarshalRequest(req, false, &params)
	if err != nil {
		return lhttp.HandleError(fmt.Errorf(errUnmarshallRequest, err), nil)
	}
	if params.ID == "" {
		return lhttp.HandleError(errors.New(errInvalidID), nil)
	}

	err = h.repo.Update(ctx, tableName, entity)
	if err != nil {
		return lhttp.HandleError(err, nil)
	}

	return lmdrouter.MarshalResponse(http.StatusOK, nil, entity)
}

func (h *Handler) DeleteEntity(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	var params entityParams
	err = lmdrouter.UnmarshalRequest(req, false, &params)
	if err != nil {
		return lhttp.HandleError(fmt.Errorf(errUnmarshallRequest, err), nil)
	}

	err = h.repo.Delete(ctx, tableName, dynamo.K1(params.ID))
	if err != nil {
		return lhttp.HandleError(err, nil)
	}

	return lmdrouter.MarshalResponse(http.StatusNoContent, nil, nil)
}
