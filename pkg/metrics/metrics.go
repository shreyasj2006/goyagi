package metrics

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/labstack/echo"
	"github.com/shreyasj2006/goyagi/pkg/config"
)

type statsdClient interface {
	Histogram(name string, value float64, tags []string, rate float64) error
	Count(name string, value int64, tags []string, rate float64) error
}

// Metrics functions for metrics clients.
type Metrics struct {
	client statsdClient
}

// Timer facilitates timing and tagging a metric before sending it to Datadog
// as a Histogram datapoint.
type Timer struct {
	name    string
	metrics *Metrics
	begin   time.Time
	tags    []string
}

const namespace = "goyagi."

// New sets up metric package with a Datadog client.
func New(cfg config.Config) (Metrics, error) {
	address := fmt.Sprintf("%s:%d", cfg.StatsdHost, cfg.StatsdPort)

	client, err := statsd.New(address)
	if err != nil {
		return Metrics{}, err
	}

	client.Namespace = namespace
	client.Tags = []string{
		fmt.Sprintf("environment:%s", cfg.Environment),
	}

	return Metrics{client}, nil
}

// Count increments an event counter in Datadog while disregarding potential
// errors.
func (m *Metrics) Count(name string, count int64, tags ...string) {
	m.client.Count(name, count, tags, 1) // nolint:gosec
}

// Histogram sends statistical distribution data to Datadog while disregarding
// potential errors.
func (m *Metrics) Histogram(name string, value float64, tags ...string) {
	m.client.Histogram(name, value, tags, 1) // nolint:gosec
}

// NewTimer returns a Timer object with a set start time
func (m *Metrics) NewTimer(name string, tags ...string) Timer {
	return Timer{
		begin:   time.Now(),
		metrics: m,
		name:    name,
		tags:    tags,
	}
}

// End ends a Timer and sends the metric and duration to Datadog as a
// Histogram datapoint.
func (t *Timer) End(additionalTags ...string) float64 {
	duration := time.Since(t.begin)
	durationInMS := float64(duration / time.Millisecond)

	t.tags = append(t.tags, additionalTags...)

	t.metrics.Histogram(t.name, durationInMS, t.tags...)

	return durationInMS
}

// Middleware returns an Echo middleware function that begins a timer before a
// request is handled and ends afterwards.
func Middleware(m Metrics) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			methodTag := fmt.Sprintf("method:%s", c.Request().Method)

			// Create a new timer
			t := m.NewTimer("http.request", methodTag)

			// Continue the execution down the middleware/handler stack
			if err := next(c); err != nil {
				c.Error(err)
			}

			statusCodeTag := fmt.Sprintf("status_code:%d", c.Response().Status)
			pathTag := fmt.Sprintf("path:%s", c.Path())

			// End the timer once we have succeeded calling the middleware/handler
			// stack. This function call emits the metric to datadog.
			t.End(statusCodeTag, pathTag)

			return nil
		}
	}
}
