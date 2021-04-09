package wiki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Structs to get the author's Wikidata ID
type O_object struct {
	head    `json:"head"`
	Results O_results `json:"results"`
}

type O_head struct {
	Vars O_vars `json:"vars"`
}

type O_vars struct {
	Vars []string `json:"type"`
}

type O_results struct {
	Bindings []O_bindings `json:"bindings"`
}

type O_bindings struct {
	Item            O_item            `json:"item"`
	InstanceOf      O_instanceOf      `json:"instanceOf"`
	InstanceOfLabel O_instanceOfLabel `json:"instanceOfLabel"`
}

type O_item struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type O_instanceOf struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type O_instanceOfLabel struct {
	Xmllang string `json:"xml:lang"`
	Type    string `json:"type"`
	Value   string `json:"value"`
}

// GetOccupationsWiki - get the Wikidata ID's of the noncoincidential occupations (the author's occupations that exist in the author's Library occupations field (Unimarc 830 field) and not in the Wikidata occupations values)
// sends the occupation and its label and the occupation, its label and the occupation instanceof and its label
func GetOccupationsWiki(nonCoincidentalOccupations []string) ([]string, []string) {

	const empty = ""
	var occupArray = make([]string, 0, 4)
	var occupArrayWithoutInstanceOf = make([]string, 0, 2)

	for x := 0; x < len(nonCoincidentalOccupations); x++ {

		if nonCoincidentalOccupations[x] != "" {
			// time.Sleep(300 * time.Millisecond)
			noCoincOccup := strings.TrimLeft(nonCoincidentalOccupations[x], " ")
			nonCoincOccup := replacer.Replace(noCoincOccup)

			// Just search the words that have more than 3 letters
			if len(nonCoincOccup) > 3 {

				// Search for the occupation (in portuguese) and that has P31 (instance of) with: Q28640 (profession), Q4164871 (position) or Q12737077 (occupation)

				url := `https://query.wikidata.org/sparql?format=json&query=SELECT%20?item%20?instanceOf%20?instanceOfLabel%20WHERE%20{SERVICE%20wikibase:mwapi%20{bd:serviceParam%20wikibase:endpoint%20%22www.wikidata.org%22;wikibase:api%20%22EntitySearch%22;mwapi:search%20"` + strings.TrimLeft(nonCoincOccup, " ") + `";mwapi:language%20%22pt%22.?item%20wikibase:apiOutputItem%20mwapi:item.?num%20wikibase:apiOrdinal%20true.}?item%20(wdt:P31)%20?instanceOf.SERVICE%20wikibase:label%20{%20bd:serviceParam%20wikibase:language%20%22pt%22.%20}}%20ORDER%20BY%20ASC(?num)%20LIMIT%201`
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

				data := O_object{}
				errP := json.Unmarshal(body, &data)
				if err != nil {
					fmt.Println(errP)
				}

				if len(data.Results.Bindings) <= 0 {
					occupArray = append(occupArray, "none")
					occupArray = append(occupArray, "none")
					occupArray = append(occupArray, "none")
					occupArray = append(occupArray, "none")
					occupArrayWithoutInstanceOf = append(occupArrayWithoutInstanceOf, "none")
					occupArrayWithoutInstanceOf = append(occupArrayWithoutInstanceOf, "none")
				} else {
					for _, v := range data.Results.Bindings {
						q := fmt.Sprintf("%v", v.Item.Value)

						qOccupation := q[strings.LastIndex(q, "/")+1:]
						occupArray = append(occupArray, qOccupation)
						occupArray = append(occupArray, strings.ReplaceAll(nonCoincOccup, "%20", " "))

						occupArrayWithoutInstanceOf = append(occupArrayWithoutInstanceOf, qOccupation)
						occupArrayWithoutInstanceOf = append(occupArrayWithoutInstanceOf, strings.ReplaceAll(nonCoincOccup, "%20", " "))

						instanceOf := fmt.Sprintf("%v", v.InstanceOf.Value)
						qInstanceOf := instanceOf[strings.LastIndex(instanceOf, "/")+1:]

						occupArray = append(occupArray, qInstanceOf)

						instanceOfLabel := fmt.Sprintf("%v", v.InstanceOfLabel.Value)

						occupArray = append(occupArray, instanceOfLabel)

						if q == "" {
							occupArray = append(occupArray, "none")
							occupArray = append(occupArray, "none")
							occupArray = append(occupArray, "none")
							occupArray = append(occupArray, "none")
							occupArrayWithoutInstanceOf = append(occupArrayWithoutInstanceOf, "none")
							occupArrayWithoutInstanceOf = append(occupArrayWithoutInstanceOf, "none")
						}

					}
				}

			}
		}
	}

	if len(occupArrayWithoutInstanceOf) < 2 {
		occupArrayWithoutInstanceOf = append(occupArrayWithoutInstanceOf, empty)
	} else {
		// If a occupation is more complete, ex.: "artista plastico" and "artista", add only "artista plastico"
		for i := 1; i < len(occupArrayWithoutInstanceOf)-1; i += 2 {

			if i <= len(occupArrayWithoutInstanceOf[i])-2 {
				for j := 1; j < len(occupArrayWithoutInstanceOf)-1; j += 2 {

					if len(occupArrayWithoutInstanceOf[i]) > len(occupArrayWithoutInstanceOf[j+2]) && strings.Contains(occupArrayWithoutInstanceOf[i], occupArrayWithoutInstanceOf[j+2]) {
						occupArrayWithoutInstanceOf[j+1] = ""
						occupArrayWithoutInstanceOf[j+2] = ""
					} else if len(occupArrayWithoutInstanceOf[i]) < len(occupArrayWithoutInstanceOf[j+2]) && strings.Contains(occupArrayWithoutInstanceOf[j+2], occupArrayWithoutInstanceOf[i]) {
						occupArrayWithoutInstanceOf[i-1] = ""
						occupArrayWithoutInstanceOf[i] = ""

					}

				}
			}
		}
	}
	return occupArray, occupArrayWithoutInstanceOf
}
