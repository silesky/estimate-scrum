package router

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestRouter(t *testing.T) {
	var r = createRoutes()
	assert.Equal(t, Routes.routeExists(r, "/foo"), true)
}
