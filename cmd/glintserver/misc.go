package main

import (
	"net/http"
	"strings"
)

func acceptsHtml(r *http.Request) bool {
	var h http.Header = r.Header
	var accept []string = h["Accept"]
	var a string
	for _, a = range accept {
		var b string
		for _, b = range strings.Split(a, ",") {
			if b == "text/html" {
				return true
			}
		}
	}
	return false
}
