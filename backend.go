package main

//StorageBackend is the storage backend interface
type StorageBackend interface {
	Get(URI string) (title string, date string, lang string, content []byte, err error)
	Put(title, date, lang string, content []byte, destDir string) (string, error)
	Flush() error
}
