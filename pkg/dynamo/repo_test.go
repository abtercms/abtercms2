package dynamo_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/abtercms/abtercms2/pkg/dynamo"
	"github.com/abtercms/abtercms2/pkg/lhttp"
	"github.com/abtercms/abtercms2/pkg/mocks"
)

func TestK1(t *testing.T) {
	t.Parallel()

	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want dynamo.Key
	}{
		{
			name: "default",
			args: args{
				id: "foo",
			},
			want: dynamo.Key{
				"pk": &types.AttributeValueMemberS{Value: "foo"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := dynamo.K1(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("K1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_List(t *testing.T) {
	ctx := context.WithValue(context.Background(), "foo", "bar")

	type T struct {
		Foo string
	}

	t.Run("fail returning nil from DynamoDB causes 500 internal server error", func(t *testing.T) {
		t.Parallel()

		// stubs
		var limitStub int32 = 25
		var exclusiveStartKeyStub dynamo.Key
		var itemStub *dynamodb.ScanOutput
		actualList := []T{}

		// system under test
		sut, dbMock := createTestRepo()

		// mocks
		dbMock.On("Scan", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).
			Once().
			Return(itemStub, nil)

		// execute
		_, _, err := sut.List(ctx, limitStub, exclusiveStartKeyStub, &actualList)

		// asserts
		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, lhttp.ToProblem(err).Status)
	})

	t.Run("fail error in retrieving item causes 500 internal server error", func(t *testing.T) {
		t.Parallel()

		// stubs
		var limitStub int32 = 25
		var exclusiveStartKeyStub dynamo.Key
		actualList := []T{}

		// system under test
		sut, dbMock := createTestRepo()

		// mocks
		dbMock.On("Scan", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).
			Once().
			Return(nil, assert.AnError)

		// execute
		_, _, err := sut.List(ctx, limitStub, exclusiveStartKeyStub, &actualList)

		// asserts
		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, lhttp.ToProblem(err).Status)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// stubs
		var limitStub int32 = 25
		exclusiveStartKeyStub := dynamo.K1("foo")
		actualList := []T{}

		var scannedCount int32 = 15
		exclusiveLastEvaluatedKey := dynamo.K1("foo")
		itemStubs := &dynamodb.ScanOutput{
			Items: []map[string]types.AttributeValue{
				{
					"Foo": &types.AttributeValueMemberS{Value: "bar"},
				},
				{
					"Foo": &types.AttributeValueMemberS{Value: "baz"},
				},
			},
			ScannedCount:     scannedCount,
			LastEvaluatedKey: exclusiveLastEvaluatedKey,
		}
		expectedResult := []T{{Foo: "bar"}, {Foo: "baz"}}

		// system under test
		sut, dbMock := createTestRepo()

		// mocks
		dbMock.On("Scan", ctx, mock.AnythingOfType("*dynamodb.ScanInput")).
			Once().
			Return(itemStubs, nil)

		// execute
		actualLastEvaluatedKey, actualScannedCount, err := sut.List(ctx, limitStub, exclusiveStartKeyStub, &actualList)

		// asserts
		require.NoError(t, err, "Get() error = %v", err)
		assert.Equal(t, exclusiveLastEvaluatedKey, actualLastEvaluatedKey)
		assert.Equal(t, scannedCount, actualScannedCount)
		assert.Equal(t, expectedResult, actualList)
	})
}

func TestRepo_Create(t *testing.T) {
	ctx := context.WithValue(context.Background(), "foo", "bar")

	t.Run("fail marshaling item causes 400 bad request", func(t *testing.T) {
		t.Parallel()

		// it's unclear if this can ever happen as it calls `attributevalue.MarshalMap` under the hood,
		// which is designed to always return something
		t.Skip()
	})

	t.Run("fail error in updating item causes 500 internal server error", func(t *testing.T) {
		t.Parallel()

		// stubs
		itemStub := map[string]string{
			"foo": "bar",
		}

		// system under test
		sut, dbMock := createTestRepo()

		// mocks
		dbMock.On("PutItem", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).
			Once().
			Return(nil, assert.AnError)

		// execute
		err := sut.Create(ctx, itemStub)

		// asserts
		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, lhttp.ToProblem(err).Status)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// stubs
		itemStub := struct{ foo string }{
			foo: "bar",
		}

		// system under test
		sut, dbMock := createTestRepo()

		dbMock.On("PutItem", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).
			Once().
			Return(nil, nil)

		err := sut.Create(ctx, itemStub)
		require.NoError(t, err, "Update() error = %v", err)
	})
}

func TestRepo_Get(t *testing.T) {
	ctx := context.WithValue(context.Background(), "foo", "bar")

	type T struct {
		Foo string
	}

	t.Run("fail returning nil from DynamoDB causes 500 internal server error", func(t *testing.T) {
		t.Parallel()

		// stubs
		keyStub := dynamo.K1("foo")
		var itemStub *dynamodb.GetItemOutput
		actualResult := T{}

		// system under test
		sut, dbMock := createTestRepo()

		// mocks
		dbMock.On("GetItem", ctx, mock.AnythingOfType("*dynamodb.GetItemInput")).
			Once().
			Return(itemStub, nil)

		// execute
		err := sut.Get(ctx, keyStub, &actualResult)

		// asserts
		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, lhttp.ToProblem(err).Status)
	})

	t.Run("fail error in retrieving item causes 500 internal server error", func(t *testing.T) {
		t.Parallel()

		// stubs
		keyStub := dynamo.K1("foo")
		actualResult := T{}

		// system under test
		sut, dbMock := createTestRepo()

		// mocks
		dbMock.On("GetItem", ctx, mock.AnythingOfType("*dynamodb.GetItemInput")).
			Once().
			Return(nil, assert.AnError)

		// execute
		err := sut.Get(ctx, keyStub, &actualResult)

		// asserts
		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, lhttp.ToProblem(err).Status)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// stubs
		keyStub := dynamo.K1("foo")
		itemStub := &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"Foo": &types.AttributeValueMemberS{Value: "bar"},
			},
		}
		expectedResult := T{Foo: "bar"}
		actualResult := T{}

		// system under test
		sut, dbMock := createTestRepo()

		// mocks
		dbMock.On("GetItem", ctx, mock.AnythingOfType("*dynamodb.GetItemInput")).
			Once().
			Return(itemStub, nil)

		// execute
		err := sut.Get(ctx, keyStub, &actualResult)

		// asserts
		require.NoError(t, err, "Get() error = %v", err)
		assert.Equal(t, expectedResult, actualResult)
	})
}

func TestRepo_Update(t *testing.T) {
	ctx := context.WithValue(context.Background(), "foo", "bar")

	t.Run("fail marshaling item causes 400 bad request", func(t *testing.T) {
		t.Parallel()

		// it's unclear if this can ever happen as it calls `attributevalue.MarshalMap` under the hood,
		// which is designed to always return something
		t.Skip()
	})

	t.Run("fail error in updating item causes 500 internal server error", func(t *testing.T) {
		t.Parallel()

		// stubs
		itemStub := map[string]string{
			"foo": "bar",
		}

		// system under test
		sut, dbMock := createTestRepo()

		// mocks
		dbMock.On("PutItem", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).
			Once().
			Return(nil, assert.AnError)

		// execute
		err := sut.Update(ctx, itemStub)

		// asserts
		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, lhttp.ToProblem(err).Status)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// stubs
		itemStub := struct{ foo string }{
			foo: "bar",
		}

		// system under test
		sut, dbMock := createTestRepo()

		// mocks
		dbMock.On("PutItem", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).
			Once().
			Return(nil, nil)

		// execute
		err := sut.Update(ctx, itemStub)

		// asserts
		require.NoError(t, err, "Update() error = %v", err)
	})
}

func TestRepo_Delete(t *testing.T) {
	ctx := context.WithValue(context.Background(), "foo", "bar")

	t.Run("fail error in deleting item causes 500 internal server error", func(t *testing.T) {
		t.Parallel()

		// stubs
		keyStub := dynamo.K1("foo")

		// system under test
		sut, dbMock := createTestRepo()

		// mocks
		dbMock.On("DeleteItem", ctx, mock.AnythingOfType("*dynamodb.DeleteItemInput")).
			Once().
			Return(nil, assert.AnError)

		// execute
		err := sut.Delete(ctx, keyStub)

		// asserts
		assert.Equal(t, http.StatusInternalServerError, lhttp.ToProblem(err).Status)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// stubs
		keyStub := dynamo.K1("foo")

		// system under test
		sut, dbMock := createTestRepo()

		// mocks
		dbMock.On("DeleteItem", ctx, mock.AnythingOfType("*dynamodb.DeleteItemInput")).
			Once().
			Return(nil, nil)

		// execute
		err := sut.Delete(ctx, keyStub)

		// asserts
		require.NoError(t, err, "Delete() error = %v", err)
	})
}

func createTestRepo() (*dynamo.Repo, *mocks.DB) {
	sut := dynamo.NewRepo(aws.Config{}, "fooTable", "")

	db := &mocks.DB{}
	sut.SetDB(db)

	return sut, db
}
