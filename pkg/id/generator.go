package id

import (
	"github.com/oklog/ulid/v2"
)

type Generator struct{}

func (g *Generator) New() ulid.ULID {
	return ulid.Make()
}

func (g *Generator) NewString() string {
	return ulid.Make().String()
}

func NewGenerator() *Generator {
	return &Generator{}
}
