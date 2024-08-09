package do

import (
	"github.com/sllt/af/di"
)

var DefaultInjector = di.New()

type Provider[T any] func() (T, error)

func warpProvider[T any](provider Provider[T]) di.Provider[T] {
	return func(i di.Injector) (T, error) {
		return provider()
	}
}

func New(packages ...func(di.Injector)) *di.RootScope {
	return di.NewWithOpts(&di.InjectorOpts{}, packages...)
}

func NameOf[T any]() string {
	return di.NameOf[T]()
}

func Provide[T any](provider Provider[T]) {
	di.Provide[T](DefaultInjector, warpProvider[T](provider))
}

func ProvideNamed[T any](name string, provider Provider[T]) {
	di.ProvideNamed[T](DefaultInjector, name, warpProvider[T](provider))
}

func ProvideValue[T any](value T) {
	di.ProvideValue[T](DefaultInjector, value)
}

func ProvideNamedValue[T any](name string, value T) {
	di.ProvideNamedValue[T](DefaultInjector, name, value)
}

func ProvideTransient[T any](provider Provider[T]) {
	di.ProvideTransient[T](DefaultInjector, warpProvider[T](provider))
}

func ProvideNamedTransient[T any](name string, provider Provider[T]) {
	di.ProvideNamedTransient[T](DefaultInjector, name, warpProvider[T](provider))
}

func Override[T any](provider Provider[T]) {
	di.Override[T](DefaultInjector, warpProvider[T](provider))
}

func OverrideNamed[T any](name string, provider Provider[T]) {
	di.OverrideNamed[T](DefaultInjector, name, warpProvider[T](provider))
}

func OverrideValue[T any](value T) {
	di.OverrideValue[T](DefaultInjector, value)
}

func OverrideNamedValue[T any](name string, value T) {
	di.OverrideNamedValue[T](DefaultInjector, name, value)
}

func OverrideTransient[T any](provider Provider[T]) {
	di.OverrideTransient[T](DefaultInjector, warpProvider[T](provider))
}

func OverrideNamedTransient[T any](name string, provider Provider[T]) {
	di.OverrideNamedTransient[T](DefaultInjector, name, warpProvider[T](provider))
}

func Invoke[T any]() (T, error) {
	return di.Invoke[T](DefaultInjector)
}

// MustInvoke invokes a service in the DI container, using type inference to determine the service name. It panics on error.
func MustInvoke[T any]() T {
	return di.MustInvoke[T](DefaultInjector)
}

func InvokeNamed[T any](name string) (T, error) {
	return di.InvokeNamed[T](DefaultInjector, name)
}

func MustInvokeNamed[T any](name string) T {
	return di.MustInvokeNamed[T](DefaultInjector, name)
}

func InvokeStruct[T any]() (*T, error) {
	return di.InvokeStruct[T](DefaultInjector)
}

func MustInvokeStruct[T any]() *T {
	return di.MustInvokeStruct[T](DefaultInjector)
}

func Package(services ...func(i di.Injector)) func(di.Injector) {
	return func(injector di.Injector) {
		for i := range services {
			services[i](injector)
		}
	}
}

func Lazy[T any](p Provider[T]) func(di.Injector) {
	return func(injector di.Injector) {
		Provide(p)
	}
}

func LazyNamed[T any](serviceName string, p Provider[T]) func(di.Injector) {
	return func(injector di.Injector) {
		ProvideNamed(serviceName, p)
	}
}

func Eager[T any](value T) func(di.Injector) {
	return func(injector di.Injector) {
		ProvideValue(value)
	}
}

func EagerNamed[T any](serviceName string, value T) func(di.Injector) {
	return func(injector di.Injector) {
		ProvideNamedValue(serviceName, value)
	}
}

func Transient[T any](p Provider[T]) func(di.Injector) {
	return func(injector di.Injector) {
		ProvideTransient(p)
	}
}

func TransientNamed[T any](serviceName string, p Provider[T]) func(di.Injector) {
	return func(injector di.Injector) {
		ProvideNamedTransient(serviceName, p)
	}
}

func Bind[Initial any, Alias any]() func(di.Injector) {
	return func(injector di.Injector) {
		di.MustAs[Initial, Alias](injector)
	}
}

func BindNamed[Initial any, Alias any](initial string, alias string) func(di.Injector) {
	return func(injector di.Injector) {
		di.MustAsNamed[Initial, Alias](injector, initial, alias)
	}
}
