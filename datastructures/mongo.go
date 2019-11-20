package datastructures

// MongoLinks holds a set of links relevant to a MongoData item
type MongoLinks struct {
	Self string `bson:"self"`
}

// MongoData holds relevant company data
type MongoData struct {
	CompanyName   string     `bson:"company_name"`
	CompanyNumber string     `bson:"company_number"`
	CompanyStatus string     `bson:"company_status"`
	CompanyType   string     `bson:"type"`
	Links         MongoLinks `bson:"links"`
}

// MongoCompany wraps a MongoData item
type MongoCompany struct {
	ID   string     `bson:"_id"`
	Data *MongoData `bson:"data"`
}
