package write

import (
	"log"
	"os"
)

const (
	postError          = "errors/error-posting-request.txt"
	unexpectedResponse = "errors/unexpected-put-response.txt"
	missingCompanyName = "errors/missing-company-name.txt"
)

// Writer provides an interface by which to write error messages to log files
type Writer interface {
	LogPostError(msg string)
	LogUnexpectedResponse(msg string)
	LogMissingCompanyName(msg string)
	Close()
}

// Write provides a concrete implementation of the Writer interface
type Write struct {
	pe  *os.File
	ur  *os.File
	mcn *os.File
}

// NewWriter returns a concrete implementation of the Writer interface
func NewWriter() Writer {

	postErrorFile, err := os.OpenFile(postError, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening [%s] file", postError)
	}

	unexpectedResponseFile, err := os.OpenFile(unexpectedResponse, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening [%s] file", unexpectedResponse)
	}

	missingCompanyNameFile, err := os.OpenFile(missingCompanyName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening [%s] file", missingCompanyName)
	}

	return &Write{
		pe:  postErrorFile,
		ur:  unexpectedResponseFile,
		mcn: missingCompanyNameFile,
	}
}

// Close closes a Writer
func (w *Write) Close() {

	if err := w.pe.Close(); err != nil {
		log.Fatalf("error closing file: %s", err)
	}
	if err := w.ur.Close(); err != nil {
		log.Fatalf("error closing file: %s", err)
	}
	if err := w.mcn.Close(); err != nil {
		log.Fatalf("error closing file: %s", err)
	}
}

// LogPostError logs an error to the 'error-posting-request' file
func (w *Write) LogPostError(msg string) {
	writeToFile(w.pe, postError, msg)
}

// LogUnexpectedResponse logs an error to the 'unexpected-put-response' file
func (w *Write) LogUnexpectedResponse(msg string) {
	writeToFile(w.ur, unexpectedResponse, msg)
}

// LogMissingCompanyName logs an error to the 'missing-company-name' file
func (w *Write) LogMissingCompanyName(msg string) {
	writeToFile(w.mcn, missingCompanyName, msg)
}

func writeToFile(connection *os.File, fileName string, msg string) {
	_, err := connection.WriteString(msg + "\n")
	if err != nil {
		log.Printf("error writing [%s] to file: [%s]", msg, fileName)
	}
}
