package main

import (
	"encoding/json"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/companieshouse/elasticsearch-data-loader/datastructures"
	"github.com/companieshouse/elasticsearch-data-loader/eshttp"
	"github.com/companieshouse/elasticsearch-data-loader/format"
	"github.com/companieshouse/elasticsearch-data-loader/transform"
	"github.com/companieshouse/elasticsearch-data-loader/write"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoTimeout = time.Duration(5) * time.Second

var (
	alphakeyURL = "http://chs-alphakey-pp.internal.ch"
	esDestURL   = "http://localhost:9200"
	mongoURL    = "mongodb://envb.mongo.ch.gov.uk:15107"
)

var (
	mongoDatabase   = "company_profile"
	mongoCollection = "company_profile"
	mongoSize       = 500
)

var (
	esDestIndex = "companies"
	esDestType  = "company"
)

var (
	syncWaitGroup sync.WaitGroup

	countChannel  = make(chan int)
	insertChannel = make(chan int)
	skipChannel   = make(chan int)
	semaphore     = make(chan int, 5)
)

// Function variables to facilitate testing.
var (
	marshal   = json.Marshal
	unmarshal = json.Unmarshal
	fatalf    = log.Fatalf
)

// ---------------------------------------------------------------------------

type esBulkResponse struct {
	Took   int                  `json:"took"`
	Errors bool                 `json:"errors"`
	Items  []esBulkItemResponse `json:"items"`
}

type esBulkItemResponse map[string]esBulkItemResponseData

type esBulkItemResponseData struct {
	Index  string `json:"_index"`
	ID     string `json:"_id"`
	Status int    `json:"status"`
	Error  string `json:"error"`
}

// ---------------------------------------------------------------------------

func main() {
	flag.StringVar(&mongoURL, "mongo-url", mongoURL, "mongoDB URL")
	flag.StringVar(&mongoDatabase, "mongo-database", mongoDatabase, "mongoDB database")
	flag.StringVar(&mongoCollection, "mongo-collection", mongoCollection, "mongoDB collection")
	flag.IntVar(&mongoSize, "mongo-source-size", mongoSize, "mongo page size")
	flag.StringVar(&esDestURL, "es-dest-url", esDestURL, "elasticsearch destination URL")
	flag.StringVar(&esDestIndex, "es-dest-index", esDestIndex, "elasticsearch destination index")
	flag.StringVar(&esDestType, "es-dest-type", esDestType, "elasticsearch destination type")
	flag.StringVar(&alphakeyURL, "alphakey-url", alphakeyURL, "alphakey service url")
	flag.Parse()

	w := write.NewWriter()
	f := format.NewFormatter()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		fatalf("error creating mongoDB session: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
	defer cancel()
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			fatalf("error disconnecting from client: %s", err)
		}
	}(client, ctx)

	go status()
	companyProfileCollection := client.Database(mongoDatabase).Collection(mongoCollection)
	findOptions := options.Find()
	findOptions.SetBatchSize(int32(mongoSize))
	ctx2, cancel2 := context.WithTimeout(context.Background(), mongoTimeout)
	defer cancel2()
	cur, err := companyProfileCollection.Find(ctx2, bson.D{}, findOptions)
	if err != nil {
		fatalf("error reading from collection: %s", err)
	}

	ctx3, cancel3 := context.WithCancel(context.Background())
	defer cancel3()

	sendCompaniesToES(cur, ctx3, err, w, f)

	time.Sleep(5 * time.Second)
	syncWaitGroup.Wait()

	w.Close()

	log.Println("SUCCESSFULLY LOADED: company data to alpha_search index")
}

func sendCompaniesToES(cur *mongo.Cursor, ctx3 context.Context, err error, w write.Writer, f format.Formatter) {
	for {
		companies := make([]*datastructures.MongoCompany, mongoSize)
		itx := 0
		for ; itx < len(companies); itx++ {
			if !cur.Next(ctx3) {
				break
			}
			result := datastructures.MongoCompany{}
			if err = cur.Decode(&result); err != nil {
				log.Fatal(err)
			}
			companies[itx] = &result
		}

		if err := cur.Err(); err != nil {
			fatalf("error iterating the collection: %s", err)
		}

		// No results read from iterator. Nothing more to do.
		if itx == 0 {
			break
		}

		// This will block if we've reached our concurrency limit (sem buffer size)
		sendToES(&companies, itx, w, f)
	}
}

// ---------------------------------------------------------------------------

/*
 pass a reference to the slice of mongoCompany pointers, for efficiency,
 otherwise golang will create a copy of the slice on the stack!
*/

func sendToES(companies *[]*datastructures.MongoCompany, length int, w write.Writer, f format.Formatter) {

	// Wait on semaphore if we've reached our concurrency limit
	syncWaitGroup.Add(1)
	semaphore <- 1

	t := transform.NewTransformer(w, f)
	c := eshttp.NewClient(w)

	go func() {
		defer func() {
			<-semaphore
			syncWaitGroup.Done()
		}()

		countChannel <- length
		target := length

		var bulk []byte
		var companyNumbers []byte

		err, alphaKeys := getAlphaKeys(t, companies, length, c)

		bulk, companyNumbers, target =
			transformMongoCompaniesToEsCompanies(
				length,
				t,
				companies,
				alphaKeys,
				bulk,
				companyNumbers,
				target)

		if submitBulkToES(err, c, bulk, companyNumbers) {
			return
		}

		insertChannel <- target
	}()
}

func submitBulkToES(err error, c eshttp.Client, bulk []byte, companyNumbers []byte) bool {
	b, err := c.SubmitBulkToES(bulk, companyNumbers, esDestURL, esDestIndex)
	if err != nil {
		return true
	}

	var bulkRes esBulkResponse
	if err := unmarshal(b, &bulkRes); err != nil {
		fatalf("error unmarshalling json: [%s] actual response: [%s]", err, b)
	}

	if bulkRes.Errors {
		for _, r := range bulkRes.Items {
			if r["create"].Status != 201 {
				fatalf("error inserting doc: %s", r["create"].Error)
			}
		}
	}
	return false
}

func getAlphaKeys(
	t transform.Transformer,
	companies *[]*datastructures.MongoCompany,
	length int,
	c eshttp.Client) (error, []datastructures.AlphaKey) {
	companyNames := t.GetCompanyNames(companies, length)
	compNamesBody, err := marshal(companyNames)
	if err != nil {
		fatalf("error marshal to json: %s", err)
	}

	keys, err := c.GetAlphaKeys(compNamesBody, alphakeyURL)
	if err != nil {
		fatalf("error fetching alpha keys: %s", err)
	}

	var alphaKeys []datastructures.AlphaKey
	if err := unmarshal(keys, &alphaKeys); err != nil {
		fatalf("error %v unmarshalling alphakey response for %s", err, compNamesBody)
	}
	return err, alphaKeys
}

func transformMongoCompaniesToEsCompanies(
	length int,
	t transform.Transformer,
	companies *[]*datastructures.MongoCompany,
	alphaKeys []datastructures.AlphaKey,
	bulk []byte,
	companyNumbers []byte,
	target int) ([]byte, []byte, int) {
	i := 0
	for i < length {
		company := t.TransformMongoCompanyToEsCompany((*companies)[i], &alphaKeys[i])

		if company != nil {
			b, err := marshal(company)
			if err != nil {
				fatalf("error marshal to json: %s", err)
			}

			bulk = append(bulk, []byte("{ \"create\": { \"_id\": \""+company.ID+"\" } }\n")...)
			bulk = append(bulk, b...)
			bulk = append(bulk, []byte("\n")...)
			companyNumbers = append(companyNumbers, []byte("\n"+company.ID+"")...)
		} else {
			skipChannel <- 1
			target--
		}

		i++
	}
	return bulk, companyNumbers, target
}

// ---------------------------------------------------------------------------

func status() {
	var (
		rpsCounter  = 0
		insCounter  = 0
		skipCounter = 0
		reqTotal    = 0
		insTotal    = 0
		skipTotal   = 0
	)

	t := time.NewTicker(time.Second)

	for {
		select {
		case n := <-skipChannel:
			skipCounter += n
			skipTotal += n
		case n := <-countChannel:
			rpsCounter += n
			reqTotal += n
		case n := <-insertChannel:
			insCounter += n
			insTotal += n
		case <-t.C:
			log.Printf("Read: %6d  Written: %6d  Skipped: %6d  |  rps: %6d  ips: %6d  sps: %6d", reqTotal, insTotal, skipTotal, rpsCounter, insCounter, skipCounter)
			rpsCounter = 0
			insCounter = 0
			skipCounter = 0
		}
	}
}

// ------------------------------------------------------------------------------
