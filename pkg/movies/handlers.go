package movies

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"github.com/shreyasj2006/goyagi/pkg/application"
	"github.com/shreyasj2006/goyagi/pkg/model"
)

const (
	resultError   = "result:error"
	resultSuccess = "result:success"

	timerPrefix = "goyagi.movies"
)

type handler struct {
	app application.App
}

func (h *handler) createHandler(c echo.Context) error {
	// params is a struct that will have our payload bound and validated against
	params := createParams{}
	if err := c.Bind(&params); err != nil {
		// if there is an error binding or validating the payload, return early with an error
		return err
	}

	movie := model.Movie{
		Title:       params.Title,
		ReleaseDate: params.ReleaseDate,
	}

	insertTimer := h.app.Metrics.NewTimer(fmt.Sprintf("%s.create.db", timerPrefix))
	_, err := h.app.DB.Model(&movie).Insert()
	if err != nil {
		insertTimer.End(resultError)
		return err
	}
	insertTimer.End(resultSuccess)

	return c.JSON(http.StatusOK, movie)
}

func (h *handler) listHandler(c echo.Context) error {
	params := listParams{}
	if err := c.Bind(&params); err != nil {
		return err
	}

	var movies []*model.Movie

	selectTimer := h.app.Metrics.NewTimer(fmt.Sprintf("%s.select.db", timerPrefix))
	err := h.app.DB.
		Model(&movies).
		Limit(params.Limit).
		Offset(params.Offset).
		Order("id DESC").
		Select()
	if err != nil {
		selectTimer.End(resultError)
		return err
	}
	selectTimer.End(resultSuccess)

	return c.JSON(http.StatusOK, movies)
}

func (h *handler) retrieveHandler(c echo.Context) error {
	id := c.Param("id")

	var movie model.Movie

	selectTimer := h.app.Metrics.NewTimer(fmt.Sprintf("%s.select.db", timerPrefix))
	err := h.app.DB.Model(&movie).Where("id = ?", id).First()
	if err != nil {
		selectTimer.End(resultError)
		if err == pg.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "movie not found")
		}
		return err
	}
	selectTimer.End(resultSuccess)

	return c.JSON(http.StatusOK, movie)
}
