package server

type StorageFiles struct {
	dataDir string
}

func (fs *StorageFiles) Open(dataSourceName string) error {
	fs.dataDir = dataSourceName
	return nil
}

func (fs *StorageFiles) Close() error {
	return nil
}
