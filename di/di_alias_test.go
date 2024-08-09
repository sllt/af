package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/////////////////////////////////////////////////////////////////////////////
// 							Explicit aliases
/////////////////////////////////////////////////////////////////////////////

func TestAs(t *testing.T) {
	is := assert.New(t)

	i := New()
	Provide(i, func(i Injector) (*lazyTestHeathcheckerOK, error) { return &lazyTestHeathcheckerOK{}, nil })

	is.Nil(As[*lazyTestHeathcheckerOK, Healthchecker](i))
	is.EqualError(As[*lazyTestShutdownerOK, Healthchecker](i), "DI: `*github.com/sllt/af/di.lazyTestShutdownerOK` is not `github.com/sllt/af/di.Healthchecker`")
	is.EqualError(As[*lazyTestHeathcheckerKO, Healthchecker](i), "DI: service `*github.com/sllt/af/di.lazyTestHeathcheckerKO` has not been declared")
	is.EqualError(As[*lazyTestShutdownerOK, *lazyTestShutdownerOK](i), "DI: service `*github.com/sllt/af/di.lazyTestShutdownerOK` has not been declared")
}

func TestMustAs(t *testing.T) {
	// @TODO
}

func TestAsNamed(t *testing.T) {
	is := assert.New(t)

	i := New()
	Provide(i, func(i Injector) (*lazyTestHeathcheckerOK, error) { return &lazyTestHeathcheckerOK{}, nil })

	is.Nil(AsNamed[*lazyTestHeathcheckerOK, Healthchecker](i, "*github.com/sllt/af/di.lazyTestHeathcheckerOK", "github.com/sllt/af/di.Healthchecker"))
	is.EqualError(AsNamed[*lazyTestShutdownerOK, Healthchecker](i, "*github.com/sllt/af/di.lazyTestShutdownerOK", "github.com/sllt/af/di.Healthchecker"), "DI: `*github.com/sllt/af/di.lazyTestShutdownerOK` is not `github.com/sllt/af/di.Healthchecker`")
	is.EqualError(AsNamed[*lazyTestHeathcheckerKO, Healthchecker](i, "*github.com/sllt/af/di.lazyTestHeathcheckerKO", "github.com/sllt/af/di.Healthchecker"), "DI: service `*github.com/sllt/af/di.lazyTestHeathcheckerKO` has not been declared")
	is.EqualError(AsNamed[*lazyTestShutdownerOK, *lazyTestShutdownerOK](i, "*github.com/sllt/af/di.lazyTestShutdownerOK", "*github.com/sllt/af/di.lazyTestShutdownerOK"), "DI: service `*github.com/sllt/af/di.lazyTestShutdownerOK` has not been declared")
}

func TestMustAsNamed(t *testing.T) {
	// @TODO
}

/////////////////////////////////////////////////////////////////////////////
// 							Implicit aliases
/////////////////////////////////////////////////////////////////////////////

func TestInvokeAs(t *testing.T) {
	is := assert.New(t)

	i := New()
	Provide(i, func(i Injector) (*lazyTestHeathcheckerOK, error) {
		return &lazyTestHeathcheckerOK{foobar: "hello world"}, nil
	})

	// found
	svc0, err := InvokeAs[*lazyTestHeathcheckerOK](i)
	is.EqualValues(&lazyTestHeathcheckerOK{foobar: "hello world"}, svc0)
	is.Nil(err)

	// found via interface
	svc1, err := InvokeAs[Healthchecker](i)
	is.EqualValues(&lazyTestHeathcheckerOK{foobar: "hello world"}, svc1)
	is.Nil(err)

	// not found
	svc2, err := InvokeAs[Shutdowner](i)
	is.Empty(svc2)
	is.EqualError(err, "DI: could not find service satisfying interface `github.com/sllt/af/di.Shutdowner`, available services: `*github.com/sllt/af/di.lazyTestHeathcheckerOK`")
}

func TestMustInvokeAs(t *testing.T) {
	// @TODO
}
