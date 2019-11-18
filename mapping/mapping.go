package mapping

import (
	"fmt"
	"github.com/companieshouse/elasticsearch-data-loader/format"
	"github.com/companieshouse/elasticsearch-data-loader/mongo"
	"github.com/companieshouse/elasticsearch-data-loader/write"
	"log"
)

const recordKind = "searchresults#company"

type Mapping interface {
	MapResult(source *mongo.MongoCompany) *esCompany
}
type Mapper struct {
	Writer write.Write
}

type esCompany struct {
	Id          string
	CompanyType string   `json:"company_type"`
	Items       esItem   `json:"items"`
	Kind        string   `json:"kind"`
	Links       *esLinks `json:"links"`
}

type esItem struct {
	CompanyNumber       string `json:"company_number"`
	CompanyStatus       string `json:"company_status,omitempty"`
	CorporateName       string `json:"corporate_name"`
	CorporateNameStart  string `json:"corporate_name_start"`
	CorporateNameEnding string `json:"corporate_name_ending,omitempty"`
	RecordType          string `json:"record_type"`
}

type esLinks struct {
	Self string `json:"self"`
}
/*
Pass in a reference to mongoCompany, as golang is pass-by-value. This version, golang
will create a copy of mongoCompany on the stack for every call (which is good, as it
ensures immutability, but we want efficiency! Passing a ref to mongoCompany will be
MUCH quicker.
*/

func (m *Mapper) MapResult(source *mongo.MongoCompany) *esCompany {
	if source.Data == nil {
		log.Printf("Missing company data element")
		return nil
	}

	if source.Data.CompanyName == "" {
		m.Writer.WriteToFile3(source.ID)
		return nil
	}

	dest := esCompany{
		Id:          source.ID,
		CompanyType: source.Data.CompanyType,
		Kind:        recordKind,
		Links:       &esLinks{fmt.Sprintf("/company/%s", source.ID)},
	}

	name := source.Data.CompanyName

	f := &format.Format{}
	nameStart, nameEnding := f.SplitCompanyNameEndings(source.Data.CompanyName)

	items := esItem{
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

