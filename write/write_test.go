package write

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

// temporaryFileIndexes provides a mapping from the actual real world file names to the index used for
// corresponding temporary files so that we can refer to the real file names for test intelligibility
// whilst actually working with their temporary substitutes.
var temporaryFileIndexes = map[string]int{
	postRequestErrors:  0,
	unexpectedResponse: 1,
	missingCompanyName: 2,
	missingCompanyData: 3,
	alphaKeyErrors:     4,
}

func TestNewWriter(t *testing.T) {

	testNewWriterFileOpeningFailure(t, postRequestErrors)
	testNewWriterFileOpeningFailure(t, unexpectedResponse)
	testNewWriterFileOpeningFailure(t, missingCompanyName)
	testNewWriterFileOpeningFailure(t, missingCompanyData)
	testNewWriterFileOpeningFailure(t, alphaKeyErrors)

}

func TestClose(t *testing.T) {

	testCloseFileClosingFailure(t, postRequestErrors)
	testCloseFileClosingFailure(t, unexpectedResponse)
	testCloseFileClosingFailure(t, missingCompanyName)
	testCloseFileClosingFailure(t, missingCompanyData)
	testCloseFileClosingFailure(t, alphaKeyErrors)

}

func testCloseFileClosingFailure(t *testing.T, failingFileName string) {

	Convey("Should handle failure to close file "+failingFileName+" by exiting program", t, func() {

		// Keep track of files created during test so that we can remove them.
		temporaryFiles := make(map[int]*os.File)
		defer removeTemporaryFiles(temporaryFiles)

		restoreOpenFile := stubOpenFileWithTempFileCreator(temporaryFiles)
		defer restoreOpenFile()

		restoreClose := stubClose(failingFileName, temporaryFiles)
		defer restoreClose()

		restoreLogFatalf := stubLogFatalf()
		defer restoreLogFatalf()

		writer := NewWriter()

		So(func() { writer.Close() },
			ShouldPanicWith,
			"error closing file: Test generated error closing "+failingFileName)
	})

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

// stubOpenFileWithTempFileCreator stubs out a version of os.OpenFile that creates a valid temporary file
// so that functionality that must have a file present can be tested using the temporary file.
func stubOpenFileWithTempFileCreator(temporaryFiles map[int]*os.File) func() {
	// Mock out os.OpenFile
	realOpenFile := openFile
	openFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		file, err := os.CreateTemp(".", "test")
		if err == nil {
			temporaryFiles[len(temporaryFiles)] = file
		}
		return file, err
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

func stubClose(failingFileName string, temporaryFiles map[int]*os.File) func() {
	// Mock out closeFile
	realCloseFile := closeFile
	closeFile = func(file *os.File) error {
		index := temporaryFileIndexes[failingFileName]
		temporaryFile := temporaryFiles[index]
		if temporaryFile.Name() == file.Name() {
			return errors.New("Test generated error closing " + failingFileName)
		} else {
			return nil
		}
	}
	// Return function to restore closeFile
	return func() { closeFile = realCloseFile }
}

func removeTemporaryFiles(temporaryFiles map[int]*os.File) {
	fmt.Println()
	for index := range temporaryFiles {
		file := temporaryFiles[index]
		fmt.Println("Removing temporary file [" + file.Name() + "]")
		err := os.Remove(file.Name())
		if err != nil {
			fmt.Println("Error encountered removing file " + file.Name() + " [" + err.Error() + "]")
		}
	}
}
