package router

import (
	"net/http"
	"path"
	"strings"
)

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

// ServeHTTP ...
func ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var p = req.URL.Path
	// var method = req.Method
	// var Body = req.Body
	// search through list, if request matches that list, make res
	if p == "/foo" {
		println("found!")
	}
	http.Error(res, "Not Found", http.StatusNotFound)
}

// IRoutes ...
type IRoutes interface {
	routeExists(string) bool
	createRoutes()
}

// Route ...
type Route struct {
	method string // GET or POST
	path   string
}

// Routes ...
type Routes struct {
	routes []Route
}

// DataResponse ...
type DataResponse struct {
	name     string
	greeting string
}

func (r Routes) getResponse(route Route) DataResponse {
	var res DataResponse
	if route.path == "api/foo" {
		res = DataResponse{
			name:     "Foo",
			greeting: "How are you?",
		}
	}
	return res
}

// routeExists ...
func (r Routes) routeExists(match string) bool {
	var res = false
	for _, v := range r.routes {
		if v.path == match {
			res = true
		}
	}
	return res
}

func createRoutes() Routes {
	r := Routes{}
	r.routes = []Route{
		{path: "api/foo", method: "GET"},
		{path: "api/bar", method: "GET"},
		{path: "api/baz", method: "GET"},
	}
	return r
}
