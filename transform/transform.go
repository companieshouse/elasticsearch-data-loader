package transform

import (
	"fmt"
	"log"

	"github.com/companieshouse/elasticsearch-data-loader/datastructures"
	"github.com/companieshouse/elasticsearch-data-loader/format"
	"github.com/companieshouse/elasticsearch-data-loader/write"
)

const recordKind = "searchresults#company"

// Transformer provides an interface by which to transform data from one form to another
type Transformer interface {
	TransformMongoCompanyToEsCompany(source *datastructures.MongoCompany) *datastructures.EsCompany
}

// Transform provides a concrete implementation of the Transformer interface
type Transform struct {
	w write.Writer
	f format.Formatter
}

// NewTransformer returns a concrete implementation of the Transformer interface
func NewTransformer(writer write.Writer, formatter format.Formatter) Transformer {

	return &Transform{
		w: writer,
		f: formatter,
	}
}

// TransformMongoCompanyToEsCompany transforms a MongoCompany into its EsCompany counterpart
func (t *Transform) TransformMongoCompanyToEsCompany(source *datastructures.MongoCompany) *datastructures.EsCompany {
	if source.Data == nil {
		log.Printf("Missing company data element")
		return nil
	}

	if source.Data.CompanyName == "" {
		t.w.LogMissingCompanyName(source.ID)
		return nil
	}

	dest := datastructures.EsCompany{
		ID:          source.ID,
		CompanyType: source.Data.CompanyType,
		Kind:        recordKind,
		Links:       &datastructures.EsLinks{Self: fmt.Sprintf("/company/%s", source.ID)},
	}

	name := source.Data.CompanyName

	nameStart, nameEnding := t.f.SplitCompanyNameEndings(source.Data.CompanyName)

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
