package server

import (
	"fmt"
	"net/http"
)

func (srv *Server) unknownRequestHandler(w http.ResponseWriter,
	r *http.Request) {

	statusCode := http.StatusNotFound

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "<html><head><title>404 Not Found</title></head><body>"+
		"<h1>404 Not Found</h1></body></html>")
}
