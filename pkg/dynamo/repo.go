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
)

// Key represents a key ready to be used to find an entity in DynamoDB.
type Key = map[string]types.AttributeValue

type DB interface {
	Scan(context.Context, *dynamodb.ScanInput, ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

// Repo represents a repository capable of returning values for DynamoDB.
type Repo struct {
	db        DB
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

// SetDB sets a database client.
func (r *Repo) SetDB(db DB) *Repo {
	r.db = db

	return r
}

// List lists existing records in the table assigned to the repository.
func (r *Repo) List(ctx context.Context, limit int32, exclusiveStartKey *Key, result interface{}) (Key, int32, error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
		Limit:     &limit,
	}
	if exclusiveStartKey != nil {
		params.ExclusiveStartKey = *exclusiveStartKey
	}

	out, err := r.db.Scan(ctx, params)
	if err != nil {
		return Key{}, 0, lhttp.WrapProblem(err, http.StatusInternalServerError, errFetchingItems)
	}

	if out == nil {
		return Key{}, 0, lhttp.NewProblem(http.StatusInternalServerError, errFetchingItem)
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, result)
	if err != nil {
		return Key{}, 0, lhttp.WrapProblem(err, http.StatusInternalServerError, errUnmarshallItems)
	}

	return out.LastEvaluatedKey, out.ScannedCount, nil
}

// Create creates a new record in the table assigned to the repository.
func (r *Repo) Create(ctx context.Context, item interface{}) error {
	itemMarshalled, err := attributevalue.MarshalMap(item)
	if err != nil {
		return lhttp.WrapProblem(err, http.StatusBadRequest, errCreatingItem)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      itemMarshalled,
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return lhttp.WrapProblem(err, http.StatusInternalServerError, errCreatingItem)
	}

	return nil
}

// Get retrieves a record in the table assigned to the repository by key.
func (r *Repo) Get(ctx context.Context, key Key, result interface{}) error {
	out, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return lhttp.WrapProblem(err, http.StatusInternalServerError, errFetchingItem)
	}

	if out == nil {
		return lhttp.NewProblem(http.StatusInternalServerError, errFetchingItem)
	}

	err = attributevalue.UnmarshalMap(out.Item, result)
	if err != nil {
		return lhttp.WrapProblem(err, http.StatusInternalServerError, errUnmarshallItem)
	}

	return nil
}

// Update updates the existing record in the table assigned to the repository.
func (r *Repo) Update(ctx context.Context, item interface{}) error {
	itemMarshalled, err := attributevalue.MarshalMap(item)
	if err != nil {
		return lhttp.WrapProblem(err, http.StatusBadRequest, errMarshallItem)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      itemMarshalled,
		TableName: aws.String(r.tableName),
	})

	if err != nil {
		return lhttp.WrapProblem(err, http.StatusInternalServerError, errUpdatingItem)
	}

	return nil
}

// Delete deletes an existing record in the table assigned to the repository.
func (r *Repo) Delete(ctx context.Context, key Key) error {
	_, err := r.db.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(r.tableName),
	})

	if err != nil {
		return lhttp.WrapProblem(err, http.StatusInternalServerError, errDeletingItem)
	}

	return nil
}
