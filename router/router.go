package router

import (
	"fmt"
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

// Routes ...
type Routes struct {
	paths []string
}

func (r Routes) routeExists(match string) bool {
	var res = false
	for _, v := range r.paths {
		if v == match {
			res = true
		}
	}
	return res
}

func createRoutes() Routes {
	r := Routes{}
	r.paths = []string{
		"api/foo",
		"api/bar",
		"api/baz",
	}
	return r
}

func initMe(i IRoutes) {
	fmt.Println(i.routeExists("foo"))
}
