package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/glintdb/glintweb/server"
	"github.com/nassibnassar/goconfig/ini"
	"github.com/urfave/cli"
)

const errorPrefix = "glintserver: "

func serverErr(err error) error {
	return fmt.Errorf("%s%v", errorPrefix, err)
}

func coalesce(s1, s2, s3 string) string {
	if s1 != "" {
		return s1
	}
	if s2 != "" {
		return s2
	}
	return s3
}

// readConfig reads a configuration file specified by the command line flag
// --config-file, or if not available then by the environment variable
// GLINTSERVER_CONFIG_FILE.
func readConfig(c *cli.Context) (*ini.Config, error) {

	file := coalesce(c.GlobalString("config-file"),
		os.Getenv("GLINTSERVER_CONFIG_FILE"),
		"")
	if file == "" {
		return ini.NewConfig(), nil
	}
	return ini.NewConfigFile(file)
}

// logToFile opens the specified file for append and sets the log package's
// standard logger to write to that file.  The file pointer is returned and
// should be eventually closed.
func logToFile(logfile string) (*os.File, error) {

	if logfile == "" {
		return nil, nil
	}

	// Open the file.
	file, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}
	// Set the logger to write to the file.
	log.SetOutput(file)
	return file, nil
}

// closeLogFile closes the file f if f != nil, unlike f.Close() which returns
// an error if f == nil.  If an error occurs, panic() is called.
func closeLog(f *os.File) {
	if f != nil {
		if err := f.Close(); err != nil {
			panic("Error closing log file: " + err.Error())
		}
	}
}

func warnConfigChange(config *ini.Config, key, newkey string) {
	keysp := strings.Split(key, ".")
	if config.Get(keysp[0], keysp[1]) != "" {
		fmt.Fprintf(os.Stderr, errorPrefix+"Warning:\n\n    Configuration key "+
			"\"\033[1m\033[31m%s\033[0m\033[0m\" is no longer supported", key)
		if newkey != "" {
			fmt.Fprintf(os.Stderr, ";\n                  use "+
				"\"\033[1m\033[32m%s\033[0m\033[0m\" instead.", newkey)
		}
		fmt.Fprintf(os.Stderr, "\n\n")
	}
}

func cliNewServer(c *cli.Context) error {

	config, err := readConfig(c)
	if err != nil {
		return serverErr(fmt.Errorf("Error reading configuration file: %v", err))
	}

	warnConfigChange(config, "debug.port", "http.port")
	warnConfigChange(config, "http.sslcert", "http.tlscert")
	warnConfigChange(config, "http.sslkey", "http.tlskey")

	logf, err := logToFile(coalesce("", config.Get("log", "file"), ""))
	if err != nil {
		return serverErr(fmt.Errorf("Error writing to log file: %v", err))
	}
	defer closeLog(logf)

	srv := &server.Server{
		DataDir: coalesce("", config.Get("core", "datadir"), ""),
		StaticDir: coalesce("", config.Get("core", "staticdir"),
			"/var/glint/html"),
		Host: coalesce("", config.Get("http", "host"),
			"localhost"),
		Port: coalesce(c.String("port"), config.Get("http",
			"port"), "8080"),
		TLSCertFile: coalesce("", config.Get("http", "tlscert"), ""),
		TLSKeyFile:  coalesce("", config.Get("http", "tlskey"), ""),
		Debug:       c.GlobalBool("debug"),
		PostgresHost: coalesce("", config.Get("database",
			"host"), ""),
		PostgresPort: coalesce("", config.Get("database",
			"port"), ""),
		PostgresUser: coalesce("", config.Get("database",
			"user"), ""),
		PostgresPassword: coalesce("", config.Get("database",
			"password"), ""),
		PostgresDBName: coalesce("", config.Get("database",
			"dbname"), ""),
		DebugAllowInsecureCORS: c.Bool("debug-allow-insecure-cors"),
		StorageModule: coalesce("", config.Get("storage",
			"module"), ""),
	}

	err = srv.ListenAndServe()
	if err != http.ErrServerClosed {
		return serverErr(fmt.Errorf("Server exited with error"))
	}
	return nil
}
