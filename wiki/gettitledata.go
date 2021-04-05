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

	replacer := strings.NewReplacer(" ", "%20", "à", "%C3%9F", "á", "%C3%A1", "â", "%C3%A2", "ã", "%C3%A3", "ä", "%C3%A4", "ç", "%C3%A7", "è", "%C3%A8", "é", "%C3%A9", "ê", "%C3%AA", "ë", "%C3%AB", "ì", "%C3%AC", "í", "%C3%AD", "î", "%C3%AE", "ï", "C3%AF", "ñ", "%C3%B1", "ò", "%C3%B2", "ó", "%C3%B3", "ô", "%C3%B4", "õ", "%C3%B5", "ö", "%C3%B6", "ù", "%C3%B9", "ú", "%C3%BA", "û", "%C3%BB", "ü", "%C3%BC", "ý", "%C3%BD", "\"", "'", "º", "%C2%BA", "ª", "%C2%AA", "&", "%26", ",", "%2C", "!", "%21", "#", "%23", "$", "%24", "%", "%25", "'", "%27", "(", "%28", ")", "%29", "-", "%2D", "[", "%5B", "]", "%5D", "^", "%5E", "_", "%5F", "_", "%60", "{", "%7B", "{", "%7C", "}", "%7D")
	ti = replacer.Replace(titleLowercase)

	if titleLowercase != "" {
		if len(authors) > 0 {
			aut = replacer.Replace(authors[1])

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
