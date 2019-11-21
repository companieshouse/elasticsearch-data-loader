package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/companieshouse/elasticsearch-data-loader/datastructures"
	"github.com/companieshouse/elasticsearch-data-loader/format"
	"github.com/companieshouse/elasticsearch-data-loader/transform"
	"github.com/companieshouse/elasticsearch-data-loader/write"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const applicationJSON = "application/json"

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

	s, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Fatalf("error creating mongoDB session: %s", err)
	}
	go status()

	it := s.DB(mongoDatabase).C(mongoCollection).Find(bson.M{}).Batch(mongoSize).Iter()

	for {
		companies := make([]*datastructures.MongoCompany, mongoSize)

		itx := 0
		for ; itx < len(companies); itx++ {
			result := datastructures.MongoCompany{}

			if !it.Next(&result) {
				break
			}
			companies[itx] = &result
		}
		// No results read from iterator. Nothing more to do.
		if itx == 0 {
			break
		}

		// This will block if we've reached our concurrency limit (sem buffer size)
		sendToES(&companies, itx, w, f)
	}

	time.Sleep(5 * time.Second)
	syncWaitGroup.Wait()

	w.Close()

	log.Println("SUCCESSFULLY LOADED: company data to alpha_search index")
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

	go func() {
		defer func() {
			<-semaphore
			syncWaitGroup.Done()
		}()

		countChannel <- length
		target := length

		var bulk []byte
		var bunchOfNamesAndNumbers []byte

		i := 0
		for i < length {
			company := t.TransformMongoCompanyToEsCompany((*companies)[i])

			if company != nil {
				b, err := json.Marshal(company)
				if err != nil {
					log.Fatalf("error marshal to json: %s", err)
				}

				bulk = append(bulk, []byte("{ \"create\": { \"_id\": \""+company.ID+"\" } }\n")...)
				bulk = append(bulk, b...)
				bulk = append(bulk, []byte("\n")...)
				bunchOfNamesAndNumbers = append(bunchOfNamesAndNumbers, []byte("\n"+company.ID+"")...)
			} else {
				skipChannel <- 1
				target--
			}

			i++
		}

		r, err := http.Post(esDestURL+"/"+esDestIndex+"/_bulk", applicationJSON, bytes.NewReader(bulk))
		if err != nil {
			w.LogPostError(string(bunchOfNamesAndNumbers))
			log.Printf("error posting request %s: data %s", err, string(bulk))
			return
		}
		defer r.Body.Close()

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("error reading response body: %s", err)
		}

		if r.StatusCode > 299 {
			w.LogUnexpectedResponse(string(bunchOfNamesAndNumbers))
			log.Printf("unexpected put response %s: data %s", r.Status, string(bulk))
			return
		}

		var bulkRes esBulkResponse
		if err := json.Unmarshal(b, &bulkRes); err != nil {
			log.Fatalf("error unmarshaling json: %s", err)
		}

		if bulkRes.Errors {
			for _, r := range bulkRes.Items {
				if r["create"].Status != 201 {
					log.Fatalf("error inserting doc: %s", r["create"].Error)
				}
			}
		}

		insertChannel <- target
	}()
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
