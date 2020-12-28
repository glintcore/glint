package main

import (
	"fmt"
	"log"
	"os"

	"github.com/glintdb/glintweb/server"
	"github.com/urfave/cli"
)

// These functions are still used by the adduser and passwd server commands.

func openLogFile() (*os.File, error) {
	// Set log file path.
	var path string = glintconfig.Get("log", "file")
	if path == "" {
		return nil, nil
	}
	// Open log file.
	var file *os.File
	var err error
	file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0640)
	if err != nil {
		return nil, err
	}
	// Redirect messages to log file.
	log.SetOutput(file)
	return file, nil
}

func closeLogFile(f *os.File) error {
	if f != nil {
		var err = f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func setup(c *cli.Context, logging bool) (*os.File, server.Storage, error) {
	const configPathEnv = "GLINTSERVER_CONFIG_FILE"
	// The location of the server configuration file is specified by the
	// GLINTSERVER_CONFIG_FILE environment variable, unless overridden by
	// the --config-file command line flag.
	var configPath string
	var configFileFlag = c.GlobalString("config-file")
	if configFileFlag == "" {
		configPath = os.Getenv(configPathEnv)
	} else {
		configPath = configFileFlag
	}
	// At the moment the configuration file is the only way to configure
	// the server; so exit with an error if the environment variable has
	// not been set or if the file cannot be read.
	if configPath == "" {
		return nil, nil, fmt.Errorf(
			"Environment variable '%s' has not been set",
			configPathEnv)
	}
	var err error
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return nil, nil, fmt.Errorf(
			"Unable to read server configuration file: %v", err)
	}
	// Read configuration file.
	glintconfig = oldReadConfig(configPath)
	// Set up logging.
	var logfile *os.File = nil
	if logging {
		logfile, err = openLogFile()
		if err != nil {
			return nil, nil, err
		}
	}
	// Set up Postgres.
	/*
		if err := server.PostgresSetup(glintconfig.Get("database", "host"),
			glintconfig.Get("database", "port"), glintconfig.Get("database", "user"),
			glintconfig.Get("database", "password"),
			glintconfig.Get("database", "dbname")); err != nil {
			return nil, fmt.Errorf("Database error: %v", err)
		}
	*/
	///////////////////////////////////////////////////////////////////////
	var storage server.Storage
	storageModule := glintconfig.Get("storage", "module")
	if storageModule == "" {
		storage = new(server.Postgres)
	} else {
		var err error
		if storage, err = server.StoragePlugin(
			storageModule); err != nil {
			return nil, nil, fmt.Errorf(
				"Error loading storage module: %v", err)
		}
	}
	if err := storage.Open(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		glintconfig.Get("database", "host"),
		glintconfig.Get("database", "port"),
		glintconfig.Get("database", "user"),
		glintconfig.Get("database", "password"),
		glintconfig.Get("database", "dbname"))); err != nil {
		return nil, nil, fmt.Errorf("Error setting up storage: %v", err)
	}
	if err := storage.Setup(); err != nil {
		return nil, nil, fmt.Errorf("Error setting up storage: %v", err)
	}
	///////////////////////////////////////////////////////////////////////
	return logfile, storage, nil
}

func cleanup(logfile *os.File, storage server.Storage) {
	storage.Close()
	if logfile != nil {
		closeLogFile(logfile)
	}
}
