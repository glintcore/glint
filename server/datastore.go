package server

import (
	"database/sql"
	"fmt"
	"plugin"
)

type Storage interface {
	Open(dataSourceName string) error
	//Connect(host, port, user, password, dbname string) error

	Close() error

	Setup() error

	//AddMetadata()
	AddMetadata(personId int64, path string, attribute string,
		metadata string) error

	//LookupMetadata()
	LookupMetadata(personId int64, path string, attribute string) (string,
		error)

	//LookupData()
	LookupData(person_id int64, path string) (string, error)

	//LookupDataList()
	LookupDataList(person_id int64) (string, error)

	//AddFile()
	AddFile(person_id int64, path string, data string) (int64, error)

	//DeleteFile()
	DeleteFile(personId int64, path string) error

	LookupPassword(username string) (string, error)
	Authenticate(username string, password string) (bool, error)
	CreateTablePerson(tx *sql.Tx) error
	ChangePassword(username string, password string) error
	LookupPersonId(username string) (int64, error)
	AddPerson(username string, fullname string, email string,
		password string) error
	LookupFileId(personId int64, path string) (int64, error)
	AddAttributes(file_id int64, attrs []string) error
	CreateTableAttribute(tx *sql.Tx) error
	CreateTableFile(tx *sql.Tx) error
	CreateSchema() error
}

func StoragePlugin(pluginFile string) (Storage, error) {

	p, err := plugin.Open(pluginFile)
	if err != nil {
		return nil, err
	}

	sym, err := p.Lookup("StorageModule")
	if err != nil {
		return nil, err
	}

	var module Storage
	module, ok := sym.(Storage)
	if !ok {
		return nil, fmt.Errorf(
			"Module does not match Storage interface: "+
				"StorageModule in module %v", pluginFile)
	}
	return module, nil
}
