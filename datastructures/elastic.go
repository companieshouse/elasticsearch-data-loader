package datastructures

// EsCompany holds elastic search company data for a set of EsItems
type EsCompany struct {
	ID          string
	CompanyType string   `json:"company_type"`
	Items       EsItem   `json:"items"`
	Kind        string   `json:"kind"`
	Links       *EsLinks `json:"links"`
}

// EsItem holds elastic search data for each company
type EsItem struct {
	CompanyNumber       string `json:"company_number"`
	CompanyStatus       string `json:"company_status,omitempty"`
	CorporateName       string `json:"corporate_name"`
	CorporateNameStart  string `json:"corporate_name_start"`
	CorporateNameEnding string `json:"corporate_name_ending,omitempty"`
	RecordType          string `json:"record_type"`
}

// EsLinks holds a set of links relevant to the EsCompany
type EsLinks struct {
	Self string `json:"self"`
}
