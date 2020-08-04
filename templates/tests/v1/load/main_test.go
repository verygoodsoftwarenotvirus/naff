package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_mainDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := mainDotGo(proj)

		expected := `
package example

import (
	"bytes"
	"context"
	hazana "github.com/emicklei/hazana"
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// ServiceAttacker implements hazana's Attacker interface.
type ServiceAttacker struct {
	todoClient *http.V1Client
}

// Setup implements hazana's Attacker interface.
func (a *ServiceAttacker) Setup(_ hazana.Config) error {
	return nil
}

// Do implements hazana's Attacker interface.
func (a *ServiceAttacker) Do(_ context.Context) hazana.DoResult {
	// Do performs one request and is executed in a separate goroutine.
	// The context is used to cancel the request on timeout.
	act := RandomAction(a.todoClient)

	req, err := act.Action()
	if err != nil || req == nil {
		if err == ErrUnavailableYet {
			return hazana.DoResult{
				RequestLabel: act.Name,
				Error:        nil,
				StatusCode:   200,
			}
		}
		log.Printf("something has gone awry: %v\n", err)
		return hazana.DoResult{Error: err}
	}

	var (
		sc int
		bo int64
		bi []byte
	)
	if req.Body != nil {
		bi, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return hazana.DoResult{Error: err}
		}
		rdr := ioutil.NopCloser(bytes.NewBuffer(bi))
		req.Body = rdr
	}

	res, err := a.todoClient.AuthenticatedClient().Do(req)
	if res != nil {
		sc = res.StatusCode
		bo = res.ContentLength
	}

	dr := hazana.DoResult{
		RequestLabel: act.Name,
		Error:        err,
		StatusCode:   sc,
		BytesIn:      int64(len(bi)),
		BytesOut:     bo,
	}
	return dr
}

// Teardown implements hazana's Attacker interface.
func (a *ServiceAttacker) Teardown() error {
	return nil
}

// Clone implements hazana's Attacker interface.
func (a *ServiceAttacker) Clone() hazana.Attack {
	return a
}

func main() {
	todoClient := initializeClient(oa2Client)

	var runTime = 10 * time.Minute
	if rt := os.Getenv("LOADTEST_RUN_TIME"); rt != "" {
		_rt, err := time.ParseDuration(rt)
		if err != nil {
			panic(err)
		}
		runTime = _rt
	}

	attacker := &ServiceAttacker{todoClient: todoClient}
	cfg := hazana.Config{
		RPS:           50,
		AttackTimeSec: int(runTime.Seconds()),
		RampupTimeSec: 5,
		MaxAttackers:  50,
		Verbose:       true,
		DoTimeoutSec:  10,
	}

	r := hazana.Run(attacker, cfg)

	// inspect the report and compute whether the test has failed
	// e.g by looking at the success percentage and mean response time of each metric.
	r.Failed = false

	hazana.PrintReport(r)
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMainServiceAttacker(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildMainServiceAttacker(proj)

		expected := `
package example

import (
	http "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/client/v1/http"
)

// ServiceAttacker implements hazana's Attacker interface.
type ServiceAttacker struct {
	todoClient *http.V1Client
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMainSetup(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMainSetup()

		expected := `
package example

import (
	hazana "github.com/emicklei/hazana"
)

// Setup implements hazana's Attacker interface.
func (a *ServiceAttacker) Setup(_ hazana.Config) error {
	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMainDo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMainDo()

		expected := `
package example

import (
	"bytes"
	"context"
	hazana "github.com/emicklei/hazana"
	"io/ioutil"
	"log"
)

// Do implements hazana's Attacker interface.
func (a *ServiceAttacker) Do(_ context.Context) hazana.DoResult {
	// Do performs one request and is executed in a separate goroutine.
	// The context is used to cancel the request on timeout.
	act := RandomAction(a.todoClient)

	req, err := act.Action()
	if err != nil || req == nil {
		if err == ErrUnavailableYet {
			return hazana.DoResult{
				RequestLabel: act.Name,
				Error:        nil,
				StatusCode:   200,
			}
		}
		log.Printf("something has gone awry: %v\n", err)
		return hazana.DoResult{Error: err}
	}

	var (
		sc int
		bo int64
		bi []byte
	)
	if req.Body != nil {
		bi, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return hazana.DoResult{Error: err}
		}
		rdr := ioutil.NopCloser(bytes.NewBuffer(bi))
		req.Body = rdr
	}

	res, err := a.todoClient.AuthenticatedClient().Do(req)
	if res != nil {
		sc = res.StatusCode
		bo = res.ContentLength
	}

	dr := hazana.DoResult{
		RequestLabel: act.Name,
		Error:        err,
		StatusCode:   sc,
		BytesIn:      int64(len(bi)),
		BytesOut:     bo,
	}
	return dr
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMainTeardown(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMainTeardown()

		expected := `
package example

import ()

// Teardown implements hazana's Attacker interface.
func (a *ServiceAttacker) Teardown() error {
	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMainClone(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMainClone()

		expected := `
package example

import (
	hazana "github.com/emicklei/hazana"
)

// Clone implements hazana's Attacker interface.
func (a *ServiceAttacker) Clone() hazana.Attack {
	return a
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMainMain(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMainMain()

		expected := `
package example

import (
	hazana "github.com/emicklei/hazana"
	"os"
	"time"
)

func main() {
	todoClient := initializeClient(oa2Client)

	var runTime = 10 * time.Minute
	if rt := os.Getenv("LOADTEST_RUN_TIME"); rt != "" {
		_rt, err := time.ParseDuration(rt)
		if err != nil {
			panic(err)
		}
		runTime = _rt
	}

	attacker := &ServiceAttacker{todoClient: todoClient}
	cfg := hazana.Config{
		RPS:           50,
		AttackTimeSec: int(runTime.Seconds()),
		RampupTimeSec: 5,
		MaxAttackers:  50,
		Verbose:       true,
		DoTimeoutSec:  10,
	}

	r := hazana.Run(attacker, cfg)

	// inspect the report and compute whether the test has failed
	// e.g by looking at the success percentage and mean response time of each metric.
	r.Failed = false

	hazana.PrintReport(r)
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
