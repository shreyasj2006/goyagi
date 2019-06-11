package sentry

import (
	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"github.com/shreyasj2006/goyagi/pkg/config"
)

// RavenClient provides an interface for the Sentry client.
type ravenClient interface {
	Capture(packet *raven.Packet, captureTags map[string]string) (string, chan error)
}

// Sentry contains a client used to send exceptions to Sentry.io.
// We are using an interface in our struct instead of the actual
// client type in order to allow us to mock the clients behavior
// in our tests.
type Sentry struct {
	Client ravenClient
}

// New returns an instance of Sentry.
func New(cfg config.Config) (Sentry, error) {
	defaultTags := map[string]string{
		"environment": cfg.Environment,
	}

	client, err := raven.NewWithTags(cfg.SentryDSN, defaultTags)
	if err != nil {
		return Sentry{}, errors.Wrap(err, "sentry")
	}

	return Sentry{client}, nil
}
