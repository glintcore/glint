package server

import (
	"fmt"
	"net/http"
	"os"
)

func (srv *Server) rootHandler(w http.ResponseWriter, r *http.Request) {

	statusCode := http.StatusOK

	index := srv.StaticDir + r.URL.Path
	if _, err := os.Stat(index); err == nil {
		http.ServeFile(w, r, index)
	} else {
		setContentTypeTextHtml(w)
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, "The Glint server is running at: "+
			"<pre>%s</pre>\n", srv.baseURL)
	}
}
