package main

// ReadWriter is the interface that groups Read and Write for pastes
type ReadWriter interface {
	ReadPaste(URI string) (title string, date string, lang string, content string, err error)
	WritePaste(title, date, lang, content, destDir string) (string, error)
}
