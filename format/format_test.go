package format

import (
	"testing"
)

func TestSplitCompanyNameEndings(t *testing.T) {

	f := &Format{}

	coName := "TEST LIMITED"

	nameStart, nameEnd := f.SplitCompanyNameEndings(coName)

	if nameStart != "TEST" {
		t.Errorf("Expected %v got %v", "TEST", nameStart)
	}

	if nameEnd != " LIMITED" {
		t.Errorf("Expected %v got %v", " LIMITED", nameEnd)
	}

}
