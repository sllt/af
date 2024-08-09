package di

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceAlias(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	i := New()

	service1 := newServiceAlias[string, string]("foobar1", i, "foobar2")
	is.Equal("foobar1", service1.name)
	is.Equal(i, service1.scope)
	is.Equal("foobar2", service1.targetName)
}

func TestServiceAlias_getName(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	i := New()

	service1 := newServiceAlias[string, string]("foobar1", i, "foobar2")
	is.Equal("foobar1", service1.getName())
}

func TestServiceAlias_getTypeName(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	i := New()

	service1 := newServiceAlias[string, int]("foobar1", i, "foobar2")
	is.Equal("int", service1.getTypeName())
}

func TestServiceAlias_getServiceType(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	i := New()

	service1 := newServiceAlias[string, string]("foobar1", i, "foobar2")
	is.Equal(ServiceTypeAlias, service1.getServiceType())
}

func TestServiceAlias_getEmptyInstance(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	svc := newServiceAlias[string, lazyTest]("foo", nil, "bar")
	is.Empty(svc.getEmptyInstance())
	is.EqualValues(lazyTest{}, svc.getEmptyInstance())
}

func TestServiceAlias_getInstanceAny(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	i := New()
	Provide(i, func(i Injector) (*lazyTestHeathcheckerOK, error) {
		return &lazyTestHeathcheckerOK{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTestHeathcheckerOK, Healthchecker](i))

	// basic type
	service1 := newServiceAlias[*lazyTestHeathcheckerOK, Healthchecker]("github.com/sllt/af/di.Healthchecker", i, "*github.com/sllt/af/di.lazyTestHeathcheckerOK")
	instance1, err1 := service1.getInstanceAny(i)
	is.Nil(err1)
	is.EqualValues(&lazyTestHeathcheckerOK{foobar: "foobar"}, instance1)

	// target service not found
	service2 := newServiceAlias[*lazyTestHeathcheckerOK, int]("github.com/sllt/af/di.Healthchecker", i, "int")
	instance2, err2 := service2.getInstanceAny(i)
	is.EqualError(err2, "DI: could not find service `int`, available services: `*github.com/sllt/af/di.lazyTestHeathcheckerOK`, `github.com/sllt/af/di.Healthchecker`")
	is.EqualValues(0, instance2)

	Provide(i, func(i Injector) (int, error) {
		return 42, nil
	})

	// target service found but not convertible type
	service3 := newServiceAlias[*lazyTestHeathcheckerOK, int]("github.com/sllt/af/di.Healthchecker", i, "int")
	instance3, err3 := service3.getInstanceAny(i)
	is.EqualError(err3, "DI: service found, but type mismatch: invoking `*github.com/sllt/af/di.lazyTestHeathcheckerOK` but registered `int`")
	is.EqualValues(0, instance3)

	// @TODO: missing test with child scopes
	// @TODO: missing test with stacktrace
}

func TestServiceAlias_getInstance(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	i := New()
	Provide(i, func(i Injector) (*lazyTestHeathcheckerOK, error) {
		return &lazyTestHeathcheckerOK{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTestHeathcheckerOK, Healthchecker](i))

	// basic type
	service1 := newServiceAlias[*lazyTestHeathcheckerOK, Healthchecker]("github.com/sllt/af/di.Healthchecker", i, "*github.com/sllt/af/di.lazyTestHeathcheckerOK")
	instance1, err1 := service1.getInstance(i)
	is.Nil(err1)
	is.EqualValues(&lazyTestHeathcheckerOK{foobar: "foobar"}, instance1)

	// target service not found
	service2 := newServiceAlias[*lazyTestHeathcheckerOK, int]("github.com/sllt/af/di.Healthchecker", i, "int")
	instance2, err2 := service2.getInstance(i)
	is.EqualError(err2, "DI: could not find service `int`, available services: `*github.com/sllt/af/di.lazyTestHeathcheckerOK`, `github.com/sllt/af/di.Healthchecker`")
	is.EqualValues(0, instance2)

	Provide(i, func(i Injector) (int, error) {
		return 42, nil
	})

	// target service found but not convertible type
	service3 := newServiceAlias[*lazyTestHeathcheckerOK, int]("github.com/sllt/af/di.Healthchecker", i, "int")
	instance3, err3 := service3.getInstance(i)
	is.EqualError(err3, "DI: service found, but type mismatch: invoking `*github.com/sllt/af/di.lazyTestHeathcheckerOK` but registered `int`")
	is.EqualValues(0, instance3)

	// @TODO: missing test with child scopes
	// @TODO: missing test with stacktrace
}

func TestServiceAlias_isHealthchecker(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// no healthcheck
	i1 := New()
	Provide(i1, func(i Injector) (*lazyTest, error) {
		return &lazyTest{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTest, any](i1))
	service1, _ := i1.serviceGet("interface {}")
	is.False(service1.(Service[any]).isHealthchecker())

	// healthcheck ok
	i2 := New()
	Provide(i2, func(i Injector) (*lazyTestHeathcheckerOK, error) {
		return &lazyTestHeathcheckerOK{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTestHeathcheckerOK, Healthchecker](i2))
	service2, _ := i2.serviceGet("github.com/sllt/af/di.Healthchecker")
	is.False(service2.(serviceIsHealthchecker).isHealthchecker())
	_, _ = service2.(serviceGetInstanceAny).getInstanceAny(nil)
	is.True(service2.(serviceIsHealthchecker).isHealthchecker())

	// healthcheck ko
	i3 := New()
	Provide(i3, func(i Injector) (*lazyTestHeathcheckerKO, error) {
		return &lazyTestHeathcheckerKO{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTestHeathcheckerKO, Healthchecker](i3))
	service3, _ := i3.serviceGet("github.com/sllt/af/di.Healthchecker")
	is.False(service3.(serviceIsHealthchecker).isHealthchecker())
	_, _ = service3.(serviceGetInstanceAny).getInstanceAny(nil)
	is.True(service3.(serviceIsHealthchecker).isHealthchecker())

	// service not found (wrong type)
	i4 := New()
	service4 := newServiceAlias[*lazyTestHeathcheckerKO, Healthchecker]("github.com/sllt/af/di.Healthchecker", i4, "*github.com/sllt/af/di.lazyTestHeathcheckerKO")
	is.False(service4.isHealthchecker())
	_, _ = service4.getInstanceAny(nil)
	is.False(service4.isHealthchecker())

	// service not found (wrong name)
	i5 := New()
	service5 := newServiceAlias[*lazyTestHeathcheckerOK, Healthchecker]("github.com/sllt/af/di.Healthchecker", i5, "*github.com/sllt/af/di.lazyTestHeathcheckerKO")
	is.False(service5.isHealthchecker())
	_, _ = service5.getInstanceAny(nil)
	is.False(service5.isHealthchecker())
}

// @TODO: missing tests for context
func TestServiceAlias_healthcheck(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := context.Background()

	// no healthcheck
	i1 := New()
	Provide(i1, func(i Injector) (*lazyTest, error) {
		return &lazyTest{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTest, any](i1))
	service1, _ := i1.serviceGet("interface {}")
	is.Nil(service1.(Service[any]).healthcheck(ctx))

	// healthcheck ok
	i2 := New()
	Provide(i2, func(i Injector) (*lazyTestHeathcheckerOK, error) {
		return &lazyTestHeathcheckerOK{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTestHeathcheckerOK, Healthchecker](i2))
	service2, _ := i2.serviceGet("github.com/sllt/af/di.Healthchecker")
	is.Nil(service2.(Service[Healthchecker]).healthcheck(ctx))
	_, _ = service2.(Service[Healthchecker]).getInstance(nil)
	is.Nil(service2.(Service[Healthchecker]).healthcheck(ctx))

	// healthcheck ko
	i3 := New()
	Provide(i3, func(i Injector) (*lazyTestHeathcheckerKO, error) {
		return &lazyTestHeathcheckerKO{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTestHeathcheckerKO, Healthchecker](i3))
	service3, _ := i3.serviceGet("github.com/sllt/af/di.Healthchecker")
	is.Nil(service3.(Service[Healthchecker]).healthcheck(ctx))
	_, _ = service3.(Service[Healthchecker]).getInstance(nil)
	is.Equal(assert.AnError, service3.(Service[Healthchecker]).healthcheck(ctx))

	// service not found (wrong type)
	i4 := New()
	service4 := newServiceAlias[*lazyTestHeathcheckerKO, Healthchecker]("github.com/sllt/af/di.Healthchecker", i4, "*github.com/sllt/af/di.lazyTestHeathcheckerKO")
	is.Nil(service4.healthcheck(ctx))
	_, _ = service4.getInstanceAny(nil)
	is.Nil(service4.healthcheck(ctx))

	// service not found (wrong name)
	i5 := New()
	service5 := newServiceAlias[*lazyTestHeathcheckerOK, Healthchecker]("github.com/sllt/af/di.Healthchecker", i5, "*github.com/sllt/af/di.lazyTestHeathcheckerKO")
	is.Nil(service5.healthcheck(ctx))
	_, _ = service5.getInstanceAny(nil)
	is.Nil(service5.healthcheck(ctx))
}

// @TODO: missing tests for context
func TestServiceAlias_isShutdowner(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	// no shutdown
	i1 := New()
	Provide(i1, func(i Injector) (*lazyTest, error) {
		return &lazyTest{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTest, any](i1))
	service1, _ := i1.serviceGet("interface {}")
	is.False(service1.(Service[any]).isShutdowner())

	// shutdown ok
	i2 := New()
	Provide(i2, func(i Injector) (*lazyTestShutdownerOK, error) {
		return &lazyTestShutdownerOK{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTestShutdownerOK, ShutdownerWithContextAndError](i2))
	service2, _ := i2.serviceGet("github.com/sllt/af/di.ShutdownerWithContextAndError")
	is.False(service2.(Service[ShutdownerWithContextAndError]).isShutdowner())
	_, _ = service2.(Service[ShutdownerWithContextAndError]).getInstance(nil)
	is.True(service2.(Service[ShutdownerWithContextAndError]).isShutdowner())

	// shutdown ko
	i3 := New()
	Provide(i3, func(i Injector) (*lazyTestShutdownerKO, error) {
		return &lazyTestShutdownerKO{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTestShutdownerKO, ShutdownerWithError](i3))
	service3, _ := i3.serviceGet("github.com/sllt/af/di.ShutdownerWithError")
	is.False(service3.(Service[ShutdownerWithError]).isShutdowner())
	_, _ = service3.(Service[ShutdownerWithError]).getInstance(nil)
	is.True(service3.(Service[ShutdownerWithError]).isShutdowner())

	// service not found (wrong type)
	i4 := New()
	service4 := newServiceAlias[*lazyTestShutdownerKO, Healthchecker]("*github.com/sllt/af/di.Healthchecker", i4, "*github.com/sllt/af/di.lazyTestShutdownerKO")
	is.False(service4.isShutdowner())
	_, _ = service4.getInstanceAny(nil)
	is.False(service4.isShutdowner())

	// service not found (wrong name)
	i5 := New()
	service5 := newServiceAlias[*lazyTestShutdownerOK, Healthchecker]("*github.com/sllt/af/di.Healthchecker", i5, "*github.com/sllt/af/di.lazyTestShutdownerKO")
	is.False(service5.isShutdowner())
	_, _ = service5.getInstanceAny(nil)
	is.False(service5.isShutdowner())
}

func TestServiceAlias_shutdown(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	ctx := context.Background()

	// no shutdown
	i1 := New()
	Provide(i1, func(i Injector) (*lazyTest, error) {
		return &lazyTest{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTest, any](i1))
	service1, _ := i1.serviceGet("interface {}")
	is.Nil(service1.(Service[any]).shutdown(ctx))

	// shutdown ok
	i2 := New()
	Provide(i2, func(i Injector) (*lazyTestShutdownerOK, error) {
		return &lazyTestShutdownerOK{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTestShutdownerOK, ShutdownerWithContextAndError](i2))
	service2, _ := i2.serviceGet("github.com/sllt/af/di.ShutdownerWithContextAndError")
	is.Nil(service2.(Service[ShutdownerWithContextAndError]).shutdown(ctx))
	_, _ = service2.(Service[ShutdownerWithContextAndError]).getInstance(nil)
	is.Nil(service2.(Service[ShutdownerWithContextAndError]).shutdown(ctx))

	// shutdown ko
	i3 := New()
	Provide(i3, func(i Injector) (*lazyTestShutdownerKO, error) {
		return &lazyTestShutdownerKO{foobar: "foobar"}, nil
	})
	is.Nil(As[*lazyTestShutdownerKO, ShutdownerWithError](i3))
	service3, _ := i3.serviceGet("github.com/sllt/af/di.ShutdownerWithError")
	is.Nil(service3.(Service[ShutdownerWithError]).shutdown(ctx))
	_, _ = service3.(Service[ShutdownerWithError]).getInstance(nil)
	is.Equal(assert.AnError, service3.(Service[ShutdownerWithError]).shutdown(ctx))

	// service not found (wrong type)
	i4 := New()
	service4 := newServiceAlias[*lazyTestShutdownerKO, Healthchecker]("github.com/sllt/af/di.Healthchecker", i4, "*github.com/sllt/af/di.lazyTestShutdownerKO")
	is.Nil(service4.shutdown(ctx))
	_, _ = service4.getInstanceAny(nil)
	is.Nil(service4.shutdown(ctx))

	// service not found (wrong name)
	i5 := New()
	service5 := newServiceAlias[*lazyTestShutdownerOK, Healthchecker]("github.com/sllt/af/di.Healthchecker", i5, "*github.com/sllt/af/di.lazyTestHeathcheckerKO")
	is.Nil(service5.shutdown(ctx))
	_, _ = service5.getInstanceAny(nil)
	is.Nil(service5.shutdown(ctx))
}

func TestServiceAlias_clone(t *testing.T) {
	// @TODO
}

func TestServiceAlias_source(t *testing.T) {
	// @TODO
}
