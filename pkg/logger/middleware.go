package logger

import (
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo"
	"github.com/lob/logger-go"
	"github.com/pkg/errors"
)

const key = "logger"

// Middleware attaches a Logger instance with a request ID onto the context. It
// also logs every request along with metadata about the request.
func Middleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	l := logger.New()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Record the time start time of the middleware invocation
			t1 := time.Now()

			// Generate a new UUID that will be used to recognize this particular
			// request
			id, err := uuid.NewV4()
			if err != nil {
				return errors.WithStack(err)
			}

			// Create a child logger with the unique UUID created and attach it to
			// the echo.Context. By attaching it to the context it can be fetched by
			// later middleware or handler functions to emit events with a logger
			// that contains this ID. This is useful as it allows us to emit all
			// events with the same request UUID.
			log := l.ID(id.String())
			c.Set(key, log)

			// Execute the next middleware/handler function in the stack.
			if err := next(c); err != nil {
				c.Error(err)
			}

			// We have now succeeded executing all later middlewares in the stack and
			// have come back to the logger middleware. Record the time at which we
			// came back to this middleware. We can use the difference between t2 and
			// t1 to calculate the request duration.
			t2 := time.Now()

			// Get the request IP address.
			var ipAddress string
			if xff := c.Request().Header.Get("x-forwarded-for"); xff != "" {
				split := strings.Split(xff, ",")
				ipAddress = strings.TrimSpace(split[len(split)-1])
			} else {
				ipAddress = c.Request().RemoteAddr
			}

			// Emit a log event with as much metadata as we can.
			log.Root(logger.Data{
				"status_code":   c.Response().Status,
				"method":        c.Request().Method,
				"path":          c.Request().URL.Path,
				"route":         c.Path(),
				"response_time": t2.Sub(t1).Seconds() * 1000,
				"referer":       c.Request().Referer(),
				"user_agent":    c.Request().UserAgent(),
				"ip_address":    ipAddress,
				"trace_id":      c.Request().Header.Get("x-amzn-trace-id"),
			}).Info("request handled")

			// Succeeded executing the middleware invocation. A nil response
			// represents no errors happened.
			return nil
		}
	}
}

// FromContext returns a Logger from the given echo.Context. If there is no
// attached logger, then it will return a new Logger instance.
func FromContext(c echo.Context) logger.Logger {
	if log, ok := c.Get(key).(logger.Logger); ok {
		return log
	}

	return logger.New()
}
