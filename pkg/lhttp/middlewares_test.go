package lhttp_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/abtercms/abtercms2/pkg/lhttp"
)

func TestLoggerMiddleware(t *testing.T) {
	var buf bytes.Buffer

	// hack needed because zerolog gets a global log builder
	{
		l := log.Logger

		log.Logger = zerolog.New(&buf)
		defer func() {
			log.Logger = l
		}()
	}

	ctx := context.WithValue(context.Background(), "foo", "bar")
	requestStub := events.APIGatewayProxyRequest{
		HTTPMethod: http.MethodGet,
		Path:       "foo",
	}

	type args struct {
		next lmdrouter.Handler
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantLogRegex string
	}{
		{
			name: "no error",
			args: args{
				next: func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
					return events.APIGatewayProxyResponse{StatusCode: http.StatusNoContent}, nil
				},
			},
			wantErr:      false,
			wantLogRegex: "",
		},
		{
			name: "error",
			args: args{
				next: func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, errors.New("foo")
				},
			},
			wantErr:      true,
			wantLogRegex: `500.*foo`,
		},
		{
			name: "specific error",
			args: args{
				next: func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, lhttp.NewProblem(http.StatusBadRequest, "foo")
				},
			},
			wantErr:      true,
			wantLogRegex: `foo \(status 400\)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel() commented out because the log hack should not be run concurrently

			_, err := lhttp.LoggerMiddleware(tt.args.next)(ctx, requestStub)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			loggedText := buf.String()
			assert.Regexp(t, tt.wantLogRegex, loggedText)
			buf.Reset()
		})
	}
}
