package di

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sllt/af/di/stacktrace"
)

var _ Service[int] = (*serviceLazy[int])(nil)
var _ serviceHealthcheck = (*serviceLazy[int])(nil)
var _ serviceShutdown = (*serviceLazy[int])(nil)
var _ serviceClone = (*serviceLazy[int])(nil)

type serviceLazy[T any] struct {
	mu       sync.RWMutex
	name     string
	typeName string
	instance T

	// lazy loading
	built     bool
	buildTime time.Duration
	provider  Provider[T]

	providerFrame           stacktrace.Frame
	invokationFrames        map[stacktrace.Frame]struct{} // map garanties uniqueness
	invokationFramesCounter uint32
}

func newServiceLazy[T any](name string, provider Provider[T]) *serviceLazy[T] {
	providerFrame, _ := stacktrace.NewFrameFromPtr(reflect.ValueOf(provider).Pointer())

	return &serviceLazy[T]{
		mu:       sync.RWMutex{},
		name:     name,
		typeName: inferServiceName[T](),

		built:     false,
		buildTime: 0,
		provider:  provider,

		providerFrame:           providerFrame,
		invokationFrames:        map[stacktrace.Frame]struct{}{},
		invokationFramesCounter: 0,
	}
}

func (s *serviceLazy[T]) getName() string {
	return s.name
}

func (s *serviceLazy[T]) getTypeName() string {
	return s.typeName
}

func (s *serviceLazy[T]) getServiceType() ServiceType {
	return ServiceTypeLazy
}

func (s *serviceLazy[T]) getEmptyInstance() any {
	return empty[T]()
}

func (s *serviceLazy[T]) getInstanceAny(i Injector) (any, error) {
	return s.getInstance(i)
}

func (s *serviceLazy[T]) getInstance(i Injector) (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Collect up to 100 invokation frames.
	// In the future, we can implement a LFU list, to evict the oldest
	// frames and keep the most recent ones, but it would be much more costly.
	if atomic.AddUint32(&s.invokationFramesCounter, 1) < MaxInvokationFrames {
		frame, ok := stacktrace.NewFrameFromCaller()
		if ok {
			s.invokationFrames[frame] = struct{}{}
		}
	}

	if !s.built {
		err := s.build(i)
		if err != nil {
			return empty[T](), err
		}
	}

	return s.instance, nil
}

func (s *serviceLazy[T]) build(i Injector) (err error) {
	start := time.Now()

	instance, err := handleProviderPanic(s.provider, i)
	if err != nil {
		return err
	}

	s.instance = instance
	s.built = true
	s.buildTime = time.Since(start)

	return nil
}

func (s *serviceLazy[T]) isHealthchecker() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.built {
		return false
	}

	_, ok1 := any(s.instance).(HealthcheckerWithContext)
	_, ok2 := any(s.instance).(Healthchecker)
	return ok1 || ok2
}

func (s *serviceLazy[T]) healthcheck(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.built {
		return nil
	}

	if instance, ok := any(s.instance).(HealthcheckerWithContext); ok {
		return instance.HealthCheck(ctx)
	} else if instance, ok := any(s.instance).(Healthchecker); ok {
		return instance.HealthCheck()
	}

	return nil
}

func (s *serviceLazy[T]) isShutdowner() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.built {
		return false
	}

	_, ok1 := any(s.instance).(ShutdownerWithContextAndError)
	_, ok2 := any(s.instance).(ShutdownerWithError)
	_, ok3 := any(s.instance).(ShutdownerWithContext)
	_, ok4 := any(s.instance).(Shutdowner)
	return ok1 || ok2 || ok3 || ok4
}

func (s *serviceLazy[T]) shutdown(ctx context.Context) error {
	s.mu.Lock()

	defer func() {
		// whatever the outcome, reset `build` flag and instance
		s.built = false
		s.instance = empty[T]()
		s.mu.Unlock()
	}()

	if !s.built {
		return nil
	}

	if instance, ok := any(s.instance).(ShutdownerWithContextAndError); ok {
		return instance.Shutdown(ctx)
	} else if instance, ok := any(s.instance).(ShutdownerWithError); ok {
		return instance.Shutdown()
	} else if instance, ok := any(s.instance).(ShutdownerWithContext); ok {
		instance.Shutdown(ctx)
		return nil
	} else if instance, ok := any(s.instance).(Shutdowner); ok {
		instance.Shutdown()
		return nil
	}

	return nil
}

func (s *serviceLazy[T]) clone() any {
	// reset `build` flag and instance
	return &serviceLazy[T]{
		mu:       sync.RWMutex{},
		name:     s.name,
		typeName: s.typeName,

		built:     false,
		provider:  s.provider,
		buildTime: 0,

		providerFrame:           s.providerFrame,
		invokationFrames:        map[stacktrace.Frame]struct{}{},
		invokationFramesCounter: 0,
	}
}

//nolint:unused
func (s *serviceLazy[T]) source() (stacktrace.Frame, []stacktrace.Frame) {
	s.mu.RLock()
	invokationFrames := make([]stacktrace.Frame, 0, len(s.invokationFrames))
	for frame := range s.invokationFrames {
		invokationFrames = append(invokationFrames, frame)
	}
	s.mu.RUnlock()

	return s.providerFrame, invokationFrames
}

func (s *serviceLazy[T]) getBuildTime() (time.Duration, bool) {
	return s.buildTime, s.built
}
