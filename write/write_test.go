package write

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestNewWriter(t *testing.T) {

	testNewWriterFileOpeningFailure(t, postRequestErrors)
	testNewWriterFileOpeningFailure(t, unexpectedResponse)
	testNewWriterFileOpeningFailure(t, missingCompanyName)
	testNewWriterFileOpeningFailure(t, missingCompanyData)
	testNewWriterFileOpeningFailure(t, alphaKeyErrors)

}

func testNewWriterFileOpeningFailure(t *testing.T, failingFileName string) {

	Convey("Should handle failure to open "+failingFileName+" file by exiting program", t, func() {

		restoreOpenFile := stubOpenFile(failingFileName)
		defer restoreOpenFile()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

		So(func() { NewWriter() },
			ShouldPanicWith,
			"error opening ["+failingFileName+"] file")
	})

}

func stubOpenFile(failingFileName string) func() {
	// Mock out os.OpenFile
	realOpenFile := openFile
	openFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		if name == failingFileName {
			return nil, errors.New("Test generated error")
		}
		return nil, nil
	}
	// Return function to restore os.OpenFile
	return func() { openFile = realOpenFile }
}

func stubLogFatalf() func() {
	// Stub out log.Fatalf
	realFatalf := fatalf
	fatalf = func(format string, v ...interface{}) {
		errorMessage := fmt.Sprintf(format, v...)
		// We replace os.Exit() with panic() because it too exits execution at the right point,
		// but the GoConvey test framework can detect the latter only.
		panic(errorMessage)
	}
	// Return function to restore log.Fatalf
	return func() { fatalf = realFatalf }
}
