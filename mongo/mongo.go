package mongo



type MongoLinks struct {
	Self string `bson:"self"`
}

type MongoData struct {
	CompanyName   string     `bson:"company_name"`
	CompanyNumber string     `bson:"company_number"`
	CompanyStatus string     `bson:"company_status"`
	CompanyType   string     `bson:"type"`
	Links         MongoLinks `bson:"links"`
}

type MongoCompany struct {
	ID   string     `bson:"_id"`
	Data *MongoData `bson:"data"`
}