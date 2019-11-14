package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"regexp"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const recordKind = "searchresults#company"

var (
	alphakeyURL = "http://chs-alphakey-pp.internal.ch"
	esDestURL   = "http://localhost:9200"
	mongoURL    = "mongodb://envb.mongo.ch.gov.uk:15107"
)

var (
	numberRegex      = regexp.MustCompile(`\d+`)
	findPremiseRegex = regexp.MustCompile(`^(\d+[\/\-]*\w*\s*|[a-zA-Z]+\s+\d+[a-zA-Z]*\W*\s*)`)
	nonWordEndRegex  = regexp.MustCompile(`[^A-Za-z0-9_]+$`)
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
	wg sync.WaitGroup

	countCh  = make(chan int)
	insertCh = make(chan int)
	skipCh   = make(chan int)
	sem      = make(chan int, 5)
)

var companyNameEndings = [...]string{
	"AEIE",
	"ANGHYFYNGEDIG",
	"C.B.C",
	"C.C.C",
	"C.I.C",
	"CBC",
	"CBCN",
	"CBP",
	"CCC",
	"CCG CYF",
	"CCG CYFYNGEDIG",
	"CIC",
	"COMMUNITY INTEREST COMPANY",
	"COMMUNITY INTEREST P.L.C",
	"COMMUNITY INTEREST PLC",
	"COMMUNITY INTEREST PUBLIC LIMITED COMPANY",
	"CWMNI BUDDIANT C.C.C",
	"CWMNI BUDDIANT CCC",
	"CWMNI BUDDIANT CYMUNEDOL C.C.C",
	"CWMNI BUDDIANT CYMUNEDOL CCC",
	"CWMNI BUDDIANT CYMUNEDOL CYHOEDDUS CYFYNGEDIG",
	"CWMNI BUDDIANT CYMUNEDOL",
	"CWMNI BUDDSODDIA CHYFALAF NEWIDIOL",
	"CWMNI BUDDSODDIANT PENAGORED",
	"CWMNI CELL GWARCHODEDIG",
	"CWMNI CYFYNGEDIG CYHOEDDUS",
	"CYF",
	"CYFYNGEDIG",
	"EEIG",
	"EESV",
	"EOFG",
	"EOOS",
	"EUROPEAN ECONOMIC INTEREST GROUPING",
	"GEIE",
	"GELE",
	"ICVC",
	"INVESTMENT COMPANY WITH VARIABLE CAPITAL",
	"L.P",
	"L.T.D",
	"LIMITED - THE",
	"LIMITED LIABILITY PARTNERSHIP",
	"LIMITED PARTNERSHIP",
	"LIMITED THE",
	"LIMITED",
	"LIMITED-THE",
	"LIMITED...THE",
	"LIMITED..THE",
	"LIMITED.THE",
	"LLP",
	"LP",
	"LTD",
	"LTD...THE",
	"LTD..THE",
	"LTD.THE",
	"OEIC",
	"OPEN-ENDED INVESTMENT COMPANY",
	"P.L.C",
	"PAC",
	"PARTNERIAETH ATEBOLRWYDD CYFYNGEDIG",
	"PARTNERIAETH CYFYNGEDIG",
	"PCC LIMITED",
	"PCC LTD",
	"PCC",
	"PLC",
	"PROTECTED CELL COMPANY",
	"PUBLIC LIMITED COMPANY .THE",
	"PUBLIC LIMITED COMPANY THE",
	"PUBLIC LIMITED COMPANY",
	"PUBLIC LIMITED COMPANY.THE",
	"UNLIMITED",
	"UNLTD",
}

var (
	filename1 = "company-errors/error-posting-request.txt"
	filename2 = "company-errors/unexpected-put-response.txt"
	filename3 = "company-errors/missing-company-name.txt"
)

// ---------------------------------------------------------------------------
type mongoLinks struct {
	Self string `bson:"self"`
}

type mongoData struct {
	CompanyName   string     `bson:"company_name"`
	CompanyNumber string     `bson:"company_number"`
	CompanyStatus string     `bson:"company_status"`
	CompanyType   string     `bson:"type"`
	Links         mongoLinks `bson:"links"`
}

type mongoCompany struct {
	ID   string     `bson:"_id"`
	Data *mongoData `bson:"data"`
}

// ---------------------------------------------------------------------------

type esItem struct {
	CompanyNumber       string `json:"company_number"`
	CompanyStatus       string `json:"company_status,omitempty"`
	CorporateName       string `json:"corporate_name"`
	CorporateNameStart  string `json:"corporate_name_start"`
	CorporateNameEnding string `json:"corporate_name_ending,omitempty"`
	RecordType          string `json:"record_type"`
}

type esLinks struct {
	Self string `json:"self"`
}

type esCompany struct {
	id          string
	CompanyType string   `json:"company_type"`
	Items       esItem   `json:"items"`
	Kind        string   `json:"kind"`
	Links       *esLinks `json:"links"`
	SortKey     string   `json:"sort_key"`
}

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

type companyName struct {
	Name string `json:"name"`
}

type alphaKeys struct {
	SameAsKey string `json:"sameAsAlphaKey"`
	SortKey   string `json:"orderedAlphaKey"`
}

type connections struct {
	connection1 *os.File
	connection2 *os.File
	connection3 *os.File
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

	s, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Fatalf("error creating mongoDB session: %s", err)
	}

	connection1, err := os.OpenFile(filename1, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening [%s] file", filename1)
	}

	connection2, err := os.OpenFile(filename2, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening [%s] file", filename2)
	}

	connection3, err := os.OpenFile(filename3, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error opening [%s] file", filename3)
	}

	c := &connections{
		connection1: connection1,
		connection2: connection2,
		connection3: connection3,
	}

	go status()

	it := s.DB(mongoDatabase).C(mongoCollection).Find(bson.M{}).Batch(mongoSize).Iter()

	for {
		companies := make([]*mongoCompany, mongoSize)

		itx := 0
		for ; itx < len(companies); itx++ {
			result := mongoCompany{}

			if !it.Next(&result) {
				break
			}
			companies[itx] = &result
		}
		if itx == 0 { // No results read from iterator. Nothing more to do.
			time.Sleep(time.Second * 30)
			break
		}

		// This will block if we've reached our concurrecy limit (sem buffer size)
		c.sendToES(&companies, itx)
	}

	time.Sleep(20 * time.Second)
	wg.Wait()
	if err := connection1.Close(); err != nil {
		log.Fatalf("error closing file: %s", err)
	}
	if err := connection2.Close(); err != nil {
		log.Fatalf("error closing file: %s", err)
	}
	if err := connection3.Close(); err != nil {
		log.Fatalf("error closing file: %s", err)
	}

	log.Println("COMPANY SEARCH DONE")
}

// ---------------------------------------------------------------------------
// pass a reference to the slice of mongoCompany pointers, for efficiency,
// otherwise golang will create a copy of the slice on the stack!
//
func (c *connections) sendToES(companies *[]*mongoCompany, length int) {
	// Wait on semaphore if we've reached our concurrency limit
	wg.Add(1)
	sem <- 1

	go func() {
		defer func() {
			<-sem
			wg.Done()
		}()

		countCh <- length
		target := length

		var bulk []byte
		var bunchOfNamesAndNumbers []byte

		// Get a list of company names for bulk alpha key api
		companyNames := getCompanyNames(companies, length)

		keys := corporateAlphaKeys(companyNames, length)

		var j []alphaKeys
		if err := json.Unmarshal(keys, &j); err != nil {
			log.Printf("error unmarshalling alphakey response for %s", keys)
		}

		i := 0
		for i < length {
			sameAsKey, sortKey := j[i].SameAsKey, j[i].SortKey
			company := c.mapResult((*companies)[i], sameAsKey, sortKey)

			if company != nil {
				b, err := json.Marshal(company)
				if err != nil {
					log.Fatalf("error marshal to json: %s", err)
				}

				bulk = append(bulk, []byte("{ \"create\": { \"_id\": \""+company.id+"\" } }\n")...)
				bulk = append(bulk, b...)
				bulk = append(bulk, []byte("\n")...)
				bunchOfNamesAndNumbers = append(bunchOfNamesAndNumbers, []byte("\n"+company.SortKey+" : "+company.id+"")...)
			} else {
				skipCh <- 1
				target--
			}

			i++
		}

		r, err := http.Post(esDestURL+"/"+esDestIndex+"/_bulk", "application/json", bytes.NewReader(bulk))
		if err != nil {
			writeToFile(c.connection1, filename1, string(bunchOfNamesAndNumbers))
			log.Printf("error posting request %s: data %s", err, string(bulk))
			return
		}
		defer r.Body.Close()

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("error reading response body: %s", err)
		}

		if r.StatusCode > 299 {
			writeToFile(c.connection2, filename2, string(bunchOfNamesAndNumbers))
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

		insertCh <- target
	}()
}

// ---------------------------------------------------------------------------

func nameNibbles(name string) (string, string) {
	var nameStart, nameEnding string

	nameStart = name

	//Strip trailing non-word characters [^a-zA-Z0-9_]
	stripped := nonWordEndRegex.ReplaceAllString(name, "")

	//Scan company name for name ending and remove suffix
	for _, cne := range companyNameEndings {
		if strings.HasSuffix(stripped, cne) {
			nameStart = strings.TrimSuffix(stripped, " "+cne)
			// Keep the actual name ending by extracted the name start
			nameEnding = name[len(nameStart):]
			break
		}
	}

	return nameStart, nameEnding
}

// ---------------------------------------------------------------------------
// Pass in a reference to mongoCompany, as golang is pass-by-value. This version, golang
// will create a copy of mongoCompany on the stack for every call (which is good, as it
// ensures imutability, but we want efficiency! Passing a ref to mongoCompany will be
// MUCH quicker.
//
func (c *connections) mapResult(source *mongoCompany, sameAsKey string, sortKey string) *esCompany {
	if source.Data == nil {
		log.Printf("Missing company data element")
		return nil
	}

	if source.Data.CompanyName == "" {
		writeToFile(c.connection3, filename3, source.ID)
		return nil
	}

	dest := esCompany{
		id:          source.ID,
		CompanyType: source.Data.CompanyType,
		Kind:        recordKind,
		Links:       &esLinks{fmt.Sprintf("/company/%s", source.ID)},
	}

	name := source.Data.CompanyName
	nameStart, nameEnding := nameNibbles(source.Data.CompanyName)

	items := esItem{
		CompanyStatus:       source.Data.CompanyStatus,
		CompanyNumber:       source.Data.CompanyNumber,
		CorporateName:       name,
		CorporateNameStart:  nameStart,
		CorporateNameEnding: nameEnding,
		RecordType:          "companies",
	}

	// Appended sort key with value of 0 for all search; this enables companies
	// to appear ahead of disqualified officers (appended with 1) and officers
	// (appended with 2) if score of multiple results spanning the different data
	// sets are equal.
	dest.SortKey = sortKey + "0"
	dest.Items = items

	return &dest
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
		case n := <-skipCh:
			skipCounter += n
			skipTotal += n
		case n := <-countCh:
			rpsCounter += n
			reqTotal += n
		case n := <-insertCh:
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

func corporateAlphaKeys(corporateNames []string, length int) []byte {
	items := []companyName{}
	for i := 0; i < length; i++ {
		name := corporateNames[i]
		var newItem companyName
		if name != "" {
			newItem = companyName{Name: name}
		} else {
			newItem = companyName{Name: "no-name-test"}
		}
		items = append(items, newItem)
	}

	jsonStr, err := json.Marshal(items)
	if err != nil {
		log.Fatalf("error marshalling %s: %s", items, err)
	}

	uri := fmt.Sprintf("%s/alphakey-bulk", alphakeyURL)

	r, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatalf(`error: %s with items: \n %s`, items, err)
	}

	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading alphakey response body for %s: %s", items, err)
	}

	return b
}

// ------------------------------------------------------------------------------

func getCompanyNames(companies *[]*mongoCompany, length int) []string {
	var companyNameList []string
	for i := 0; i < length; i++ {
		companyName := (*companies)[i].Data.CompanyName
		companyNameList = append(companyNameList, companyName)
	}

	return companyNameList
}

// ------------------------------------------------------------------------------

func writeToFile(connection *os.File, location string, sentence string) {
	_, err := connection.WriteString(sentence + "\n")
	if err != nil {
		log.Printf("error writing [%s] to file location: [%s]", sentence, location)
	}
}
