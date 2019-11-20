package mapping

import (
	"fmt"
	"github.com/companieshouse/elasticsearch-data-loader/datastructures"
	"github.com/companieshouse/elasticsearch-data-loader/format"
	"github.com/companieshouse/elasticsearch-data-loader/write"
	"log"
)

const recordKind = "searchresults#company"

type Mapper interface {
	MapResult(source *datastructures.MongoCompany) *datastructures.EsCompany
}

type Map struct {
	writer    write.Writer
	formatter format.Formatter
}

func NewMapper(w write.Writer, f format.Formatter) Mapper {

	return &Map{
		writer:    w,
		formatter: f,
	}
}

/*
Pass in a reference to mongoCompany, as golang is pass-by-value. This version, golang
will create a copy of mongoCompany on the stack for every call (which is good, as it
ensures immutability, but we want efficiency! Passing a ref to mongoCompany will be
MUCH quicker.
*/
func (m *Map) MapResult(source *datastructures.MongoCompany) *datastructures.EsCompany {
	if source.Data == nil {
		log.Printf("Missing company data element")
		return nil
	}

	if source.Data.CompanyName == "" {
		m.writer.WriteToFile3(source.ID)
		return nil
	}

	dest := datastructures.EsCompany{
		ID:          source.ID,
		CompanyType: source.Data.CompanyType,
		Kind:        recordKind,
		Links:       &datastructures.EsLinks{Self: fmt.Sprintf("/company/%s", source.ID)},
	}

	name := source.Data.CompanyName

	nameStart, nameEnding := m.formatter.SplitCompanyNameEndings(source.Data.CompanyName)

	items := datastructures.EsItem{
		CompanyStatus:       source.Data.CompanyStatus,
		CompanyNumber:       source.Data.CompanyNumber,
		CorporateName:       name,
		CorporateNameStart:  nameStart,
		CorporateNameEnding: nameEnding,
		RecordType:          "companies",
	}

	dest.Items = items

	return &dest
}
