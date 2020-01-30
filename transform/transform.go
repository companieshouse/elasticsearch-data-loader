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
	TransformMongoCompanyToEsCompany(mongoCompany *datastructures.MongoCompany, alphaKey *datastructures.AlphaKey) *datastructures.EsCompany
	GetCompanyNames(companies *[]*datastructures.MongoCompany, length int) []datastructures.CompanyName
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

// TransformMongoCompanyToEsCompany transforms a MongoCompany and its relevant AlphaKey into its EsCompany counterpart
func (t *Transform) TransformMongoCompanyToEsCompany(mongoCompany *datastructures.MongoCompany, alphaKey *datastructures.AlphaKey) *datastructures.EsCompany {
	if mongoCompany.Data == nil {
		log.Printf("Missing company data element")
		return nil
	}

	if mongoCompany.Data.CompanyName == "" {
		t.w.LogMissingCompanyName(mongoCompany.ID)
		return nil
	}

	dest := datastructures.EsCompany{
		ID:          mongoCompany.ID,
		CompanyType: mongoCompany.Data.CompanyType,
		Kind:        recordKind,
		Links:       &datastructures.EsLinks{Self: fmt.Sprintf("/company/%s", mongoCompany.ID)},
	}

	name := mongoCompany.Data.CompanyName

	nameStart, nameEnding := t.f.SplitCompanyNameEndings(mongoCompany.Data.CompanyName)

	items := datastructures.EsItem{
		CompanyStatus:       mongoCompany.Data.CompanyStatus,
		CompanyNumber:       mongoCompany.Data.CompanyNumber,
		CorporateName:       name,
		CorporateNameStart:  nameStart,
		CorporateNameEnding: nameEnding,
		RecordType:          "companies",
		AlphaKey:            alphaKey.SameAsAlphaKey,
		OrderedAlphaKey:     alphaKey.OrderedAlphaKey,
	}

	dest.Items = items

	return &dest
}

// GetCompanyNames returns a set of 'CompanyName's for a given set of 'MongoCompany's
func (t *Transform) GetCompanyNames(companies *[]*datastructures.MongoCompany, length int) []datastructures.CompanyName {

	var companyNames []datastructures.CompanyName
	for i := 0; i < length; i++ {
		companyNames = append(companyNames, datastructures.CompanyName{Name: (*companies)[i].Data.CompanyName})
	}

	return companyNames
}
