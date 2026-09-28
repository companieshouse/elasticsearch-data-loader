package datastructures

// MongoLinks holds a set of links relevant to MongoData
type MongoLinks struct {
	Self string `bson:"self"`
}

// MongoData contains company data from MongoDB
type MongoData struct {
	CompanyName   string     `bson:"company_name"`
	CompanyNumber string     `bson:"company_number"`
	CompanyStatus string     `bson:"company_status"`
	CompanyType   string     `bson:"type"`
	Links         MongoLinks `bson:"links"`
}

// MongoCompany wraps MongoData with an accompanying ID
type MongoCompany struct {
	ID   string     `bson:"_id"`
	Data *MongoData `bson:"data"`
}
