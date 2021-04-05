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

	var idWiki string

	url := `https://query.wikidata.org/sparql?format=json&query=SELECT%20DISTINCT%20?item%20WHERE%20{?item%20wdt:P31%20wd:Q5.%20?item%20?label%20"` + authorName + `"@pt%20FILTER(BOUND(?item)).%20SERVICE%20wikibase:label%20{bd:serviceParam%20wikibase:language%20%22pt%22.}}`

	res, err := http.Get(url)

	if err != nil {
		// panic(err.Error())
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		// panic(err.Error())
		fmt.Println(err)
	}
	defer res.Body.Close()

	data := object{}
	errI := json.Unmarshal(body, &data)
	if errI != nil {
		fmt.Println(errI)
	}

	for _, p := range data.Results.Bindings {

		q := fmt.Sprintf("%v", p.Item.Value)
		idWiki = q[strings.LastIndex(q, "/")+1:]
	}

	return idWiki
}
