package write

import (
	"log"
	"os"
)

const (
	postRequestErrors  = "errors/postRequestErrors.txt"
	unexpectedResponse = "errors/unexpectedResponse.txt"
	missingCompanyName = "errors/missingCompanyName.txt"
	missingCompanyData = "errors/missingCompanyData.txt"
	alphaKeyErrors     = "errors/alphaKeyErrors.txt"
	errorOpeningFile   = "error opening [%s] file"
	errorClosingFile   = "error closing file: %s"
)

// Writer provides an interface by which to write error messages to log files
type Writer interface {
	LogPostError(msg string)
	LogUnexpectedResponse(msg string)
	LogMissingCompanyName(msg string)
	LogMissingCompanyData(msg string)
	LogAlphaKeyErrors(msg string)
	Close()
}

// Write provides a concrete implementation of the Writer interface
type Write struct {
	pe  *os.File
	ur  *os.File
	mcn *os.File
	mcd *os.File
	ake *os.File
}

// Function variables to facilitate testing.
var (
	openFile = os.OpenFile
	fatalf   = log.Fatalf
)

// NewWriter returns a concrete implementation of the Writer interface
func NewWriter() Writer {

	postErrorFile, err := openFile(postRequestErrors, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fatalf(errorOpeningFile, postRequestErrors)
	}

	unexpectedResponseFile, err := openFile(unexpectedResponse, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fatalf(errorOpeningFile, unexpectedResponse)
	}

	missingCompanyNameFile, err := openFile(missingCompanyName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fatalf(errorOpeningFile, missingCompanyName)
	}

	missingCompanyDataFile, err := openFile(missingCompanyData, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fatalf(errorOpeningFile, missingCompanyData)
	}

	alphaKeyErrorsFile, err := openFile(alphaKeyErrors, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fatalf(errorOpeningFile, alphaKeyErrors)
	}

	return &Write{
		pe:  postErrorFile,
		ur:  unexpectedResponseFile,
		mcn: missingCompanyNameFile,
		mcd: missingCompanyDataFile,
		ake: alphaKeyErrorsFile,
	}
}

// Close closes a Writer
func (w *Write) Close() {

	if err := w.pe.Close(); err != nil {
		log.Fatalf(errorClosingFile, err)
	}
	if err := w.ur.Close(); err != nil {
		log.Fatalf(errorClosingFile, err)
	}
	if err := w.mcn.Close(); err != nil {
		log.Fatalf(errorClosingFile, err)
	}
	if err := w.mcd.Close(); err != nil {
		log.Fatalf(errorClosingFile, err)
	}
	if err := w.ake.Close(); err != nil {
		log.Fatalf(errorClosingFile, err)
	}
}

// LogPostError logs an error to the 'error-posting-request' file
func (w *Write) LogPostError(msg string) {
	writeToFile(w.pe, postRequestErrors, msg)
}

// LogUnexpectedResponse logs an error to the 'unexpected-put-response' file
func (w *Write) LogUnexpectedResponse(msg string) {
	writeToFile(w.ur, unexpectedResponse, msg)
}

// LogMissingCompanyName logs an error to the 'missing-company-name' file
func (w *Write) LogMissingCompanyName(msg string) {
	writeToFile(w.mcn, missingCompanyName, msg)
}

// LogMissingCompanyData logs an error to the 'missingCompanyData' file
func (w *Write) LogMissingCompanyData(msg string) {
	log.Println(msg) // This really is very bad data, log to console too.
	writeToFile(w.mcd, missingCompanyData, msg)
}

func (w *Write) LogAlphaKeyErrors(msg string) {
	writeToFile(w.ake, alphaKeyErrors, msg)
}

func writeToFile(connection *os.File, fileName string, msg string) {
	_, err := connection.WriteString(msg + "\n")
	if err != nil {
		log.Printf("error writing [%s] to file: [%s]", msg, fileName)
	}
}
