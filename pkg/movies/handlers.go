package movies

import (
	"net/http"

	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"github.com/shreyasj2006/goyagi/pkg/application"
	"github.com/shreyasj2006/goyagi/pkg/model"
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

	insertTimer := h.app.Metrics.NewTimer("goyagi.movies.create.db")
	_, err := h.app.DB.Model(&movie).Insert()
	if err != nil {
		insertTimer.End("result:error")
		return err
	}
	insertTimer.End("result:success")

	return c.JSON(http.StatusOK, movie)
}

func (h *handler) listHandler(c echo.Context) error {
	params := listParams{}
	if err := c.Bind(&params); err != nil {
		return err
	}

	var movies []*model.Movie

	err := h.app.DB.
		Model(&movies).
		Limit(params.Limit).
		Offset(params.Offset).
		Order("id DESC").
		Select()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movies)
}

func (h *handler) retrieveHandler(c echo.Context) error {
	id := c.Param("id")

	var movie model.Movie

	err := h.app.DB.Model(&movie).Where("id = ?", id).First()
	if err != nil {
		if err == pg.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "movie not found")
		}
		return err
	}

	return c.JSON(http.StatusOK, movie)
}
