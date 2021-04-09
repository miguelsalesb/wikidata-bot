package wiki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"wikidata/functions"
)

// Structs to get the author's Wikidata properties values
type W_object struct {
	head
	Results W_results
}

type W_head struct {
	Vars W_vars
}

type W_vars struct {
	Vars []string `json:"type"`
}

type W_results struct {
	Bindings []W_bindings `json:"bindings"`
}

type W_bindings struct {
	Prop              W_prop              `json:"prop"`
	Val_              W_val_              `json:"val_"`
	ValLabel          W_val_Label         `json:"val_Label"`
	AuthorDescription W_authorDescription `json:"authorDescription"`
}

type W_prop struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type W_val_ struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type W_val_Label struct {
	Xmllang string `json:"xml:lang"`
	Type    string `json:"type"`
	Value   string `json:"value"`
}

type W_authorDescription struct {
	Xmllang string `json:"xml:lang"`
	Type    string `json:"type"`
	Value   string `json:"value"`
}

// GetWikiContent - function that sends a SPARQL query with the author name, and if the author exists in Wikidata, provides the author's data in the JSON format
func GetWikiContent(authors []string) ([]string, []string) {

	const empty = ""
	// Search the author ID in Wikidata to embed it on the SPARQL query

	var (
		idBNP, wBirthdate, wDeathdate, name1, name2, wikiAuthorDescription,
		qOccupation, occupationsFinal, notableworkFinal, wNationality, signature, image string
	)
	var wikiAuthors = make([]string, 0, 9) // To avoid errors when nothing was retrieved - Para resolver o problema de quando não recuperava nada no JSON
	var occupations, notablework []string
	var wikiChecked = make([]string, 0, 2)

	wikiAuthorDescription = empty

	var counter200, counter400 = 0, 0

	for x := 0; x < len(authors); x += 6 {
		replacer := strings.NewReplacer(" ", "%20", "À", "%C3%80", "Á", "%C3%81", "Â", "%C3%82", "Ã", "%C3%83", "Ä", "%C3%84", "Ç", "%C3%87", "È", "%C3%88",
			"É", "%C3%89", "Ê", "%C3%8A", "Ë", "%C3%8B", "Ì", "%C3%8C", "Í", "%C3%8D", "Î", "%C3%8E", "Ï", "%C3%8F", "Ò", "%C3%92", "Ó", "%C3%93", "Ô", "%C3%94",
			"Õ", "%C3%95", "Ö", "%C3%96", "Ù", "%C3%99", "Ó", "%C3%9A", "Û", "%C3%9B", "Ý", "%C3%9D", "à", "%C3%A0", "á", "%C3%A1", "â", "%C3%A2", "ã", "%C3%A3",
			"ä", "%C3%A4", "ç", "%C3%A7", "è", "%C3%A8", "é", "%C3%A9", "ê", "%C3%AA", "ë", "%C3%AB", "ì", "%C3%AC", "í", "%C3%AD", "î", "%C3%AE", "ï", "C3%AF",
			"ñ", "%C3%B1", "ò", "%C3%B2", "ó", "%C3%B3", "ô", "%C3%B4", "õ", "%C3%B5", "ö", "%C3%B6", "ù", "%C3%B9", "ú", "%C3%BA", "û", "%C3%BB", "ü", "%C3%BC",
			"ý", "%C3%BD", "\"", "'", "º", "%C2%BA", "ª", "%C2%AA", "&", "%26", ",", "%2C", "!", "%21", "#", "%23", "$", "%24", "%", "%25", "'", "%27", "(", "%28",
			")", "%29", "-", "%2D", "[", "%5B", "]", "%5D", "^", "%5E", "_", "%5F", "_", "%60", "{", "%7B", "{", "%7C", "}", "%7D")
		name2 = replacer.Replace(authors[x+2])
		name1 = replacer.Replace(authors[x+1])

		// get the author's Wikidata ID
		var idWiki = GetAuthorIDWiki(name2 + "%20" + name1)

		if name2 != "" && name1 != "" {
			time.Sleep(300 * time.Millisecond)
			// SPARQL query with the author's Wikidata ID to get the values of the library ID, birthdate, deathdate, occupations, notable work and author description
			url := `https://query.wikidata.org/sparql?format=json&query=SELECT%20?prop%20?val_%20?val_Label%20?authorDescription%20{VALUES%20(?author)%20{(wd:` + idWiki + `)}%20?author%20?p%20?statement.%20?statement%20?val%20?val_.%20?prop%20wikibase:claim%20?p.%20?prop%20wikibase:statementProperty%20?val.%20SERVICE%20wikibase:label%20{bd:serviceParam%20wikibase:language%20%22pt%22%20}}%20ORDER%20BY%20?prop%20?statement%20?val_`

			res, err := http.Get(url)

			if err != nil {
				// panic(err.Error())
				fmt.Println(err)
			}

			body, err := ioutil.ReadAll(res.Body) // get Body data from the Webpage

			if err != nil {
				// panic(err.Error())
				fmt.Println(err)
			}
			defer res.Body.Close()

			data := W_object{}
			errP := json.Unmarshal(body, &data) // decode the JSON data
			if errP != nil {
				fmt.Println(errP)
			}

			if len(data.Results.Bindings) > 0 {
				for _, v := range data.Results.Bindings {
					// get the author's Library ID from Wikidata. Change the property to get theyour organization author's ID
					if v.Prop.Value == "http://www.wikidata.org/entity/P1005" {
						if idBNP == authors[0] {
							idBNP = authors[0]
						} else {
							idBNP = fmt.Sprintf("%v", v.ValLabel.Value)
						}
					}
					// get the author's Wikidata birth date
					if v.Prop.Value == "http://www.wikidata.org/entity/P569" {
						birth := fmt.Sprintf("%v", v.ValLabel.Value)
						wBirthdate = birth[:4]
					}
					// get the author's Wikidata death date
					if v.Prop.Value == "http://www.wikidata.org/entity/P570" {
						death := fmt.Sprintf("%v", v.ValLabel.Value)
						wDeathdate = death[:4]
					}
					// get the author's Wikidata nationality
					if v.Prop.Value == "http://www.wikidata.org/entity/P27" {
						nat := fmt.Sprintf("%v", v.ValLabel.Value)
						wNationality = nat
					}
					// get the author's signature image
					if v.Prop.Value == "http://www.wikidata.org/entity/P109" {
						sig := fmt.Sprintf("%v", v.ValLabel.Value)

						signature = sig
						// fmt.Println("\n\n\nsignature", signature)
					}
					// get the author's image
					if v.Prop.Value == "http://www.wikidata.org/entity/P18" {
						img := fmt.Sprintf("%v", v.ValLabel.Value)
						image = img
					}

					if v.Prop.Value == "http://www.wikidata.org/entity/P106" {
						q := fmt.Sprintf("%v", v.Val_.Value)
						qOccupation = q[strings.LastIndex(q, "/")+1:]
						occupations = append(occupations, fmt.Sprintf("%v", qOccupation))      // id of the occupation
						occupations = append(occupations, fmt.Sprintf("%v", v.ValLabel.Value)) // name of the occupation
					}

					// get the author's notable work
					if v.Prop.Value == "http://www.wikidata.org/entity/P800" {
						replacer := strings.NewReplacer(",", "", "\\", "", "'", "\\'")
						q := fmt.Sprintf("%v", v.Val_.Value)
						qNotableWork := q[strings.LastIndex(q, "/")+1:]
						notableWk := replacer.Replace(qNotableWork)
						notablework = append(notablework, fmt.Sprintf("%v", notableWk)) // id of the occupation
					}
					if v.AuthorDescription.Value != "" {
						replacer := strings.NewReplacer(",", "", "\\", "", "'", "\\'")
						autDesc := fmt.Sprintf("%v", v.AuthorDescription.Value)
						wikiAuthorDescription = replacer.Replace(autDesc)

					} else {
						wikiAuthorDescription = empty
					}

				}

				if idWiki != "" {
					if x == 0 {
						counter200 = 1
					} else if x > 0 {
						counter400++
					}
					wikiAuthors = append(wikiAuthors, idWiki)
				} else {
					wikiAuthors = append(wikiAuthors, empty)
				}

				if idBNP != "" {
					wikiAuthors = append(wikiAuthors, idBNP)
				} else {
					wikiAuthors = append(wikiAuthors, empty)
				}
				if wBirthdate != "" {
					wikiAuthors = append(wikiAuthors, wBirthdate)
				} else {
					wikiAuthors = append(wikiAuthors, empty)
				}
				if wDeathdate != "" {
					wikiAuthors = append(wikiAuthors, wDeathdate)
				} else {
					wikiAuthors = append(wikiAuthors, empty)
				}
				if wNationality != "" {
					wikiAuthors = append(wikiAuthors, wNationality)
				} else {
					wikiAuthors = append(wikiAuthors, empty)
				}
				if signature != "" {
					wikiAuthors = append(wikiAuthors, signature)
				} else {
					wikiAuthors = append(wikiAuthors, empty)
				}
				if image != "" {
					wikiAuthors = append(wikiAuthors, image)
				} else {
					wikiAuthors = append(wikiAuthors, empty)
				}

				occupations = functions.Unique(occupations)
				occupationsFinal = strings.Join(occupations, ", ")

				if len(occupationsFinal) > 0 {
					wikiAuthors = append(wikiAuthors, occupationsFinal)
				} else {
					wikiAuthors = append(wikiAuthors, empty)
				}

				notablework = functions.Unique(notablework)
				notableworkFinal = strings.Join(notablework, ", ")

				if len(notableworkFinal) > 0 {
					wikiAuthors = append(wikiAuthors, notableworkFinal)
				} else {
					wikiAuthors = append(wikiAuthors, empty)
				}

				// clean the arrays so that it won't append the previous data
				occupations, notablework = nil, nil

				if len(wikiAuthorDescription) > 0 {
					wikiAuthors = append(wikiAuthors, wikiAuthorDescription)
				} else {
					wikiAuthors = append(wikiAuthors, empty)
				}

			} else {
				wikiAuthors = append(wikiAuthors, empty)
				wikiAuthors = append(wikiAuthors, empty)
				wikiAuthors = append(wikiAuthors, empty)
				wikiAuthors = append(wikiAuthors, empty)
				wikiAuthors = append(wikiAuthors, empty)
				wikiAuthors = append(wikiAuthors, empty)
				wikiAuthors = append(wikiAuthors, empty)
				wikiAuthors = append(wikiAuthors, empty)
				wikiAuthors = append(wikiAuthors, empty)
				wikiAuthors = append(wikiAuthors, empty)

			}
		} else {
			wikiAuthors = append(wikiAuthors, empty)
			wikiAuthors = append(wikiAuthors, empty)
			wikiAuthors = append(wikiAuthors, empty)
			wikiAuthors = append(wikiAuthors, empty)
			wikiAuthors = append(wikiAuthors, empty)
			wikiAuthors = append(wikiAuthors, empty)
			wikiAuthors = append(wikiAuthors, empty)
			wikiAuthors = append(wikiAuthors, empty)
			wikiAuthors = append(wikiAuthors, empty)
			wikiAuthors = append(wikiAuthors, empty)

		}
	}

	// find out when our alternative name is filled in Wikidata instead of the header
	if (counter200 == 1 && counter400 == 0) || (counter200 == 1 && counter400 >= 1) {
		wikiChecked = append(wikiChecked, "true")
		// The corret name (according to our authorities) is filled in Wikidata")
	} else if counter200 == 1 && counter400 >= 1 {
		wikiChecked = append(wikiChecked, "false")
		// A alternative name is filled in Wikidata instead of the corret name")
	} else {
		wikiChecked = append(wikiChecked, empty)
	}
	// Find out if there is no record in Wikidata
	if counter200 == 0 && counter400 == 0 {
		wikiChecked = append(wikiChecked, "true")
	} else {
		wikiChecked = append(wikiChecked, "false")
	}
	return wikiAuthors, wikiChecked
}
