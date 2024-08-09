package fixtures

import (
	"github.com/sllt/af/di"
)

/**
 * driver
 */
func newDriver(i di.Injector) (*Driver, error) {
	d := &Driver{
		seat:   di.MustInvokeNamed[*Seat](i, "seat-1"),
		engine: di.MustInvoke[*Engine](i),
	}

	d.seat.take()
	d.engine.start()

	return d, nil
}

type Driver struct {
	seat   *Seat
	engine *Engine
}

func (d *Driver) Shutdiwn() {
	d.engine.stop()
	d.seat.release()
}

/**
 * passenger
 */
func newPassenger(seatName string) func(i di.Injector) (*Passenger, error) {
	return func(i di.Injector) (*Passenger, error) {
		p := &Passenger{
			seat: di.MustInvokeNamed[*Seat](i, seatName),
		}

		p.seat.take()

		return p, nil
	}
}

type Passenger struct {
	seat *Seat
}

func (p *Passenger) Shutdown() {
	p.seat.release()
}

/**
 * Seat
 */
type Seat struct {
	busy bool
}

func (s *Seat) take() {
	if s.busy {
		panic("seat should be free")
	}
	s.busy = true
}

func (s *Seat) release() {
	if !s.busy {
		panic("seat should be busy")
	}
	s.busy = false
}

func (s *Seat) Shutdiwn() {
	if s.busy {
		panic("seat should be free")
	}
}

/**
 * Wheel
 */
type Wheel struct {
	active bool
}

func (w *Wheel) start() {
	if w.active {
		panic("wheel should be stopped")
	}
	w.active = true
}

func (w *Wheel) stop() {
	if !w.active {
		panic("wheel should be started")
	}
	w.active = false
}

func (w *Wheel) Shutdiwn() {
	if w.active {
		panic("wheel should be stopped")
	}
}

/**
 * engine
 */
func newEngine(i di.Injector) (*Engine, error) {
	return &Engine{
		wheels: []*Wheel{
			di.MustInvokeNamed[*Wheel](i, "wheel-1"),
			di.MustInvokeNamed[*Wheel](i, "wheel-2"),
			di.MustInvokeNamed[*Wheel](i, "wheel-3"),
			di.MustInvokeNamed[*Wheel](i, "wheel-4"),
		},
	}, nil
}

type Engine struct {
	started bool
	wheels  []*Wheel
}

func (e *Engine) start() {
	if e.started {
		panic("engine should be stopped")
	}
	e.started = true

	for _, wheel := range e.wheels {
		wheel.start()
	}
}

func (e *Engine) stop() {
	if !e.started {
		panic("engine should be started")
	}
	e.started = false

	for _, wheel := range e.wheels {
		wheel.stop()
	}
}

func (e *Engine) Shutdiwn() {
	if e.started {
		panic("engine should be stopped")
	}
}

func GetPackage() (di.Injector, di.Injector, di.Injector) {
	injector := di.New()

	driverScope := injector.Scope("driver")
	passengerScope := injector.Scope("passenger")

	// provide wheels
	di.ProvideNamedValue(injector, "wheel-1", &Wheel{})
	di.ProvideNamedValue(injector, "wheel-2", &Wheel{})
	di.ProvideNamedValue(injector, "wheel-3", &Wheel{})
	di.ProvideNamedValue(injector, "wheel-4", &Wheel{})

	// provide engine
	di.Provide(injector, newEngine)

	// provide seats
	di.ProvideNamedValue(injector, "seat-1", &Seat{})
	di.ProvideNamedValue(injector, "seat-2", &Seat{})
	di.ProvideNamedValue(injector, "seat-3", &Seat{})
	di.ProvideNamedValue(injector, "seat-4", &Seat{})

	// provide driver
	di.Provide(driverScope, newDriver)

	// provide passenger
	di.ProvideNamed(passengerScope, "passenger-1", newPassenger("seat-2"))
	di.ProvideNamed(passengerScope, "passenger-2", newPassenger("seat-3"))
	di.ProvideNamed(passengerScope, "passenger-3", newPassenger("seat-4"))

	return injector, driverScope, passengerScope
}
