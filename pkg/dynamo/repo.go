package dynamo

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/abtercms/abtercms2/pkg/lhttp"
)

const (
	privateKey = "pk"

	errMarshallItem    = "failed to marshal item"
	errFetchingItems   = "failed to fetch items"
	errFetchingItem    = "failed to fetch item"
	errCreatingItem    = "failed to create item"
	errUpdatingItem    = "failed to update item"
	errDeletingItem    = "failed to delete item"
	errUnmarshallItems = "failed to unmarshal items"
	errUnmarshallItem  = "failed to unmarshal item"

	problemDetailMarshaling    = "dynamodb marshaling issue"
	problemDetailUnmarshalling = "dynamodb unmarshalling issue"
	problemDetailDefault       = "dynamodb command issue"
)

type Key = map[string]types.AttributeValue

type Repo struct {
	db *dynamodb.Client
}

func NewRepo(sdkConfig aws.Config, isLocal bool) *Repo {
	db := dynamodb.NewFromConfig(sdkConfig, func(o *dynamodb.Options) {
		if isLocal {
			o.EndpointResolver = dynamodb.EndpointResolverFromURL("http://127.0.0.1:8000")
		}
	})
	return &Repo{
		db: db,
	}
}

func K1(id string) Key {
	return Key{
		privateKey: &types.AttributeValueMemberS{Value: id},
	}
}

func (r *Repo) ListTables(ctx context.Context, limit int32) ([]string, error) {
	out, err := r.db.ListTables(ctx, &dynamodb.ListTablesInput{
		Limit: &limit,
	})
	if err != nil {
		return nil, lhttp.Wrap(err, http.StatusNotFound, problemDetailDefault, errFetchingItems)
	}

	return out.TableNames, nil
}

func (r *Repo) Get(ctx context.Context, tableName string, key Key, result interface{}) error {
	out, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return lhttp.Wrap(err, http.StatusNotFound, problemDetailDefault, errFetchingItem)
	}

	err = attributevalue.UnmarshalMap(out.Item, result)
	if err != nil {
		return lhttp.Wrap(err, http.StatusNotFound, problemDetailUnmarshalling, errUnmarshallItem)
	}

	return nil
}

func (r *Repo) List(ctx context.Context, tableName string, limit int32, exclusiveStartKey Key, result interface{}) error {
	out, err := r.db.Scan(ctx, &dynamodb.ScanInput{
		TableName:         aws.String(tableName),
		Limit:             &limit,
		ExclusiveStartKey: exclusiveStartKey,
	})
	if err != nil {
		return lhttp.Wrap(err, http.StatusNotFound, problemDetailDefault, errFetchingItems)
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, result)
	if err != nil {
		return lhttp.Wrap(err, http.StatusNotFound, problemDetailUnmarshalling, errUnmarshallItems)
	}

	return nil
}

func (r *Repo) Create(ctx context.Context, tableName string, item interface{}) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return lhttp.Wrap(err, http.StatusNotFound, problemDetailMarshaling, errCreatingItem)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return lhttp.Wrap(err, http.StatusNotFound, problemDetailDefault, errCreatingItem)
	}

	return nil
}

func (r *Repo) Update(ctx context.Context, tableName string, item interface{}) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return lhttp.Wrap(err, http.StatusNotFound, problemDetailMarshaling, errMarshallItem)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	})

	if err != nil {
		return lhttp.Wrap(err, http.StatusNotFound, problemDetailDefault, errUpdatingItem)
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, tableName string, key Key) error {
	_, err := r.db.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(tableName),
	})

	if err != nil {
		return lhttp.Wrap(err, http.StatusNotFound, problemDetailDefault, errDeletingItem)
	}

	return nil
}
