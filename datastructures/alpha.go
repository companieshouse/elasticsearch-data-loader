package datastructures

// AlphaKey holds data returned by the alpha key service
type AlphaKey struct {
	SameAsAlphaKey  string `json:"sameAsAlphaKey"`
	OrderedAlphaKey string `json:"orderedAlphaKey"`
}

// CompanyName holds the name of a company to be sent to the alpha key service
type CompanyName struct {
	Name string `json:"name"`
}
