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

type Writer interface {
	LogPostError(msg string)
	LogUnexpectedResponse(msg string)
	LogMissingCompanyName(msg string)
	Close()
}

type Write struct {
	pe  *os.File
	ur  *os.File
	mcn *os.File
}

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

func (w *Write) LogPostError(msg string) {
	writeToFile(w.pe, postError, msg)
}

func (w *Write) LogUnexpectedResponse(msg string) {
	writeToFile(w.ur, unexpectedResponse, msg)
}

func (w *Write) LogMissingCompanyName(msg string) {
	writeToFile(w.mcn, missingCompanyName, msg)
}

func writeToFile(connection *os.File, fileName string, msg string) {
	_, err := connection.WriteString(msg + "\n")
	if err != nil {
		log.Printf("error writing [%s] to file fileName: [%s]", msg, fileName)
	}
}
