package defaults

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Simple struct {
	A string `default:"str"`
	B int    `default:"1"`
	C bool   `default:"true"`
	D string
	E int `default:"11"`
}

func TestStruct(t *testing.T) {
	s := new(Simple)
	s.E = 10
	Set(s)
	assert.Equal(t, "str", s.A)
	assert.Equal(t, 1, s.B)
	assert.Equal(t, true, s.C)
	assert.Equal(t, 10, s.E)
}
