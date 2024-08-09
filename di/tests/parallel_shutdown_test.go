package tests

import (
	"github.com/sllt/af/di"
	"github.com/sllt/af/di/tests/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParallelShutdown(t *testing.T) {
	is := assert.New(t)

	root, driver, passenger := fixtures.GetPackage()
	is.NotPanics(func() {
		_ = di.MustInvoke[*fixtures.Driver](driver)
		_ = di.MustInvokeNamed[*fixtures.Passenger](passenger, "passenger-1")
		_ = di.MustInvokeNamed[*fixtures.Passenger](passenger, "passenger-2")
		_ = di.MustInvokeNamed[*fixtures.Passenger](passenger, "passenger-3")
		root.Shutdown()
	})
}
