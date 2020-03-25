package datastructures

// EsCompany holds a set of items containing company data relevant to Elastic Search
type EsCompany struct {
	ID                    string
	CompanyType           string   `json:"company_type"`
	Items                 EsItem   `json:"items"`
	Kind                  string   `json:"kind"`
	Links                 *EsLinks `json:"links"`
	OrderedAlphaKeyWithID string   `json:"ordered_alpha_key_with_id"`
}

// EsItem holds an individual company's data
type EsItem struct {
	CompanyNumber       string `json:"company_number"`
	CompanyStatus       string `json:"company_status,omitempty"`
	CorporateName       string `json:"corporate_name"`
	CorporateNameStart  string `json:"corporate_name_start"`
	CorporateNameEnding string `json:"corporate_name_ending,omitempty"`
	RecordType          string `json:"record_type"`
	AlphaKey            string `json:"alpha_key"`
	OrderedAlphaKey     string `json:"ordered_alpha_key"`
}

// EsLinks holds a set of links relevant to an EsCompany
type EsLinks struct {
	Self string `json:"self"`
}
