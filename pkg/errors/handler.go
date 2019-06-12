package errors

import (
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	loggergo "github.com/lob/logger-go"
	"github.com/shreyasj2006/goyagi/pkg/application"
	"github.com/shreyasj2006/goyagi/pkg/logger"
)

type handler struct {
	app application.App
}

// RegisterErrorHandler takes in an Echo router and registers routes onto it.
func RegisterErrorHandler(e *echo.Echo, app application.App) {
	h := handler{app}

	e.HTTPErrorHandler = h.handleError
}

// handleError is an Echo error handler that uses HTTP errors accordingly, and any
// generic error will be interpreted as an internal server error.
func (h *handler) handleError(err error, c echo.Context) {
	// Fetch the logger set into the echo.Context by the logger middleware.
	// This logger has the request ID associated with this particular request.
	log := logger.FromContext(c)

	// Default our status code and message to a 500 error.
	code := http.StatusInternalServerError
	msg := http.StatusText(code)

	// If we are bubbling up an HTTP error it must be because we specifically
	// wanted to return this particular error and message. Set our status code
	// and message to the ones specificed in the HTTP error.
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = http.StatusText(code)
	}

	if code == http.StatusInternalServerError {
		stacktrace := raven.NewException(err, raven.GetOrNewStacktrace(err, 0, 2, nil))
		httpContext := raven.NewHttp(c.Request())
		packet := raven.NewPacket(msg, stacktrace, httpContext)

		h.app.Sentry.Client.Capture(packet, map[string]string{})
	}

	// Log our error with our child logger.
	log.Root(loggergo.Data{"status_code": code}).Err(err).Error("request error")

	// Return the error in the form of an HTTP response.
	err = c.JSON(code, map[string]interface{}{"error": map[string]interface{}{"message": msg, "status_code": code}})
	if err != nil {
		log.Err(err).Error("error handler json error")
	}
}
