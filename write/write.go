package write

import (
	"log"
	"os"
)

const (
	filename1 = "company-errors/error-posting-request.txt"
	filename2 = "company-errors/unexpected-put-response.txt"
	filename3 = "company-errors/missing-company-name.txt"
)

type Write interface {
	WriteToFile1(sentence string)
	WriteToFile2(sentence string)
	WriteToFile3(sentence string)
	Close()
}

type Writer struct {
	f1 *os.File
	f2 *os.File
	f3 *os.File
}

func New() *Writer {

	connection1, err := os.OpenFile(filename1, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening [%s] file", filename1)
	}

	connection2, err := os.OpenFile(filename2, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening [%s] file", filename2)
	}

	connection3, err := os.OpenFile(filename3, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening [%s] file", filename3)
	}

	return &Writer{
		f1: connection1,
		f2: connection2,
		f3: connection3,
	}
}

func (w *Writer) Close() {

	if err := w.f1.Close(); err != nil {
		log.Fatalf("error closing file: %s", err)
	}
	if err := w.f2.Close(); err != nil {
		log.Fatalf("error closing file: %s", err)
	}
	if err := w.f3.Close(); err != nil {
		log.Fatalf("error closing file: %s", err)
	}
}

func (w *Writer) WriteToFile1(sentence string) {
	writeToFile(w.f1, filename1, sentence)
}

func (w *Writer) WriteToFile2(sentence string) {
	writeToFile(w.f2, filename2, sentence)
}

func (w *Writer) WriteToFile3(sentence string) {
	writeToFile(w.f3, filename3, sentence)
}

func writeToFile(connection *os.File, location string, sentence string) {
	_, err := connection.WriteString(sentence + "\n")
	if err != nil {
		log.Printf("error writing [%s] to file location: [%s]", sentence, location)
	}
}
