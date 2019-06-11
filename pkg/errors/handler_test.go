package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	"github.com/shreyasj2006/goyagi/internal/test"
	"github.com/shreyasj2006/goyagi/pkg/application"
	"github.com/shreyasj2006/goyagi/pkg/sentry"
	"github.com/stretchr/testify/assert"
)

// mockSentryClient is a dummy struct that adheres to the ravenClient
// interface defined by in our sentry package. This allows us to create a mock
// client that has the same function signature but doesn't actually make a
// request to Sentry when invoked. While we do not assert the values that get
// passed into the function you can assert or fail tests based on the values
// that get pased in to the mock client functions.
type mockSentryClient struct{}

func (m mockSentryClient) Capture(packet *raven.Packet, captureTags map[string]string) (string, chan error) {
	return "", nil
}

func TestHandler(t *testing.T) {
	app := application.App{
		Sentry: sentry.Sentry{
			Client: mockSentryClient{}, // set the mock client instead of a real raven client
		},
	}
	h := handler{app}

	t.Run("surfaces generic errors as internal server errors", func(tt *testing.T) {
		c, rr := test.NewContext(t, nil, echo.MIMEApplicationJSON, "")
		err := errors.New("foo")

		h.handleError(err, c)

		assert.Equal(tt, http.StatusInternalServerError, rr.Code, "expected generic errors to be 500s")
		assert.Contains(tt, rr.Body.String(), "Internal Server Error", "expected generic errors to have the correct message")
	})

	t.Run("surfaces HTTP errors transparently but obfuscates message", func(tt *testing.T) {
		c, rr := test.NewContext(t, nil, echo.MIMEApplicationJSON, "")
		err := echo.NewHTTPError(http.StatusForbidden, "foo")

		h.handleError(err, c)

		assert.Equal(tt, http.StatusForbidden, rr.Code, "expected HTTP errors to be correct")
		assert.Contains(tt, rr.Body.String(), "Forbidden", "expected HTTP errors to have the correct message")
	})

	t.Run("overwrites HTTP 400 error messages", func(tt *testing.T) {
		c, rr := test.NewContext(t, nil, echo.MIMEApplicationJSON, "")
		err := echo.NewHTTPError(http.StatusBadRequest, "this shouldn't be sent to customers")

		h.handleError(err, c)

		assert.Equal(tt, http.StatusBadRequest, rr.Code, "expected HTTP errors to be correct")
		assert.Contains(tt, rr.Body.String(), "Bad Request", "expected HTTP errors to have the correct message")
	})

	t.Run("overwrites HTTP 403 error messages", func(tt *testing.T) {
		c, rr := test.NewContext(t, nil, echo.MIMEApplicationJSON, "")
		err := echo.NewHTTPError(http.StatusForbidden, "this shouldn't be sent to customers")

		h.handleError(err, c)

		assert.Equal(tt, http.StatusForbidden, rr.Code, "expected HTTP errors to be correct")
		assert.Contains(tt, rr.Body.String(), "Forbidden", "expected HTTP errors to have the correct message")
	})

	t.Run("overwrites HTTP 404 error messages", func(tt *testing.T) {
		c, rr := test.NewContext(t, nil, echo.MIMEApplicationJSON, "")
		err := echo.NewHTTPError(http.StatusNotFound, "this shouldn't be sent to customers")

		h.handleError(err, c)

		assert.Equal(tt, http.StatusNotFound, rr.Code, "expected HTTP errors to be correct")
		assert.Contains(tt, rr.Body.String(), "Not Found", "expected HTTP errors to have the correct message")
	})
}
