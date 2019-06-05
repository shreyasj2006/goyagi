package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	app := New()
	assert.NotNil(t, app)
	assert.NotNil(t, app.Config)
}
