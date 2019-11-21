package format

import (
	"regexp"
	"strings"
)

var nonWordEndRegex = regexp.MustCompile(`[^A-Za-z0-9_]+$`)

var companyNameEndings = [...]string{
	"AEIE",
	"ANGHYFYNGEDIG",
	"C.B.C",
	"C.C.C",
	"C.I.C",
	"CBC",
	"CBCN",
	"CBP",
	"CCC",
	"CCG CYF",
	"CCG CYFYNGEDIG",
	"CIC",
	"COMMUNITY INTEREST COMPANY",
	"COMMUNITY INTEREST P.L.C",
	"COMMUNITY INTEREST PLC",
	"COMMUNITY INTEREST PUBLIC LIMITED COMPANY",
	"CWMNI BUDDIANT C.C.C",
	"CWMNI BUDDIANT CCC",
	"CWMNI BUDDIANT CYMUNEDOL C.C.C",
	"CWMNI BUDDIANT CYMUNEDOL CCC",
	"CWMNI BUDDIANT CYMUNEDOL CYHOEDDUS CYFYNGEDIG",
	"CWMNI BUDDIANT CYMUNEDOL",
	"CWMNI BUDDSODDIA CHYFALAF NEWIDIOL",
	"CWMNI BUDDSODDIANT PENAGORED",
	"CWMNI CELL GWARCHODEDIG",
	"CWMNI CYFYNGEDIG CYHOEDDUS",
	"CYF",
	"CYFYNGEDIG",
	"EEIG",
	"EESV",
	"EOFG",
	"EOOS",
	"EUROPEAN ECONOMIC INTEREST GROUPING",
	"GEIE",
	"GELE",
	"ICVC",
	"INVESTMENT COMPANY WITH VARIABLE CAPITAL",
	"L.P",
	"L.T.D",
	"LIMITED - THE",
	"LIMITED LIABILITY PARTNERSHIP",
	"LIMITED PARTNERSHIP",
	"LIMITED THE",
	"LIMITED",
	"LIMITED-THE",
	"LIMITED...THE",
	"LIMITED..THE",
	"LIMITED.THE",
	"LLP",
	"LP",
	"LTD",
	"LTD...THE",
	"LTD..THE",
	"LTD.THE",
	"OEIC",
	"OPEN-ENDED INVESTMENT COMPANY",
	"P.L.C",
	"PAC",
	"PARTNERIAETH ATEBOLRWYDD CYFYNGEDIG",
	"PARTNERIAETH CYFYNGEDIG",
	"PCC LIMITED",
	"PCC LTD",
	"PCC",
	"PLC",
	"PROTECTED CELL COMPANY",
	"PUBLIC LIMITED COMPANY .THE",
	"PUBLIC LIMITED COMPANY THE",
	"PUBLIC LIMITED COMPANY",
	"PUBLIC LIMITED COMPANY.THE",
	"UNLIMITED",
	"UNLTD",
}

type Formatter interface {
	SplitCompanyNameEndings(name string) (string, string)
}

type Format struct{}

// NewFormatter returns a concrete implementation of the Formatter interface
func NewFormatter() Formatter {

	return &Format{}
}

//SplitCompanyNameEndings splits company name into nameStart and nameEnding in order to remove common name endings
func (f *Format) SplitCompanyNameEndings(name string) (string, string) {
	var nameStart, nameEnding string

	nameStart = name

	//Strip trailing non-word characters [^a-zA-Z0-9_]
	stripped := nonWordEndRegex.ReplaceAllString(name, "")

	//Scan company name for name ending and remove suffix
	for _, cne := range companyNameEndings {
		if strings.HasSuffix(stripped, cne) {
			nameStart = strings.TrimSuffix(stripped, " "+cne)
			// Keep the actual name ending by extracted the name start
			nameEnding = name[len(nameStart):]
			break
		}
	}

	return nameStart, nameEnding
}
