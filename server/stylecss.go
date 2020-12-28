package server

import (
	"fmt"
	"net/http"
)

func (srv *Server) styleCssHandler(w http.ResponseWriter, r *http.Request) {

	statusCode := http.StatusOK

	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s", styleCss())
}
