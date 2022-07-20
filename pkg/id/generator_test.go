package id_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtercms/abtercms2/pkg/id"
)

func TestGenerator_New(t *testing.T) {
	t.Parallel()

	// system under test
	sut := id.NewGenerator()

	// execute
	got0 := sut.New()
	got1 := sut.New()

	// asserts
	assert.NotEmpty(t, got0.String())
	assert.NotEmpty(t, got1.String())
	assert.NotEqual(t, got0.String(), got1.String())
}

func TestGenerator_NewString(t *testing.T) {
	t.Parallel()

	// system under test
	sut := id.NewGenerator()

	// execute
	got0 := sut.NewString()
	got1 := sut.NewString()

	// asserts
	assert.NotEmpty(t, got0)
	assert.NotEmpty(t, got1)
	assert.NotEqual(t, got0, got1)
}
