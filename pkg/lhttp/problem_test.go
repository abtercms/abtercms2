package lhttp_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtercms/abtercms2/pkg/lhttp"
)

func TestNewProblem(t *testing.T) {
	t.Parallel()

	type args struct {
		status int
		msg    string
		args   []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "default",
			args: args{
				status: http.StatusBadRequest,
				msg:    "bar",
				args:   []interface{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// execute
			sut := lhttp.NewProblem(tt.args.status, tt.args.msg, tt.args.args...)

			// asserts
			assert.Equal(t, tt.args.status, sut.Status)
			assert.Contains(t, sut.Type, fmt.Sprintf("%d", tt.args.status))
			assert.Empty(t, sut.Detail)
		})
	}
}

func TestToProblem(t *testing.T) {
	t.Parallel()

	type args struct {
		err    error
		status int
		args   []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "default",
			args: args{
				err:    errors.New("foo"),
				status: http.StatusInternalServerError,
				args:   []interface{}{},
			},
		},
		{
			name: "bad request",
			args: args{
				err:    lhttp.NewProblem(http.StatusBadRequest, "foo"),
				status: http.StatusBadRequest,
				args:   []interface{}{},
			},
		},
		{
			name: "bad request wrapped",
			args: args{
				err:    fmt.Errorf("wrapped, err: %w", lhttp.NewProblem(http.StatusBadRequest, "foo")),
				status: http.StatusBadRequest,
				args:   []interface{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// execute
			sut := lhttp.ToProblem(tt.args.err)

			// asserts
			assert.Equal(t, tt.args.status, sut.Status)
			assert.Contains(t, sut.Type, fmt.Sprintf("%d", tt.args.status))
			assert.Equal(t, tt.args.err.Error(), sut.Detail)
		})
	}
}

func TestWrapProblem(t *testing.T) {
	t.Parallel()

	type args struct {
		err    error
		status int
		msg    string
		args   []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "default",
			args: args{
				err:    errors.New("foo"),
				status: http.StatusBadRequest,
				msg:    "bar",
				args:   []interface{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// execute
			sut := lhttp.WrapProblem(tt.args.err, tt.args.status, tt.args.msg, tt.args.args...)

			// asserts
			assert.Equal(t, tt.args.status, sut.Status)
			assert.Contains(t, sut.Type, fmt.Sprintf("%d", tt.args.status))
			assert.Empty(t, sut.Detail)
			assert.Contains(t, sut.Error(), tt.args.err.Error())
		})
	}
}
