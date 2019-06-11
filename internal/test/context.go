package test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
)

// NewContext returns a new echo.Context, and *httptest.ResponseRecorder to be
// used for tests.
func NewContext(t *testing.T, payload []byte, mime string, queryParamString string) (echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	e := echo.New()
	endpoint := fmt.Sprintf("/%s", queryParamString)
	req := httptest.NewRequest(echo.GET, endpoint, bytes.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, mime)
	rr := httptest.NewRecorder()
	c := e.NewContext(req, rr)
	return c, rr
}
