package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/glintdb/glintweb/api"
)

var glintbaseurl string

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

// joinURLPath concatenates two parts (s1, s2) of a URL path to create a single
// path.  The separator string "/" is placed between s1 and s2 in the resulting
// string.  If s1 ends with "/" or s2 begins with "/", the extra separators are
// removed.
func joinURLPath(s1, s2 string) string {
	s1last := s1[len(s1)-1:]
	s2first := s2[0:1]
	var b strings.Builder
	if s1last == "/" {
		b.WriteString(s1[0 : len(s1)-1])
	} else {
		b.WriteString(s1)
	}
	b.WriteString("/")
	if s2first == "/" {
		b.WriteString(s2[1:])
	} else {
		b.WriteString(s2)
	}
	return b.String()
}

func setContentTypeTextHtml(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func setContentTypeTextPlain(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

// Authenticate user provided via HTTP basic authentication, returning the
// username if possible.
func (srv *Server) handleBasicAuth(w http.ResponseWriter, r *http.Request) (
	string, bool) {
	var user, password string
	var ok bool
	user, password, ok = r.BasicAuth()
	if !ok {
		var m = "Unauthorized: Invalid HTTP Basic Authentication"
		log.Println(m)
		//w.Header().Set("WWW-Authenticate", "Basic")
		http.Error(w, m, http.StatusForbidden)
		return user, false
	}
	var match bool
	var err error
	match, err = srv.storage.Authenticate(user, password)
	if err != nil {
		var m = "Unauthorized (user '" + user + "')"
		log.Println(m + ": " + err.Error())
		//w.Header().Set("WWW-Authenticate", "Basic")
		http.Error(w, m, http.StatusForbidden)
		return user, false
	}
	if !match {
		var m = "Unauthorized (user '" + user + "'): " +
			"Unable to authenticate username/password"
		log.Println(m)
		//w.Header().Set("WWW-Authenticate", "Basic")
		http.Error(w, m, http.StatusForbidden)
		return user, false
	}
	return user, true
}

func handleError(w http.ResponseWriter, err error, statusCode int) {
	var m string = err.Error()
	log.Println(m)
	http.Error(w, m, statusCode)
}

func (srv *Server) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		var m = "HTTP method " + r.Method +
			" is not supported by this URL"
		http.Error(w, m, http.StatusMethodNotAllowed)
		log.Println(m)
		return
	}
	// Authenticate user.
	var user string
	var ok bool
	user, ok = srv.handleBasicAuth(w, r)
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
	var p api.AccountPasswordRequest
	err = json.Unmarshal(body, &p)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
	// Set the new password.
	err = srv.storage.ChangePassword(user, p.Password)
	if err != nil {
		var m = "Unable to update password: " + err.Error()
		http.Error(w, m, http.StatusBadRequest)
		log.Println(m)
		return
	}
	// Respond with success.
	//w.Header().Set("Content-Type", "text/plain")
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//fmt.Fprintf(w, "Password successfully changed\n")
}

func parsePathElements(r *http.Request) []string {
	var split []string = strings.Split(r.URL.Path, "/")
	var parsed []string
	var s string
	for _, s = range split {
		if strings.TrimSpace(s) != "" {
			parsed = append(parsed, s)
		}
	}
	return parsed
}

type rootData struct {
	Root string
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var t *template.Template
		t, _ = template.New("index").Parse(index())
		t.Execute(w, rootData{Root: "ping"})
	}
	// TODO Access error for non-GET requests.
}

func parsePathBasic(r *http.Request) (string, string, error) {
	var p []string = parsePathElements(r)
	if len(p) < 1 {
		return "", "", fmt.Errorf("Path not supported: %s", r.URL.Path)
	}
	if len(p) > 2 {
		return "", "", fmt.Errorf("Path not supported: %s", r.URL.Path)
	}
	var pathUser, pathDataName string
	pathUser = p[0]
	if len(p) == 2 {
		pathDataName = p[1]
	}
	return pathUser, pathDataName, nil
}

func header() string {
	return `
<html>

<head>
    <meta http-equiv="Content-type" content="text/html; charset=utf-8">
        <link rel="stylesheet" href="/static/style.css" type="text/css"
	          media="screen" title="no title">
</head>

<body>
`
}

func footer() string {
	return `
</body>

</html>
`
}

func (srv *Server) fprintData(w http.ResponseWriter, html bool, thp map[string][]string,
	personId int64, user string, path string, data string) {

	if html {
		fmt.Fprintf(w, "%s", header())
		fmt.Fprintf(w, "<h1><a href=\"/%s\">%s</a> / %s</h1>\n",
			user, user, path)
		fmt.Fprintf(w, "<table>\n")
	}
	var rows []string = strings.Split(data, "\n")
	var r int
	for r = range rows {
		if strings.TrimSpace(rows[r]) == "" {
			continue
		}
		var cells []string = strings.Split(rows[r], ",")
		if html {
			fmt.Fprintf(w, "<tr>")
		}
		var c int
		for c = range cells {
			if html {
				if r == 0 {
					var md string
					if thp["md"] != nil {
						md, _ = srv.storage.LookupMetadata(
							personId, path,
							cells[c])
					}
					if md == "" {
						md = "&nbsp;"
					}
					fmt.Fprintf(w,
						"<th><div>%s</div><div>%s</div></th>",
						cells[c], md)
				} else {
					if path == "" {
						fmt.Fprintf(w, "<td>"+
							"<a href=\"/%s/%s\">"+
							"%s"+
							"</a>"+
							"</td>",
							user, cells[c],
							cells[c])
					} else {
						fmt.Fprintf(w, "<td>%s</td>",
							cells[c])
					}
				}
			} else {
				if c > 0 {
					fmt.Fprintf(w, ",")
				}
				fmt.Fprintf(w, "%s", cells[c])
				if r == 0 && thp["md"] != nil {
					var md string
					md, _ = srv.storage.LookupMetadata(personId,
						path, cells[c])
					fmt.Fprintf(w, "%s", md)
				}
			}
		}
		if html {
			fmt.Fprintf(w, "</tr>")
		}
		fmt.Fprintf(w, "\n")
	}
	if html {
		fmt.Fprintf(w, "</table>\n")
		fmt.Fprintf(w, "%s", footer())
	}
}

func writeStatusCode(w http.ResponseWriter, code int) {
	setContentTypeTextHtml(w)
	w.WriteHeader(code)
	fmt.Fprintf(w, "%d %s\n", code, http.StatusText(code))
}

func (srv *Server) handleDataGet(w http.ResponseWriter, r *http.Request) {

	var pathUser, pathDataName string
	var err error
	pathUser, pathDataName, err = parsePathBasic(r)
	if err != nil {
		writeStatusCode(w, http.StatusBadRequest)
		return
	}

	var thp map[string][]string = thumpParseBasic(r)

	var personId int64
	personId, err = srv.storage.LookupPersonId(pathUser)
	if err != nil {
		writeStatusCode(w, http.StatusNotFound)
		return
	}

	var data string
	if pathDataName == "" {
		data, err = srv.storage.LookupDataList(personId)
		if err != nil {
			writeStatusCode(w, http.StatusNotFound)
			return
		}
	} else {
		data, err = srv.storage.LookupData(personId, pathDataName)
		if err != nil {
			writeStatusCode(w, http.StatusNotFound)
			return
		}
	}

	if thp["show"] != nil {
		data = thumpShowBasic(data, thp["show"])
	}
	if acceptsHtml(r) {
		setContentTypeTextHtml(w)
		w.WriteHeader(http.StatusOK)
	} else {
		var contentType string
		if thp["as"] != nil && thp["as"][0] == "tsv" {
			contentType = "text/tab-separated-values"
		} else {
			contentType = "text/csv"
		}
		w.Header().Set("Content-Type", contentType+"; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if thp["as"] != nil && thp["as"][0] == "tsv" {
			data = strings.Replace(data, ",", "\t", -1)
		}
	}
	srv.fprintData(w, acceptsHtml(r), thp, personId, pathUser, pathDataName,
		data)
}

func handleCss(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", styleCss())
}

func (srv *Server) handleMetadataPut(w http.ResponseWriter, r *http.Request) {
	// Authenticate user.
	var user string
	var ok bool
	user, ok = srv.handleBasicAuth(w, r)
	if !ok {
		return
	}

	var pathUser, pathDataName string
	var err error
	pathUser, pathDataName, err = parsePathBasic(r)
	if err != nil {
		// TODO Handle error.
	}
	if pathUser != user {
		// TODO Handle error.
	}

	var sp []string = strings.Split(pathDataName, ".")
	var path = sp[0]
	var attribute = sp[1]

	// Read the json request.
	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
	var req api.MetadataRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	var personId int64
	personId, err = srv.storage.LookupPersonId(user)
	if err != nil {
		log.Print(err)
		// TODO Handle error.
		return
	}

	err = srv.storage.AddMetadata(personId, path, attribute, req.Metadata)
	if err != nil {
		log.Print(err)
		// TODO Handle error.
		return
	}

	var resp api.PostResponse
	resp.Url = glintbaseurl + "/" + pathUser + "/" + path
	var respbody []byte
	respbody, err = json.Marshal(resp)
	if err != nil {
		// TODO Handle error.
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respbody)

}

func (srv *Server) handleDataPut(w http.ResponseWriter, r *http.Request) {
	// Authenticate user.
	var user string
	var ok bool
	user, ok = srv.handleBasicAuth(w, r)
	if !ok {
		return
	}

	var pathUser string
	var pathDataName string
	var err error
	pathUser, pathDataName, err = parsePathBasic(r)
	if err != nil {
		// TODO Handle error.
	}
	if pathUser != user {
		// TODO Handle error.
	}

	// Read the json request.
	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}
	var req api.PostRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	var personId int64
	personId, err = srv.storage.LookupPersonId(user)
	if err != nil {
		log.Print(err)
		// TODO Handle error.
		return
	}

	var data string = strings.Replace(req.Data, "\\n", "\n", -1)

	var id int64
	id, err = srv.storage.AddFile(personId, pathDataName, data)
	if err != nil {
		log.Print(err)
		// TODO Handle error.
		return
	}

	var sp []string = strings.Split(data, "\n")
	var attrs = strings.Split(sp[0], ",")

	err = srv.storage.AddAttributes(id, attrs)
	if err != nil {
		log.Print(err)
		// TODO Handle error.
		return
	}

	var resp api.PostResponse
	resp.Url = joinURLPath(glintbaseurl, pathUser+"/"+pathDataName)
	var respbody []byte
	respbody, err = json.Marshal(resp)
	if err != nil {
		// TODO Handle error.
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// TODO Return http.StatusOK if updated rather than created.
	w.Write(respbody)

}

func (srv *Server) handleDataDelete(w http.ResponseWriter, r *http.Request) {
	// Authenticate user.
	var user string
	var ok bool
	user, ok = srv.handleBasicAuth(w, r)
	if !ok {
		return
	}

	var pathUser string
	var pathDataName string
	var err error
	pathUser, pathDataName, err = parsePathBasic(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if pathUser != user {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var personId int64
	personId, err = srv.storage.LookupPersonId(user)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = srv.storage.DeleteFile(personId, pathDataName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (srv *Server) dataHandler(w http.ResponseWriter, r *http.Request) {
	srv.logRequest(r, 0)
	switch r.Method {
	case http.MethodGet:
		srv.handleDataGet(w, r)
	case http.MethodPut:
		if strings.ContainsRune(r.URL.Path, '.') {
			srv.handleMetadataPut(w, r)
		} else {
			srv.handleDataPut(w, r)
		}
	case http.MethodDelete:
		if strings.ContainsRune(r.URL.Path, '.') {
			w.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			srv.handleDataDelete(w, r)
		}
	default:
		// TODO Access error for non-GET/PUT requests.
	}
}
