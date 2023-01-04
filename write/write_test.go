package write

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestNewWriter(t *testing.T) {

	Convey("Should handle failure to open postRequestErrors file by exiting program", t, func() {

		restoreOpenFile := stubOpenFile(postRequestErrors)
		defer restoreOpenFile()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

		So(func() { NewWriter() },
			ShouldPanicWith,
			"error opening [errors/postRequestErrors.txt] file")
	})

	Convey("Should handle failure to open unexpectedResponse file by exiting program", t, func() {

		restoreOpenFile := stubOpenFile(unexpectedResponse)
		defer restoreOpenFile()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

		So(func() { NewWriter() },
			ShouldPanicWith,
			"error opening [errors/unexpectedResponse.txt] file")
	})

	Convey("Should handle failure to open missingCompanyName file by exiting program", t, func() {

		restoreOpenFile := stubOpenFile(missingCompanyName)
		defer restoreOpenFile()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

		So(func() { NewWriter() },
			ShouldPanicWith,
			"error opening [errors/missingCompanyName.txt] file")
	})

	Convey("Should handle failure to open missingCompanyData file by exiting program", t, func() {

		restoreOpenFile := stubOpenFile(missingCompanyData)
		defer restoreOpenFile()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

		So(func() { NewWriter() },
			ShouldPanicWith,
			"error opening [errors/missingCompanyData.txt] file")
	})

	Convey("Should handle failure to open alphaKeyErrors file by exiting program", t, func() {

		restoreOpenFile := stubOpenFile(alphaKeyErrors)
		defer restoreOpenFile()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

		So(func() { NewWriter() },
			ShouldPanicWith,
			"error opening [errors/alphaKeyErrors.txt] file")
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
