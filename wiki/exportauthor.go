package wiki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"wikidata/db"
	"wikidata/functions"

	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/net/publicsuffix"
)

type A_object struct {
	A_Entity  A_Entity `json:"entity"`
	A_Success int      `json:"success"`
}

type A_Entity struct {
	A_Label       A_Labels       `json:"labels,omitempty"`
	A_Alias       A_Aliases      `json:"aliases,omitempty"`
	A_Description A_Descriptions `json:"descriptions,omitempty"`
	A_Claim       A_Claims       `json:"claims,omitempty"`
	A_Id          string         `json:"id,omitempty"`
	A_Type        string         `json:"type,omitempty"`
	A_Lastrevid   int            `json:"lastrevid,omitempty"`
}

type A_Claims struct {
	P5 []PItem   `json:"P5,omitempty"` // Entity type
	P8 []PItem   `json:"P8,omitempty"` // Entity type
	P1 []*PTime  `json:"P1,omitempty"` // Time type
	P2 []*PTime  `json:"P2,omitempty"` // Time type
	P6 []PString `json:"P6,omitempty"` // Entity type
}

type A_Aliases struct {
	A_Pt []A_Pts `json:"pt,omitempty"`
}

type A_Labels struct {
	A_Pt A_Pts `json:"pt,omitempty"`
	A_En A_Ens `json:"en,omitempty"`
}

type A_Descriptions struct {
	A_Pt *A_Pts `json:"pt,omitempty"`
}

type A_Pts struct {
	A_Language string `json:"language,omitempty"`
	A_Value    string `json:"value,omitempty"`
}
type A_Ens struct {
	A_Language string `json:"language,omitempty"`
	A_Value    string `json:"value,omitempty"`
}

// Entities and properties used to insert the author's and titles data
// P1 = date of birth (P569)
// P2 = date of death (P570)
// P3 afirmado em - stated in (P248)
// P4 endereço eletrónico da referência - reference URL (P854)
// P5 instância de (P5) - instance of (P31)
// P6 identificador PTBNP - Portuguese National Library ID (P1005)
// P7 data de acesso (P7) - retrieved (P813)
// P8 obra destacada - notable work (P800)
// P9 data de publicação - publication date (P577)
// P10 país de origem - country of origin (P495)
// P11 país de nacionalidade - country of citizenship (P27)
// Q1 ser humano (Q1) - human (Q5)
// Q2 BNP - National Library of Portugal (Q245966)
// Q3 obra escrita - written work (Q47461344)
// Q4 Portugal - Portugal (Q45)

func ExportAuthor(authorWikiArray []string) {

	var (
		greaterThanOne       bool
		exists_in_wiki_array []int
		tokenCsfr            string
		counter              int
	)

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar}

	tokenCsfr = ConnectToWikidata(client)

	for x := 0; x < len(authorWikiArray); x += 9 {

		i, _ := strconv.Atoi(authorWikiArray[x+8])

		// If any of the exists_in_wiki values is greater than 1, then the entity should not be exported do Wikidata
		if i > 1 {
			greaterThanOne = true
		}

		// only append the entries that have tle lowest probability of existing in Wikidata
		if counter > 0 && authorWikiArray[x+8] == "0" || authorWikiArray[x+8] == "1" {
			exists_in_wiki_array = append(exists_in_wiki_array, i)
		}

		// append the 200 field value
		if counter == 0 {
			exists_in_wiki_array = append(exists_in_wiki_array, i)
			counter++
		}
	}

	if greaterThanOne == false {
		// If field 200 has zero probability of existing in Wikidata
		if exists_in_wiki_array[0] == 0 {
			ExportToWiki(client, tokenCsfr, authorWikiArray)
		} else if len(exists_in_wiki_array) > 0 {
			for x := 1; x < len(exists_in_wiki_array); x++ {
				if exists_in_wiki_array[x] == 0 {
					ExportToWiki(client, tokenCsfr, authorWikiArray)
					break
				}
			}
		}
	}

}

func ExportToWiki(client http.Client, tokenCsfr string, authorWikiArray []string) {
	// time.Sleep(300 * time.Millisecond)
	var (
		id_library, name, birth_date_library, death_date_library, nationality, occupations_library, field, retrieved_date string
		entity                                                                                                            A_Entity
		replacerWiki                                                                                                      = strings.NewReplacer("\\", "", "\"[", "[", "]\"", "]")
	)

	for x := 0; x < len(authorWikiArray); x += 9 {
		id_library = authorWikiArray[x]
		name = authorWikiArray[x+1]
		birth_date_library = authorWikiArray[x+2]
		death_date_library = authorWikiArray[x+3]
		nationality = authorWikiArray[x+4]
		occupations_library = authorWikiArray[x+5]
		field = authorWikiArray[x+6]
		retrieved_date = authorWikiArray[x+7]

		// The Wikidata Descriptions field can not have more that 250 characters
		if len(occupations_library) >= 250 {
			occupations_library = occupations_library[:247] + "..."
		}

		replacer_occupations_library := strings.NewReplacer("\\", "", "\"[", "[", "]\"", "]", "\"", "'")
		occupations_library = replacer_occupations_library.Replace(occupations_library)
		nationality = functions.CheckCountryName(nationality)

		// Use the data from field 200 as label
		if field == "200" {

			entity.A_Label.A_Pt =
				A_Pts{
					"pt",
					name,
				}
			entity.A_Label.A_En =
				A_Ens{
					"en",
					name,
				}

			if len(occupations_library) > 0 {
				entity.A_Description.A_Pt =
					&A_Pts{
						"pt",
						occupations_library,
					}
			} else {
				entity.A_Description.A_Pt = nil
			}

			// export the property: "instance of" is entity: "human"
			entity.A_Claim.P5 = append(entity.A_Claim.P5, ReturnItemProperty("P5", 1, retrieved_date, id_library))

			// if library author record has any information about its birth date, export it
			if len(birth_date_library) > 0 {
				entity.A_Claim.P1 = append(entity.A_Claim.P1, ReturnTimeProperty("P1", birth_date_library, retrieved_date, id_library))
			}
			// if library author record has any information about its death date, export it
			if len(death_date_library) > 0 {
				entity.A_Claim.P2 = append(entity.A_Claim.P1, ReturnTimeProperty("P2", death_date_library, retrieved_date, id_library))
			}

			// NOT WORKING YET - should export to the author's page identifiers section
			// but it is exporting to the declarations section
			entity.A_Claim.P6 = append(entity.A_Claim.P6, ReturnStringProperty("P6", retrieved_date, id_library))

		}

		// if the author has variant names, export them
		if field == "400" {

			entity.A_Alias.A_Pt = append(entity.A_Alias.A_Pt, A_Pts{
				"pt",
				name,
			})

		}

	}

	// convert the struct to JSON
	t, _ := json.Marshal(entity)
	exportToWikidata := (string(t))

	// convert the result from []byte to a string
	exportToWikidata = replacerWiki.Replace(exportToWikidata)

	// Post the data to Wikidata
	resp, err := client.PostForm("http://127.0.0.1:8181/api.php?", url.Values{
		"action": {"wbeditentity"},
		"new":    {"item"},
		"token":  {tokenCsfr},
		"bot":    {"1"},
		"data":   {exportToWikidata},
		"format": {"json"},
	})
	if err != nil {
		log.Fatal(err, "erro post")
	}

	// read the response from the Wikidata API
	dataPost, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, "erro dataPost")
	}

	// get the posted Wikidata entity ID in order to write it in the database
	data := A_object{}
	errP := json.Unmarshal(dataPost, &data) // decode the JSON data
	if errP != nil {
		fmt.Println(errP)
	}

	id := data.A_Entity.A_Id

	fmt.Printf(`
	------------------------------
	  NEW WIKIDATA ID: %v					AUTHOR ID: %v
	  Author: %v
	------------------------------`, data.A_Entity.A_Id, id_library, data.A_Entity.A_Label.A_Pt.A_Value)

	// call the function that writes the new entity Wikidata ID in the authors table
	db.WriteAuthorId(id, id_library)
}
