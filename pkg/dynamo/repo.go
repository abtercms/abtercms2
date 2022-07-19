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

// Key represents a key ready to be used to find an entity in DynamoDB.
type Key = map[string]types.AttributeValue

// Repo represents a repository capable of returning values for DynamoDB.
type Repo struct {
	db        *dynamodb.Client
	tableName string
}

// NewRepo creates a new Repo instance.
func NewRepo(sdkConfig aws.Config, tableName, dynamoDBEndpoint string) *Repo {
	return &Repo{
		db: dynamodb.NewFromConfig(sdkConfig, func(o *dynamodb.Options) {
			if dynamoDBEndpoint != "" {
				o.EndpointResolver = dynamodb.EndpointResolverFromURL(dynamoDBEndpoint)
			}
		}),
		tableName: tableName,
	}
}

// K1 converts a string into a Key for DynamoDB.
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
		return nil, lhttp.WrapProblem(err, http.StatusNotFound, problemDetailDefault, errFetchingItems)
	}

	return out.TableNames, nil
}

func (r *Repo) Get(ctx context.Context, key Key, result interface{}) error {
	out, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return lhttp.WrapProblem(err, http.StatusNotFound, problemDetailDefault, errFetchingItem)
	}

	err = attributevalue.UnmarshalMap(out.Item, result)
	if err != nil {
		return lhttp.WrapProblem(err, http.StatusNotFound, problemDetailUnmarshalling, errUnmarshallItem)
	}

	return nil
}

func (r *Repo) List(ctx context.Context, limit int32, exclusiveStartKey Key, result interface{}) (Key, int32, error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
		Limit:     &limit,
	}
	if len(exclusiveStartKey) > 0 {
		params.ExclusiveStartKey = exclusiveStartKey
	}

	out, err := r.db.Scan(ctx, params)
	if err != nil {
		return Key{}, 0, lhttp.WrapProblem(err, http.StatusNotFound, problemDetailDefault, errFetchingItems)
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, result)
	if err != nil {
		return Key{}, 0, lhttp.WrapProblem(err, http.StatusNotFound, problemDetailUnmarshalling, errUnmarshallItems)
	}

	return out.LastEvaluatedKey, out.ScannedCount, nil
}

func (r *Repo) Create(ctx context.Context, item interface{}) error {
	itemMarshalled, err := attributevalue.MarshalMap(item)
	if err != nil {
		return lhttp.WrapProblem(err, http.StatusNotFound, problemDetailMarshaling, errCreatingItem)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      itemMarshalled,
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return lhttp.WrapProblem(err, http.StatusNotFound, problemDetailDefault, errCreatingItem)
	}

	return nil
}

func (r *Repo) Update(ctx context.Context, item interface{}) error {
	itemMarshalled, err := attributevalue.MarshalMap(item)
	if err != nil {
		return lhttp.WrapProblem(err, http.StatusNotFound, problemDetailMarshaling, errMarshallItem)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      itemMarshalled,
		TableName: aws.String(r.tableName),
	})

	if err != nil {
		return lhttp.WrapProblem(err, http.StatusNotFound, problemDetailDefault, errUpdatingItem)
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, key Key) error {
	_, err := r.db.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(r.tableName),
	})

	if err != nil {
		return lhttp.WrapProblem(err, http.StatusNotFound, problemDetailDefault, errDeletingItem)
	}

	return nil
}
