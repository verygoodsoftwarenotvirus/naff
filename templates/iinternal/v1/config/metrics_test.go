package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"
	"gitlab.com/verygoodsoftwarenotvirus/naff/testutils"
)

func Test_metricsDotGo(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := metricsDotGo(proj)

		expected := `
package example

import (
	jaeger "contrib.go.opencensus.io/exporter/jaeger"
	prometheus "contrib.go.opencensus.io/exporter/prometheus"
	"errors"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	view "go.opencensus.io/stats/view"
	trace "go.opencensus.io/trace"
	"math"
	"os"
	"time"
)

const (
	// MetricsNamespace is the namespace under which we register metrics.
	MetricsNamespace = "todo_server"

	// MinimumRuntimeCollectionInterval is the smallest interval we can collect metrics at
	// this value is used to guard against zero values.
	MinimumRuntimeCollectionInterval = time.Second
)

type (
	metricsProvider string
	tracingProvider string
)

var (
	// ErrInvalidMetricsProvider is a sentinel error value.
	ErrInvalidMetricsProvider = errors.New("invalid metrics provider")
	// Prometheus represents the popular time series database.
	Prometheus metricsProvider = "prometheus"
	// DefaultMetricsProvider indicates what the preferred metrics provider is.
	DefaultMetricsProvider = Prometheus

	// ErrInvalidTracingProvider is a sentinel error value.
	ErrInvalidTracingProvider = errors.New("invalid tracing provider")
	// Jaeger represents the popular distributed tracing server.
	Jaeger tracingProvider = "jaeger"
	// DefaultTracingProvider indicates what the preferred tracing provider is.
	DefaultTracingProvider = Jaeger
)

// ProvideInstrumentationHandler provides an instrumentation handler.
func (cfg *ServerConfig) ProvideInstrumentationHandler(logger v1.Logger) metrics.InstrumentationHandler {
	logger = logger.WithValue("metrics_provider", cfg.Metrics.MetricsProvider)
	logger.Debug("setting metrics provider")

	switch cfg.Metrics.MetricsProvider {
	case Prometheus:
		p, err := prometheus.NewExporter(
			prometheus.Options{
				OnError: func(err error) {
					logger.Error(err, "setting up prometheus export")
				},
				Namespace: MetricsNamespace,
			},
		)
		if err != nil {
			logger.Error(err, "failed to create Prometheus exporter")
			return nil
		}
		view.RegisterExporter(p)
		logger.Debug("metrics provider registered")

		if err := metrics.RegisterDefaultViews(); err != nil {
			logger.Error(err, "registering default metric views")
			return nil
		}
		metrics.RecordRuntimeStats(time.Duration(
			math.Max(
				float64(MinimumRuntimeCollectionInterval),
				float64(cfg.Metrics.RuntimeMetricsCollectionInterval),
			),
		))

		return p
	default:
		return nil
	}
}

// ProvideTracing provides an instrumentation handler.
func (cfg *ServerConfig) ProvideTracing(logger v1.Logger) error {
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1)})

	log := logger.WithValue("tracing_provider", cfg.Metrics.TracingProvider)
	log.Info("setting tracing provider")

	switch cfg.Metrics.TracingProvider {
	case Jaeger:
		ah := os.Getenv("JAEGER_AGENT_HOST")
		ap := os.Getenv("JAEGER_AGENT_PORT")
		sn := os.Getenv("JAEGER_SERVICE_NAME")

		if ah != "" && ap != "" && sn != "" {
			je, err := jaeger.NewExporter(jaeger.Options{
				AgentEndpoint: fmt.Sprintf("%s:%s", ah, ap),
				Process:       jaeger.Process{ServiceName: sn},
			})
			if err != nil {
				return fmt.Errorf("failed to create Jaeger exporter: %w", err)
			}

			trace.RegisterExporter(je)
			log.Debug("tracing provider registered")
		}
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMetricsConstantsDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMetricsConstantsDeclarations()

		expected := `
package example

import (
	"time"
)

const (
	// MetricsNamespace is the namespace under which we register metrics.
	MetricsNamespace = "todo_server"

	// MinimumRuntimeCollectionInterval is the smallest interval we can collect metrics at
	// this value is used to guard against zero values.
	MinimumRuntimeCollectionInterval = time.Second
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMetricsTypeDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMetricsTypeDeclarations()

		expected := `
package example

import ()

type (
	metricsProvider string
	tracingProvider string
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildMetricsVarDeclarations(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildMetricsVarDeclarations()

		expected := `
package example

import (
	"errors"
)

var (
	// ErrInvalidMetricsProvider is a sentinel error value.
	ErrInvalidMetricsProvider = errors.New("invalid metrics provider")
	// Prometheus represents the popular time series database.
	Prometheus metricsProvider = "prometheus"
	// DefaultMetricsProvider indicates what the preferred metrics provider is.
	DefaultMetricsProvider = Prometheus

	// ErrInvalidTracingProvider is a sentinel error value.
	ErrInvalidTracingProvider = errors.New("invalid tracing provider")
	// Jaeger represents the popular distributed tracing server.
	Jaeger tracingProvider = "jaeger"
	// DefaultTracingProvider indicates what the preferred tracing provider is.
	DefaultTracingProvider = Jaeger
)
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideInstrumentationHandler(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		x := buildProvideInstrumentationHandler(proj)

		expected := `
package example

import (
	prometheus "contrib.go.opencensus.io/exporter/prometheus"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	metrics "gitlab.com/verygoodsoftwarenotvirus/naff/example_output/internal/v1/metrics"
	view "go.opencensus.io/stats/view"
	"math"
	"time"
)

// ProvideInstrumentationHandler provides an instrumentation handler.
func (cfg *ServerConfig) ProvideInstrumentationHandler(logger v1.Logger) metrics.InstrumentationHandler {
	logger = logger.WithValue("metrics_provider", cfg.Metrics.MetricsProvider)
	logger.Debug("setting metrics provider")

	switch cfg.Metrics.MetricsProvider {
	case Prometheus:
		p, err := prometheus.NewExporter(
			prometheus.Options{
				OnError: func(err error) {
					logger.Error(err, "setting up prometheus export")
				},
				Namespace: MetricsNamespace,
			},
		)
		if err != nil {
			logger.Error(err, "failed to create Prometheus exporter")
			return nil
		}
		view.RegisterExporter(p)
		logger.Debug("metrics provider registered")

		if err := metrics.RegisterDefaultViews(); err != nil {
			logger.Error(err, "registering default metric views")
			return nil
		}
		metrics.RecordRuntimeStats(time.Duration(
			math.Max(
				float64(MinimumRuntimeCollectionInterval),
				float64(cfg.Metrics.RuntimeMetricsCollectionInterval),
			),
		))

		return p
	default:
		return nil
	}
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}

func Test_buildProvideTracing(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		x := buildProvideTracing()

		expected := `
package example

import (
	jaeger "contrib.go.opencensus.io/exporter/jaeger"
	"fmt"
	v1 "gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	trace "go.opencensus.io/trace"
	"os"
)

// ProvideTracing provides an instrumentation handler.
func (cfg *ServerConfig) ProvideTracing(logger v1.Logger) error {
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1)})

	log := logger.WithValue("tracing_provider", cfg.Metrics.TracingProvider)
	log.Info("setting tracing provider")

	switch cfg.Metrics.TracingProvider {
	case Jaeger:
		ah := os.Getenv("JAEGER_AGENT_HOST")
		ap := os.Getenv("JAEGER_AGENT_PORT")
		sn := os.Getenv("JAEGER_SERVICE_NAME")

		if ah != "" && ap != "" && sn != "" {
			je, err := jaeger.NewExporter(jaeger.Options{
				AgentEndpoint: fmt.Sprintf("%s:%s", ah, ap),
				Process:       jaeger.Process{ServiceName: sn},
			})
			if err != nil {
				return fmt.Errorf("failed to create Jaeger exporter: %w", err)
			}

			trace.RegisterExporter(je)
			log.Debug("tracing provider registered")
		}
	}

	return nil
}
`
		actual := testutils.RenderOuterStatementToString(t, x...)

		assert.Equal(t, expected, actual, "expected and actual output do not match")
	})
}
