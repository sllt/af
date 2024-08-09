package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInjectorOrDefault(t *testing.T) {
	is := assert.New(t)

	is.Equal(DefaultRootScope, getInjectorOrDefault(nil))
	is.NotEqual(DefaultRootScope, getInjectorOrDefault(New()))

	type test struct {
		foobar string
	}

	DefaultRootScope = New()

	Provide(nil, func(i Injector) (*test, error) {
		return &test{foobar: "42"}, nil
	})

	is.Len(DefaultRootScope.self.services, 1)
}
