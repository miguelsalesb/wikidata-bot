package wiki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type WikiData struct {
	id      string
	mattype string
}

// Structs to get the titles Wikidata ID
type Title_object struct {
	Title_head
	Results Title_results
}

type Title_head struct {
	Vars Title_vars
}

type Title_vars struct {
	Vars []string `json:"type"`
}

type Title_results struct {
	Bindings []Title_bindings `json:"bindings"`
}

type Title_bindings struct {
	Item    Title_item    `json:"item"`
	Mattype Title_mattype `json:"mattype"`
	Author  Title_author  `json:"occurences"`
}

type Title_author struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Title_item struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Title_mattype struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func GetWikiTitleID(titleLowercase string, authors []string, language string) WikiData {
	var titleArray WikiData
	var idW, wMat, idWiki, wMaterial, ti, aut, url string
	idW, wMat, idWiki, wMaterial = "", "", "", ""
	const empty = ""

	ti = Replacer.Replace(titleLowercase)

	if titleLowercase != "" {
		if len(authors) > 0 {
			aut = Replacer.Replace(authors[1])

			url = `https://query.wikidata.org/sparql?format=json&query=SELECT%20?item%20?itemLabel%20?author%20?authorLabel%20{SERVICE%20wikibase:mwapi%20{bd:serviceParam%20wikibase:api%20%22EntitySearch%22.bd:serviceParam%20wikibase:endpoint%20%22www.wikidata.org%22.bd:serviceParam%20mwapi:search%20%22` + ti + `%22.bd:serviceParam%20mwapi:language%20%22en%22.?item%20wikibase:apiOutputItem%20mwapi:item.?num%20wikibase:apiOrdinal%20true.}?item%20?label%20?title;(wdt:P31)%20?mattype;(wdt:P50)%20?author.?author%20?label%20%22` + aut + `%22@` + language + `.SERVICE%20wikibase:label%20{bd:serviceParam%20wikibase:language%20%22` + language + `%22.}}ORDER%20BY%20?num%20LIMIT%201`

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

			data := object{}
			errI := json.Unmarshal(body, &data)
			if errI != nil {
				fmt.Println(errI)
			}
			for _, p := range data.Results.Bindings {
				q := fmt.Sprintf("%v", p.Item.Value)
				idW = q[strings.LastIndex(q, "/")+1:]

				m := fmt.Sprintf("%v", p.Mattype.Value)
				wMat = m[strings.LastIndex(m, "/")+1:]
			}

			if idW != "" {
				idWiki = idW
			} else {
				idWiki = empty
			}
			if wMat != "" {
				wMaterial = wMat
			} else {
				wMaterial = empty
			}

			titleArray = WikiData{
				idWiki,
				wMaterial,
			}

		} else {
			titleArray = WikiData{
				empty,
				empty,
			}
		}
	}
	return titleArray
}
