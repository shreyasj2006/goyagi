package movies

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/shreyasj2006/goyagi/internal/test"
	"github.com/shreyasj2006/goyagi/pkg/application"
	"github.com/shreyasj2006/goyagi/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateHandler(t *testing.T) {
	h := newHandler(t)

	t.Run("creates a new movie", func(tt *testing.T) {
		releaseDate := time.Now()
		payload := createParams{
			Title:       "a new movie",
			ReleaseDate: releaseDate,
		}
		payloadBytes, err := json.Marshal(payload)
		assert.NoError(tt, err)
		c, rr := test.NewContext(tt, payloadBytes, echo.MIMEApplicationJSON, "")

		err = h.createHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		var movie model.Movie
		err = json.Unmarshal(rr.Body.Bytes(), &movie)
		require.NoError(tt, err)
		assert.Equal(tt, movie.Title, "a new movie")
		// If you test with assert.Equal, these times won't be equal since
		// there will be a difference in monotonic clock.
		// For more info, refer https://golang.org/pkg/time/#hdr-Monotonic_Clocks
		assert.True(tt, movie.ReleaseDate.Equal(releaseDate))
	})
}

func TestListHandler(t *testing.T) {
	h := newHandler(t)

	t.Run("lists movies on success", func(tt *testing.T) {
		c, rr := test.NewContext(tt, nil, echo.MIMEApplicationJSON, "")

		err := h.listHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		var response []model.Movie
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(tt, err)
		assert.True(tt, len(response) >= 23)
	})

	t.Run("tests list params", func(tt *testing.T) {
		c, rr := test.NewContext(tt, nil, echo.MIMEApplicationJSON, "?limit=2")
		err := h.listHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		var response []model.Movie
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(tt, err)
		assert.True(tt, len(response) == 2)
	})
}

func TestRetrieveHandler(t *testing.T) {
	h := newHandler(t)

	t.Run("retrieves movie on success", func(tt *testing.T) {
		c, rr := test.NewContext(tt, nil, echo.MIMEApplicationJSON, "")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.retrieveHandler(c)
		assert.NoError(tt, err)
		assert.Equal(tt, http.StatusOK, rr.Code)

		var response model.Movie
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(tt, err)
		assert.Equal(tt, response.ID, 1)
		assert.Equal(tt, response.Title, "Iron Man")
	})

	t.Run("returns 404 if user isn't found", func(tt *testing.T) {
		c, _ := test.NewContext(tt, nil, echo.MIMEApplicationJSON, "")
		c.SetParamNames("id")
		c.SetParamValues("9999")

		err := h.retrieveHandler(c)
		assert.Contains(tt, err.Error(), "movie not found")
	})
}

// newHandler returns a new handler to be used for tests.
func newHandler(t *testing.T) handler {
	t.Helper()

	app, err := application.New()
	require.NoError(t, err)
	return handler{app}
}
