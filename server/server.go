package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"
)

// fileModeRWX is the umask "-rwx------".
const fileModeRWX = 0700

// Server defines parameters for running a Glint server.
type Server struct {

	// DataDir is the name of a directory for storing the data sets,
	// metadata, and internal files that will be managed by the server.
	DataDir string

	// Host is the TCP host address to listen on.
	Host string

	// Port is the TCP port to listen on.
	Port string

	// TLSCertFile is the name of a file containing a certificate for the
	// server.
	TLSCertFile string

	// TLSKeyFile is the name of a file containing the matching private key
	// for the server.
	TLSKeyFile string

	// Logger optionally specifies a logger for the server.  If nil, the
	// log package's standard logger is used.
	Logger *log.Logger

	// DisableCORS specifies whether to disable handling of CORS (Cross-origin
	// resource sharing) requests.
	DisableCORS bool

	// Debug specifies whether debugging output should be written to the log.
	Debug bool

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string

	DebugAllowInsecureCORS bool

	serverLog

	StaticDir string
	baseURL   string

	StorageModule string
	storage       Storage
}

func (srv *Server) setupCORS(h http.Handler) http.Handler {

	if srv.DebugAllowInsecureCORS {
		srv.log("WARNING: INSECURE CORS CONFIGURATION")
		c := cors.New(cors.Options{
			AllowOriginFunc:  func(origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT"},
			AllowedHeaders:   []string{"authorization", "accept", "content-type"},
			AllowCredentials: true,
			Debug:            true,
		})
		return c.Handler(h)
	}

	return cors.Default().Handler(h)
}

func (srv *Server) handleMain(w http.ResponseWriter, r *http.Request) {

	// Temporary file paths for glintcore website
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		if r.Method == "GET" {
			srv.rootHandler(w, r)
			return
		}
	}

	if r.URL.Path == "/static/style.css" {
		if r.Method == "GET" {
			srv.styleCssHandler(w, r)
			return
		}
	}

	if strings.HasPrefix(r.URL.Path, "/datasets/") {
		if r.Method == "GET" && acceptsHtml(r) {
			srv.datasetsHandler(w, r)
			return
		}
	}

	srv.dataHandler(w, r)

	//srv.unknownRequestHandler(w, r)
}

func (srv *Server) logHandler(h http.Handler, record bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if record {
			srv.logRequest(r, 0)
		}
		h.ServeHTTP(w, r)
	})
}

func (srv *Server) setupHandlers() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/", srv.handleMain)

	// Temporary directory paths for ui-datasets
	mux.Handle("/stripesassets/", http.StripPrefix("/stripesassets/",
		srv.logHandler(
			http.FileServer(
				http.Dir(srv.StaticDir+"/stripesassets")),
			false)))

	// Temporary directory paths for glintcore website
	mux.Handle("/assets/", http.StripPrefix("/assets/",
		srv.logHandler(
			http.FileServer(http.Dir(srv.StaticDir+"/assets")),
			false)))
	mux.Handle("/resources/", http.StripPrefix("/resources/",
		srv.logHandler(
			http.FileServer(http.Dir(srv.StaticDir+"/resources")),
			false)))
	mux.Handle("/about/", http.StripPrefix("/about/",
		srv.logHandler(
			http.FileServer(http.Dir(srv.StaticDir+"/about")),
			false)))

	// Old server handlers
	mux.HandleFunc("/login", srv.handleLogin)
	mux.HandleFunc("/account/password", srv.handleChangePassword)
	mux.HandleFunc("/plot-time-series", handlePlot)

	if !srv.DisableCORS {
		return srv.setupCORS(mux)
	}
	return mux
}

func composeURL(scheme, host, port string) string {

	if (port == "80" && scheme == "http") || (port == "443" && scheme == "https") {
		return scheme + "://" + host + "/"
	}
	return scheme + "://" + net.JoinHostPort(host, port) + "/"
}

func (srv *Server) requireFlags() error {
	if srv.DataDir == "" {
		return fmt.Errorf("Error: data directory not specified")
	}
	return nil
}

func (srv *Server) logExitError(format string, v ...interface{}) {
	srv.log(format, v...)
	srv.log("Exiting with error")
}

func (srv *Server) setupStorage() error {
	if srv.StorageModule == "" {
		srv.storage = new(Postgres)
	} else {
		var err error
		if srv.storage, err = StoragePlugin(
			srv.StorageModule); err != nil {
			return fmt.Errorf(
				"Error loading storage module: %v", err)
		}
	}
	if err := srv.storage.Open(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		srv.PostgresHost, srv.PostgresPort, srv.PostgresUser,
		srv.PostgresPassword, srv.PostgresDBName)); err != nil {
		return fmt.Errorf("Error setting up storage: %v", err)
	}
	if err := srv.storage.Setup(); err != nil {
		return fmt.Errorf("Error setting up storage: %v", err)
	}
	return nil
}

// ListenAndServe listens on the TCP host address srv.Host and port srv.Port,
// and handles requests on incoming connections.  ListenAndServe always
// returns a non-nil error; after Shutdown or Close, the returned error is
// http.ErrServerClosed.
func (srv *Server) ListenAndServe() error {

	if srv.TLSCertFile != "" || srv.TLSKeyFile != "" {
		srv.baseURL = composeURL("https", srv.Host, srv.Port)
	} else {
		srv.baseURL = composeURL("http", srv.Host, srv.Port)
	}

	// Temporary support of old server code
	glintbaseurl = srv.baseURL

	srv.serverLog = serverLog{
		logger: srv.Logger,
		pid:    0,
	}

	if srv.Debug {
		srv.log("%#v\n", srv)
	}

	if err := srv.requireFlags(); err != nil {
		srv.logExitError(err.Error())
		return err
	}

	srv.log("Starting server")

	if srv.Debug {
		srv.log("Setting up storage access")
	}
	if err := srv.setupStorage(); err != nil {
		err = fmt.Errorf("Error setting up storage access: %v", err)
		srv.logExitError(err.Error())
		return err
	}
	defer srv.storage.Close()

	if srv.Debug {
		srv.log("Ensuring data directory \"%s\" exists", srv.DataDir)
	}
	// Create datadir path if it does not exist.
	if err := os.MkdirAll(srv.DataDir, fileModeRWX); err != nil {
		err = fmt.Errorf("Error creating data directory: %v", err)
		srv.logExitError(err.Error())
		return err
	}

	if srv.Debug {
		srv.log("Registering server handlers")
	}
	handler := srv.setupHandlers()

	addr := net.JoinHostPort(srv.Host, srv.Port)

	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	if srv.TLSCertFile != "" || srv.TLSKeyFile != "" {
		srv.log("Server address: %s", srv.baseURL)
		err := server.ListenAndServeTLS(srv.TLSCertFile, srv.TLSKeyFile)
		if err != nil {
			err = fmt.Errorf("Error starting server: %v", err)
			srv.logExitError(err.Error())
			return err
		}
	} else {
		srv.log("Server address: %s", srv.baseURL)
		err := server.ListenAndServe()
		if err != nil {
			err = fmt.Errorf("Error starting server: %v", err)
			srv.logExitError(err.Error())
			return err
		}
	}

	srv.log("Shutting down server")
	return http.ErrServerClosed
}
