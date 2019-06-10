package application

import (
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/shreyasj2006/goyagi/pkg/config"
	"github.com/shreyasj2006/goyagi/pkg/database"
)

// App contains necessary references that will be persisted throughout the
// application's lifecycle.
type App struct {
	Config config.Config
	DB     *pg.DB
}

// New creates a new instance of App
func New() (App, error) {
	cfg := config.New()

	db, err := database.New(cfg)
	if err != nil {
		return App{}, errors.Wrap(err, "application")
	}

	return App{
		Config: cfg,
		DB:     db,
	}, nil
}
