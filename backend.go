package main

//StorageBackend is the storage backend interface
type StorageBackend interface {
	Get(URI string) (title string, date string, lang string, content string, err error)
	Put(title, date, lang, content, destDir string) (string, error)
	Flush() error
}
