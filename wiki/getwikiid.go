package wiki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Structs to get the author's Wikidata ID
type object struct {
	head
	Results results
}

type head struct {
	Vars vars
}

type vars struct {
	Vars []string
}

type results struct {
	Bindings []bindings `json:"bindings"`
}

type bindings struct {
	Item    item    `json:"item"`
	Mattype mattype `json:"mattype"`
	Author  author  `json:"occurences"`
}

type author struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type item struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type mattype struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// GetAuthorIDWiki - function that searches for an author and returns the its Wikidata ID
func GetAuthorIDWiki(authorName string) string {
	// time.Sleep(300 * time.Millisecond)

	var idWiki, idWikiPT, idWikiEN string

	urlPT := `https://query.wikidata.org/sparql?format=json&query=SELECT%20DISTINCT%20?item%20WHERE%20{?item%20wdt:P31%20wd:Q5.%20?item%20?label%20"` + authorName + `"@pt%20FILTER(BOUND(?item)).%20SERVICE%20wikibase:label%20{bd:serviceParam%20wikibase:language%20%22pt%22.}}`

	resPT, errPT := http.Get(urlPT)

	if errPT != nil {
		// panic(err.Error())
		fmt.Println(errPT)
	}

	bodyPT, errBPT := ioutil.ReadAll(resPT.Body)

	if errBPT != nil {
		// panic(err.Error())
		fmt.Println(errBPT)
	}
	defer resPT.Body.Close()

	dataPT := object{}
	errIPT := json.Unmarshal(bodyPT, &dataPT)
	if errIPT != nil {
		fmt.Println(errIPT)
	}

	urlEN := `https://query.wikidata.org/sparql?format=json&query=SELECT%20DISTINCT%20?item%20WHERE%20{?item%20wdt:P31%20wd:Q5.%20?item%20?label%20"` + authorName + `"@en%20FILTER(BOUND(?item)).%20SERVICE%20wikibase:label%20{bd:serviceParam%20wikibase:language%20%22en%22.}}`

	resEN, errEN := http.Get(urlEN)

	if errEN != nil {
		// panic(err.Error())
		fmt.Println(errEN)
	}

	bodyEN, errBEN := ioutil.ReadAll(resEN.Body)

	if errBEN != nil {
		// panic(err.Error())
		fmt.Println(errBEN)
	}
	defer resPT.Body.Close()

	dataEN := object{}
	errIEN := json.Unmarshal(bodyEN, &dataEN)
	if errIEN != nil {
		fmt.Println(errIEN)
	}

	for _, p := range dataPT.Results.Bindings {
		q := fmt.Sprintf("%v", p.Item.Value)
		idWikiPT = q[strings.LastIndex(q, "/")+1:]
	}

	for _, p := range dataEN.Results.Bindings {
		q := fmt.Sprintf("%v", p.Item.Value)
		idWikiEN = q[strings.LastIndex(q, "/")+1:]
	}

	if len(idWikiPT) > 0 {
		idWiki = idWikiPT
	} else if len(idWikiEN) > 0 {
		idWiki = idWikiEN
	}
	return idWiki
}
