package router

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestRouter(t *testing.T) {
	var r = createRouter()
	assert.Equal(t, Router.routeExists(r, "api/foo"), true)
}
