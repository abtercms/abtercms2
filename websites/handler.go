package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"

	"github.com/abtercms/abtercms2/pkg/dynamo"
	"github.com/abtercms/abtercms2/pkg/id"
	"github.com/abtercms/abtercms2/pkg/lhttp"
)

type repo interface {
	ListTables(context.Context, int32) ([]string, error)
	Get(context.Context, string, dynamo.Key, interface{}) error
	List(context.Context, string, int32, dynamo.Key, interface{}) error
	Create(context.Context, string, interface{}) error
	Update(context.Context, string, interface{}) error
	Delete(context.Context, string, dynamo.Key) error
}

type Handler struct {
	repo repo
}

func getHandler(sdkConfig aws.Config, isLocal bool) *Handler {
	return &Handler{
		repo: dynamo.NewRepo(sdkConfig, isLocal),
	}
}

func (h *Handler) TestAny(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	log.Log().
		Interface("request", req).
		Interface("environ", os.Environ()).
		Msg("testAny")

	collection := map[string]string{
		"foo": "bar",
		"bar": "baz",
	}

	return lmdrouter.MarshalResponse(http.StatusOK, nil, collection)
}

//func testDynamo(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
//	collection, err := repo.ListTables(ctx, limit)
//	if err != nil {
//		return lhttp.HandleError(err)
//	}
//
//	return lmdrouter.MarshalResponse(http.StatusOK, nil, collection)
//}

func (h *Handler) RetrieveCollection(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	var params listParams
	err = lmdrouter.UnmarshalRequest(req, false, &params)
	if err != nil {
		return lhttp.HandleError(fmt.Errorf(errUnmarshallRequest, err))
	}

	var collection []website
	err = h.repo.List(ctx, tableName, limit, dynamo.K1(params.ExclusiveStartKey), collection)
	if err != nil {
		return lhttp.HandleError(err)
	}

	return lmdrouter.MarshalResponse(http.StatusOK, nil, collection)
}

func (h *Handler) CreateEntity(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	var entity website
	err = lmdrouter.UnmarshalRequest(req, true, &entity)
	if err != nil {
		return lhttp.HandleError(fmt.Errorf(errUnmarshallRequest, err))
	}
	if entity.ID != "" {
		return lhttp.HandleError(fmt.Errorf(errPrimaryKeyNotAllowed, entity.ID))
	}

	entity.ID, err = id.NewGenerator().NewString()
	if err != nil {
		return lhttp.HandleError(err)
	}

	err = h.repo.Create(ctx, tableName, entity)
	if err != nil {
		return lhttp.HandleError(err)
	}

	return lmdrouter.MarshalResponse(http.StatusCreated, nil, entity)
}

func (h *Handler) RetrieveEntity(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	var params entityParams
	err = lmdrouter.UnmarshalRequest(req, false, &params)
	if err != nil {
		return lhttp.HandleError(fmt.Errorf(errUnmarshallRequest, err))
	}
	if params.ID == "" {
		return lhttp.HandleError(errors.New(errInvalidID))
	}

	var entity website
	err = h.repo.Get(ctx, tableName, dynamo.K1(params.ID), entity)
	if err != nil {
		return lhttp.HandleError(err)
	}

	return lmdrouter.MarshalResponse(http.StatusOK, nil, entity)
}

func (h *Handler) UpdateEntity(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	var entity website
	err = lmdrouter.UnmarshalRequest(req, true, &entity)
	if err != nil {
		return lhttp.HandleError(fmt.Errorf(errUnmarshallRequest, err))
	}

	var params entityParams
	err = lmdrouter.UnmarshalRequest(req, false, &params)
	if err != nil {
		return lhttp.HandleError(fmt.Errorf(errUnmarshallRequest, err))
	}
	if params.ID == "" {
		return lhttp.HandleError(errors.New(errInvalidID))
	}

	err = h.repo.Update(ctx, tableName, entity)
	if err != nil {
		return lhttp.HandleError(err)
	}

	return lmdrouter.MarshalResponse(http.StatusOK, nil, entity)
}

func (h *Handler) DeleteEntity(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	var params entityParams
	err = lmdrouter.UnmarshalRequest(req, false, &params)
	if err != nil {
		return lhttp.HandleError(fmt.Errorf(errUnmarshallRequest, err))
	}

	err = h.repo.Delete(ctx, tableName, dynamo.K1(params.ID))
	if err != nil {
		return lhttp.HandleError(err)
	}

	return lmdrouter.MarshalResponse(http.StatusNoContent, nil, nil)
}
