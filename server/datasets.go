package server

import (
	"net/http"
	"os"
)

func (srv *Server) datasetsHandler(w http.ResponseWriter, r *http.Request) {

	var statusCode int

	index := srv.StaticDir + "/stripesassets/index.html"
	if _, err := os.Stat(index); err == nil {
		statusCode = http.StatusOK
		http.ServeFile(w, r, index)
	} else {
		statusCode = http.StatusBadRequest
		w.WriteHeader(statusCode)
	}
}
