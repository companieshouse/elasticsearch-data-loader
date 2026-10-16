package format

import (
	"regexp"
	"strings"
)

var (
	wordRegex      = regexp.MustCompile(`\S+`)
	namePartsRegex = regexp.MustCompile(`(.+?)(\W*)$`)
)

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

var companyNameEndingSet = func() map[string]struct{} {
	endings := make(map[string]struct{}, len(companyNameEndings))
	for _, ending := range companyNameEndings {
		endings[ending] = struct{}{}
	}
	return endings
}()

// Formatter provides an interface by which to perform string formatting operations
type Formatter interface {
	SplitCompanyNameEndings(name string) (string, string)
}

// Format provides a concrete implementation of the Formatter interface
type Format struct{}

// NewFormatter returns a concrete implementation of the Formatter interface
func NewFormatter() Formatter {

	return &Format{}
}

// SplitCompanyNameEndings splits company name into nameStart and nameEnding in order to remove common name endings
func (f *Format) SplitCompanyNameEndings(name string) (string, string) {
	groups := namePartsRegex.FindStringSubmatch(name)
	if len(groups) < 3 {
		return name, ""
	}

	clean := groups[1]
	allWords := wordRegex.FindAllString(clean, -1)

	for i := 1; i < len(allWords); i++ {
		// Try longest whole-word suffix first so we don't match partial endings (e.g. LTD inside UNLTD).
		candidate := strings.Join(allWords[i:], " ")
		if _, ok := companyNameEndingSet[candidate]; ok {
			nameStart := strings.TrimSpace(strings.TrimSuffix(clean, candidate))
			return nameStart, name[len(nameStart):]
		}
	}

	return name, ""
}
