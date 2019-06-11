package recovery

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// Middleware recovers from any possible panics in subsequent handlers
// and funnels it to the error handler to be returned as a 500.
func Middleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				// If for some reason we panic we can recover from it and
				// invoke our error handler to handle the error surfaced by
				// the panic.
				if r := recover(); r != nil {
					// Create an error based on the recovery panic message
					err := errors.Errorf("%v", r)
					// Invoke our error handler with the error created above
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
