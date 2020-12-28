package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/glintdb/glintweb/api"
)

func (srv *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		var m = "HTTP method " + r.Method +
			" is not supported by this URL"
		http.Error(w, m, http.StatusMethodNotAllowed)
		log.Println(m)
		return
	}
	// Authenticate user.
	var ok bool
	_, ok = srv.handleBasicAuth(w, r)
	if !ok {
		return
	}
	// Read the json request.
	var body []byte
	var err error
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
	var p api.LoginRequest
	err = json.Unmarshal(body, &p)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
	// Write the json response.
	var resp api.LoginResponse
	resp.SessionId = "0"
	var respbody []byte
	respbody, err = json.Marshal(resp)
	if err != nil {
		// TODO Handle error.
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respbody)
}
