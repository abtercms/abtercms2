package id

import (
	"github.com/oklog/ulid/v2"
)

type Generator struct{}

func (g *Generator) New() ulid.ULID {
	return ulid.Make()
}

func NewGenerator() *Generator {
	return &Generator{}
}
