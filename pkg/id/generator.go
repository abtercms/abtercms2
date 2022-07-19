package id

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

const (
	errGeneratingID = "failed to generating id, err: %w"
)

type Generator struct {
	entropy *rand.Rand
	time    int64
}

func (g *Generator) getTime() time.Time {
	if g.time == 0 {
		return time.Now()
	}

	return time.UnixMilli(g.time)
}

func (g *Generator) setTime(t int64) *Generator {
	g.time = t

	return g
}

func (g *Generator) New() (ulid.ULID, error) {
	ms := ulid.Timestamp(g.getTime())

	return ulid.New(ms, g.entropy)
}

func (g *Generator) NewString() (string, error) {
	ms := ulid.Timestamp(g.getTime())

	u, err := ulid.New(ms, g.entropy)
	if err != nil {
		return "", fmt.Errorf(errGeneratingID, err)
	}

	return u.String(), nil
}

func NewGenerator() *Generator {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &Generator{
		entropy: entropy,
	}
}
