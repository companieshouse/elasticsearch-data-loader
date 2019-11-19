package datastructures

type EsCompany struct {
	Id          string
	CompanyType string   `json:"company_type"`
	Items       EsItem   `json:"items"`
	Kind        string   `json:"kind"`
	Links       *EsLinks `json:"links"`
}

type EsItem struct {
	CompanyNumber       string `json:"company_number"`
	CompanyStatus       string `json:"company_status,omitempty"`
	CorporateName       string `json:"corporate_name"`
	CorporateNameStart  string `json:"corporate_name_start"`
	CorporateNameEnding string `json:"corporate_name_ending,omitempty"`
	RecordType          string `json:"record_type"`
}

type EsLinks struct {
	Self string `json:"self"`
}
