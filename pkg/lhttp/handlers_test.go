package lhttp_test

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/abtercms/abtercms2/pkg/lhttp"
)

func TestHandleError(t *testing.T) {
	t.Parallel()

	type args struct {
		err     error
		headers map[string]string
	}

	tests := []struct {
		name       string
		args       args
		bodyRegexp string
		want       events.APIGatewayProxyResponse
	}{
		{
			name: "default w/o headers",
			args: args{
				err:     errors.New("foo"),
				headers: nil,
			},
			want: events.APIGatewayProxyResponse{
				StatusCode:      http.StatusInternalServerError,
				IsBase64Encoded: false,
				Headers: map[string]string{
					"Content-Type": "application/problem+json; charset=UTF-8",
				},
			},
			bodyRegexp: "Status/500",
		},
		{
			name: "default /w headers",
			args: args{
				err: errors.New("foo"),
				headers: map[string]string{
					"bar": "baz",
				},
			},
			want: events.APIGatewayProxyResponse{
				StatusCode:      http.StatusInternalServerError,
				IsBase64Encoded: false,
				Headers: map[string]string{
					"Content-Type": "application/problem+json; charset=UTF-8",
					"bar":          "baz",
				},
			},
			bodyRegexp: "Status/500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := lhttp.HandleError(tt.args.err, tt.args.headers)
			require.Error(t, err)

			tt.want.Body = got.Body
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleError() got = %v, want = %v", got, tt.want)
			}

			assert.Regexp(t, tt.bodyRegexp, got.Body)
		})
	}
}
