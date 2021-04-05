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
	A_P5 []A_PItem   `json:"P5,omitempty"` // Entity type
	A_P8 []A_PItem   `json:"P8,omitempty"` // Entity type
	A_P1 []A_PTime   `json:"P1,omitempty"` // Time type
	A_P2 []A_PTime   `json:"P2,omitempty"` // Time type
	A_P6 []A_PString `json:"P6,omitempty"` // Entity type
}

type A_PTime struct {
	A_Mainsnak   A_MainsnaksTime `json:"mainsnak,omitempty"`
	A_Type       string          `json:"type,omitempty"`
	A_Rank       string          `json:"rank,omitempty"`
	A_References []A_References  `json:"references,omitempty"`
}

type A_PItem struct {
	A_Mainsnak   A_MainsnaksItem `json:"mainsnak,omitempty"`
	A_Type       string          `json:"type,omitempty"`
	A_Rank       string          `json:"rank,omitempty"`
	A_References []A_References  `json:"references,omitempty"`
}

type A_PString struct {
	A_Mainsnak   A_MainsnaksString `json:"mainsnak,omitempty"`
	A_Type       string            `json:"type,omitempty"`
	A_Rank       string            `json:"rank,omitempty"`
	A_References []A_References    `json:"references,omitempty"`
}

type A_MainsnaksTime struct {
	A_Snaktype  string          `json:"snaktype,omitempty"`
	A_Property  string          `json:"property,omitempty"`
	A_Datatype  string          `json:"datatype,omitempty"`
	A_Datavalue A_DatavalueTime `json:"datavalue,omitempty"`
}

type A_MainsnaksItem struct {
	A_Snaktype  string          `json:"snaktype,omitempty"`
	A_Property  string          `json:"property,omitempty"`
	A_Datatype  string          `json:"datatype,omitempty"`
	A_Datavalue A_DatavalueItem `json:"datavalue,omitempty"`
}

type A_MainsnaksString struct {
	A_Snaktype  string            `json:"snaktype,omitempty"`
	A_Property  string            `json:"property,omitempty"`
	A_Datatype  string            `json:"datatype,omitempty"`
	A_Datavalue A_DatavalueString `json:"datavalue,omitempty"`
}

type A_DatavalueTime struct {
	A_Value A_ValueTime `json:"value,omitempty"`
	A_Type  string      `json:"type,omitempty"`
}

type A_DatavalueItem struct {
	A_Value A_ValueItem `json:"value,omitempty"`
	A_Type  string      `json:"type,omitempty"`
}

type A_DatavalueString struct {
	A_Value string `json:"value,omitempty"`
	A_Type  string `json:"type,omitempty"`
}

type A_ValueTime struct {
	A_Time          string `json:"time,omitempty"`
	A_Timezone      int    `json:"timezone"`
	A_Before        int    `json:"before"`
	A_After         int    `json:"after"`
	A_Precision     int    `json:"precision,omitempty"`
	A_CalendarModel string `json:"calendarmodel,omitempty"`
}

type A_ValueItem struct {
	A_EntityType string `json:"entity-type,omitempty"`
	A_ID         string `json:"id,omitempty"`
	A_NumericID  int    `json:"numeric-id,omitempty"`
}

type A_References struct {
	A_Snaks      A_Snaks  `json:"snaks,omitempty"`
	A_SnaksOrder []string `json:"snaks-order,omitempty"`
}

type A_Snaks struct {
	A_P3 []A_PRefItem   `json:"P3,omitempty"` // Item type
	A_P4 []A_PRefString `json:"P8,omitempty"` // String type
	A_P7 []A_PRefTime   `json:"P7,omitempty"` // Time type
}

type A_PRefString struct {
	A_Snaktype  string               `json:"snaktype,omitempty"`
	A_Property  string               `json:"property,omitempty"`
	A_Datatype  string               `json:"datatype,omitempty"`
	A_Datavalue A_DatavalueRefString `json:"datavalue,omitempty"`
}

type A_PRefItem struct {
	A_Snaktype  string             `json:"snaktype,omitempty"`
	A_Property  string             `json:"property,omitempty"`
	A_Datatype  string             `json:"datatype,omitempty"`
	A_Datavalue A_DatavalueRefItem `json:"datavalue,omitempty"`
}

type A_PRefTime struct {
	A_Snaktype  string             `json:"snaktype,omitempty"`
	A_Property  string             `json:"property,omitempty"`
	A_Datatype  string             `json:"datatype,omitempty"`
	A_Datavalue A_DatavalueRefTime `json:"datavalue,omitempty"`
}

type A_DatavalueRefItem struct {
	A_Value A_ValueRefItem `json:"value,omitempty"`
	A_Type  string         `json:"type,omitempty"`
}

type A_DatavalueRefString struct {
	A_Value string `json:"value,omitempty"`
	A_Type  string `json:"type,omitempty"`
}

type A_ValueRefItem struct {
	A_EntityType string `json:"entity-type,omitempty"`
	A_Item       string `json:"item,omitempty"`
	A_NumericId  int    `json:"numeric-id,omitempty"`
}

type A_DatavalueRefTime struct {
	A_Value A_ValueRefTime `json:"value,omitempty"`
	A_Type  string         `json:"type,omitempty"`
}

type A_ValueRefTime struct {
	A_Time          string `json:"time,omitempty"`
	A_Timezone      int    `json:"timezone"`
	A_Before        int    `json:"before"`
	A_After         int    `json:"after"`
	A_Precision     int    `json:"precision,omitempty"`
	A_CalendarModel string `json:"calendarmodel,omitempty"`
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
		replacer                                                                                                          = strings.NewReplacer("\\", "", "\"[", "[", "]\"", "]")
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
		// Q1 ser humano (Q1) - human (Q5)
		// Q2 BNP - National Library of Portugal (Q245966)
		// Q3 obra escrita - written work (Q47461344)
		// Q4 Portugal - Portugal (Q45)

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
			entity.A_Claim.A_P5 = append(entity.A_Claim.A_P5,
				A_PItem{
					A_MainsnaksItem{
						"value",
						"P5",
						"wikibase-item",
						A_DatavalueItem{
							A_ValueItem{
								"item",
								"Q1",
								1,
							},
							"wikibase-entityid",
						},
					},
					"statement",
					"normal",
					[]A_References{
						A_References{
							A_Snaks{
								[]A_PRefItem{
									A_PRefItem{
										"value",
										"P3",
										"wikibase-item",
										A_DatavalueRefItem{
											A_ValueRefItem{
												"item",
												"Q2",
												2,
											},
											"wikibase-entityid",
										},
									},
								},
								[]A_PRefString{
									A_PRefString{
										"value",
										"P4",
										"url",
										A_DatavalueRefString{
											"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + id_library,
											"string",
										},
									},
								},
								[]A_PRefTime{
									A_PRefTime{
										"value",
										"P7",
										"time",
										A_DatavalueRefTime{
											A_ValueRefTime{
												retrieved_date,
												0,
												0,
												0,
												11,
												"http://www.wikidata.org/entity/Q1985727",
											},
											"time",
										},
									},
								},
							},
							[]string{

								"P3", "P4", "P7",
							},
						},
					},
				},
			)

			// if library author record has any information about its birth date, export it
			if len(birth_date_library) > 0 {
				entity.A_Claim.A_P1 = append(entity.A_Claim.A_P1,
					A_PTime{
						A_MainsnaksTime{
							"value",
							"P1",
							"time",
							A_DatavalueTime{
								A_ValueTime{
									"+" + birth_date_library + "-00-00T00:00:00Z",
									0,
									0,
									0,
									9,
									"http://www.wikidata.org/entity/Q1985727",
								},
								"time",
							},
						},
						"statement",
						"normal",
						[]A_References{
							A_References{
								A_Snaks{
									[]A_PRefItem{
										A_PRefItem{
											"value",
											"P3",
											"wikibase-item",
											A_DatavalueRefItem{
												A_ValueRefItem{
													"item",
													"Q2",
													2,
												},
												"wikibase-entityid",
											},
										},
									},
									[]A_PRefString{
										A_PRefString{
											"value",
											"P4",
											"url",
											A_DatavalueRefString{
												"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + id_library,
												"string",
											},
										},
									},
									[]A_PRefTime{
										A_PRefTime{
											"value",
											"P7",
											"time",
											A_DatavalueRefTime{
												A_ValueRefTime{
													retrieved_date,
													0,
													0,
													0,
													11,
													"http://www.wikidata.org/entity/Q1985727",
												},
												"time",
											},
										},
									},
								},
								[]string{

									"P3", "P4", "P7",
								},
							},
						},
					},
				)
			}

			// if library author record has any information about its death date, export it
			if len(death_date_library) > 0 {
				entity.A_Claim.A_P2 = append(entity.A_Claim.A_P2,
					A_PTime{
						A_MainsnaksTime{
							"value",
							"P2",
							"time",
							A_DatavalueTime{
								A_ValueTime{
									"+" + death_date_library + "-00-00T00:00:00Z",
									0,
									0,
									0,
									9,
									"http://www.wikidata.org/entity/Q1985727",
								},
								"time",
							},
						},
						"statement",
						"normal",
						[]A_References{
							A_References{
								A_Snaks{
									[]A_PRefItem{
										A_PRefItem{
											"value",
											"P3",
											"wikibase-item",
											A_DatavalueRefItem{
												A_ValueRefItem{
													"item",
													"Q2",
													2,
												},
												"wikibase-entityid",
											},
										},
									},
									[]A_PRefString{
										A_PRefString{
											"value",
											"P4",
											"url",
											A_DatavalueRefString{
												"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + id_library,
												"string",
											},
										},
									},
									[]A_PRefTime{
										A_PRefTime{
											"value",
											"P7",
											"time",
											A_DatavalueRefTime{
												A_ValueRefTime{
													retrieved_date,
													0,
													0,
													0,
													11,
													"http://www.wikidata.org/entity/Q1985727",
												},
												"time",
											},
										},
									},
								},
								[]string{

									"P3", "P4", "P7",
								},
							},
						},
					},
				)
			}

			// NOT WORKING YET - should export to the author's page identifiers section
			// but it is exporting to the declarations section
			entity.A_Claim.A_P6 = append(entity.A_Claim.A_P6,
				A_PString{
					A_MainsnaksString{
						"value",
						"P6",
						"external-id",
						A_DatavalueString{
							"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + id_library,
							"string",
						},
					},

					"statement",
					"normal",
					[]A_References{
						A_References{
							A_Snaks{
								[]A_PRefItem{
									A_PRefItem{
										"value",
										"P3",
										"wikibase-item",
										A_DatavalueRefItem{
											A_ValueRefItem{
												"item",
												"Q2",
												2,
											},
											"wikibase-entityid",
										},
									},
								},
								[]A_PRefString{
									A_PRefString{
										"value",
										"P4",
										"url",
										A_DatavalueRefString{
											"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + id_library,
											"string",
										},
									},
								},
								[]A_PRefTime{
									A_PRefTime{
										"value",
										"P7",
										"time",
										A_DatavalueRefTime{
											A_ValueRefTime{
												retrieved_date,
												0,
												0,
												0,
												11,
												"http://www.wikidata.org/entity/Q1985727",
											},
											"time",
										},
									},
								},
							},
							[]string{

								"P3", "P4", "P7",
							},
						},
					},
				},
			)

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
	exportToWikidata = replacer.Replace(exportToWikidata)

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
