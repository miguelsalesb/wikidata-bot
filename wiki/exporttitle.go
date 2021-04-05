package wiki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"wikidata/db"

	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/net/publicsuffix"
)

type T_object struct {
	T_Entity  T_Entity `json:"entity"`
	T_Success int      `json:"success"`
}

type T_Entity struct {
	T_Label T_Labels `json:"labels,omitempty"`
	// Description Descriptions `json:"descriptions,omitempty"`
	T_Claim     T_Claims `json:"claims,omitempty"`
	T_Id        string   `json:"id,omitempty"`
	T_Type      string   `json:"type,omitempty"`
	T_Lastrevid int      `json:"lastrevid,omitempty"`
}

type T_Labels struct {
	T_Pt *T_Pts `json:"pt,omitempty"`
	T_En *T_Ens `json:"en,omitempty"`
	T_Fr *T_Frs `json:"fr,omitempty"`
	T_Es *T_Ess `json:"es,omitempty"`
}

type T_Pts struct {
	T_Language string `json:"language,omitempty"`
	T_Value    string `json:"value,omitempty"`
}

type T_Ens struct {
	T_Language string `json:"language,omitempty"`
	T_Value    string `json:"value,omitempty"`
}

type T_Frs struct {
	T_Language string `json:"language,omitempty"`
	T_Value    string `json:"value,omitempty"`
}

type T_Ess struct {
	T_Language string `json:"language,omitempty"`
	T_Value    string `json:"value,omitempty"`
}

type T_Descriptions struct {
	T_Pt *T_Pts `json:"pt,omitempty"`
	T_En *T_Ens `json:"en,omitempty"`
	T_Fr *T_Frs `json:"fr,omitempty"`
	T_Es *T_Ess `json:"es,omitempty"`
}

type T_Claims struct {
	T_P8  []T_PItem  `json:"P8,omitempty"`  // Item type
	T_P10 []T_PItem  `json:"P10,omitempty"` // Item Type
	T_P9  []*T_PTime `json:"P9,omitempty"`  // Time type
}

type T_PItem struct {
	T_Mainsnak   T_MainsnaksItem `json:"mainsnak,omitempty"`
	T_Type       string          `json:"type,omitempty"`
	T_Rank       string          `json:"rank,omitempty"`
	T_References []T_References  `json:"references,omitempty"`
}

type T_PTime struct {
	T_Mainsnak   T_MainsnaksTime `json:"mainsnak,omitempty"`
	T_Type       string          `json:"type,omitempty"`
	T_Rank       string          `json:"rank,omitempty"`
	T_References []T_References  `json:"references,omitempty"`
}

type T_MainsnaksItem struct {
	T_Snaktype  string          `json:"snaktype,omitempty"`
	T_Property  string          `json:"property,omitempty"`
	T_Datatype  string          `json:"datatype,omitempty"`
	T_Datavalue T_DatavalueItem `json:"datavalue,omitempty"`
}

type T_MainsnaksTime struct {
	T_Snaktype  string          `json:"snaktype,omitempty"`
	T_Property  string          `json:"property,omitempty"`
	T_Datatype  string          `json:"datatype,omitempty"`
	T_Datavalue T_DatavalueTime `json:"datavalue,omitempty"`
}

type T_DatavalueItem struct {
	T_Value T_ValueItem `json:"value,omitempty"`
	T_Type  string      `json:"type,omitempty"`
}

type T_DatavalueTime struct {
	T_Value T_ValueTime `json:"value,omitempty"`
	T_Type  string      `json:"type,omitempty"`
}

type T_ValueItem struct {
	T_EntityType string `json:"entity-type,omitempty"`
	T_ID         string `json:"id,omitempty"`
	T_NumericID  int    `json:"numeric-id,omitempty"`
}

type T_ValueTime struct {
	T_Time          string `json:"time,omitempty"`
	T_Timezone      int    `json:"timezone"`
	T_Before        int    `json:"before"`
	T_After         int    `json:"after"`
	T_Precision     int    `json:"precision,omitempty"`
	T_CalendarModel string `json:"calendarmodel,omitempty"`
}

type T_References struct {
	T_Snaks      T_Snaks  `json:"snaks,omitempty"`
	T_SnaksOrder []string `json:"snaks-order,omitempty"`
}

type T_Snaks struct {
	T_P3 []T_PRefItem   `json:"P3,omitempty"` // Item type
	T_P4 []T_PRefString `json:"P8,omitempty"` // String type
	T_P7 []T_PRefTime   `json:"P7,omitempty"` // Time type
}

type T_PRefString struct {
	T_Snaktype  string               `json:"snaktype,omitempty"`
	T_Property  string               `json:"property,omitempty"`
	T_Datatype  string               `json:"datatype,omitempty"`
	T_Datavalue T_DatavalueRefString `json:"datavalue,omitempty"`
}

type T_PRefItem struct {
	T_Snaktype  string             `json:"snaktype,omitempty"`
	T_Property  string             `json:"property,omitempty"`
	T_Datatype  string             `json:"datatype,omitempty"`
	T_Datavalue T_DatavalueRefItem `json:"datavalue,omitempty"`
}

type T_PRefTime struct {
	T_Snaktype  string             `json:"snaktype,omitempty"`
	T_Property  string             `json:"property,omitempty"`
	T_Datatype  string             `json:"datatype,omitempty"`
	T_Datavalue T_DatavalueRefTime `json:"datavalue,omitempty"`
}

type T_DatavalueRefItem struct {
	T_Value T_ValueRefItem `json:"value,omitempty"`
	T_Type  string         `json:"type,omitempty"`
}

type T_DatavalueRefString struct {
	T_Value string `json:"value,omitempty"`
	T_Type  string `json:"type,omitempty"`
}

type T_ValueRefItem struct {
	T_EntityType string `json:"entity-type,omitempty"`
	T_Item       string `json:"item,omitempty"`
	T_NumericId  int    `json:"numeric-id,omitempty"`
}

type T_DatavalueRefTime struct {
	T_Value T_ValueRefTime `json:"value,omitempty"`
	T_Type  string         `json:"type,omitempty"`
}

type T_ValueRefTime struct {
	T_Time          string `json:"time,omitempty"`
	T_Timezone      int    `json:"timezone"`
	T_Before        int    `json:"before"`
	T_After         int    `json:"after"`
	T_Precision     int    `json:"precision,omitempty"`
	T_CalendarModel string `json:"calendarmodel,omitempty"`
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

func ExportTitle(idLibrary, originalLanguageOfWork, title, originalTitle, authorWiki, pubDate, retrieved_date string) {

	var tokenCsfr string
	var entity T_Entity

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar}

	// token needed to write the data in Wikidata
	tokenCsfr = ConnectToWikidata(client)

	replacerTitle := strings.NewReplacer("\\", "", "\"[", "[", "]\"", "]", "\"{", "{", "}\"", "}", "<", "", ">", "", "'", "\\'", "\"", "'")
	title = replacerTitle.Replace(title)

	// origLanguage = functions.CheckLanguagesCodes(originalLanguageOfWork)

	// if the original book is in english, export the portuguese and english titles
	if len(originalLanguageOfWork) > 0 && originalLanguageOfWork == "eng" && len(originalTitle) > 0 {

		entity.T_Label.T_Pt = &T_Pts{
			"pt",
			title,
		}

		entity.T_Label.T_En = &T_Ens{
			"en",
			originalTitle,
		}

		entity.T_Label.T_Fr = nil

		entity.T_Label.T_Es = nil

		// else if the original book is in french, export the portuguese and french titles
	} else if len(originalLanguageOfWork) > 0 && originalLanguageOfWork == "fre" && len(originalTitle) > 0 {

		entity.T_Label.T_Pt = &T_Pts{
			"pt",
			title,
		}

		entity.T_Label.T_Fr = &T_Frs{
			"fr",
			originalTitle,
		}

		entity.T_Label.T_En = nil

		entity.T_Label.T_Es = nil

		// if the original book is in spanish, export the portuguese and spanish titles
	} else if len(originalLanguageOfWork) > 0 && originalLanguageOfWork == "esp" && len(originalTitle) > 0 {

		entity.T_Label.T_Pt = &T_Pts{
			"pt",
			title,
		}

		entity.T_Label.T_Es = &T_Ess{
			"es",
			originalTitle,
		}

		entity.T_Label.T_En = nil

		entity.T_Label.T_Fr = nil

		// if the book is only in portuguese, write the portuguese title
	} else {

		entity.T_Label.T_Pt = &T_Pts{
			"pt",
			title,
		}

		entity.T_Label.T_Es = nil

		entity.T_Label.T_En = nil

		entity.T_Label.T_Fr = nil
	}

	// instance of notable work
	entity.T_Claim.T_P8 = append(entity.T_Claim.T_P8,
		T_PItem{
			T_MainsnaksItem{
				"value",
				"P8",
				"wikibase-item",
				T_DatavalueItem{
					T_ValueItem{
						"item",
						"Q3",
						3,
					},
					"wikibase-entityid",
				},
			},
			"statement",
			"normal",
			[]T_References{
				T_References{
					T_Snaks{
						[]T_PRefItem{
							T_PRefItem{
								"value",
								"P3",
								"wikibase-item",
								T_DatavalueRefItem{
									T_ValueRefItem{
										"item",
										"Q2",
										2,
									},
									"wikibase-entityid",
								},
							},
						},
						[]T_PRefString{
							T_PRefString{
								"value",
								"P4",
								"url",
								T_DatavalueRefString{
									"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + idLibrary,
									"string",
								},
							},
						},
						[]T_PRefTime{
							T_PRefTime{
								"value",
								"P7",
								"time",
								T_DatavalueRefTime{
									T_ValueRefTime{
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

	// country of origin - Portugal - it only exports portuguese works
	entity.T_Claim.T_P10 = append(entity.T_Claim.T_P10,
		T_PItem{
			T_MainsnaksItem{
				"value",
				"P10",
				"wikibase-item",
				T_DatavalueItem{
					T_ValueItem{
						"item",
						"Q4",
						4,
					},
					"wikibase-entityid",
				},
			},
			"statement",
			"normal",
			[]T_References{
				T_References{
					T_Snaks{
						[]T_PRefItem{
							T_PRefItem{
								"value",
								"P3",
								"wikibase-item",
								T_DatavalueRefItem{
									T_ValueRefItem{
										"item",
										"Q2",
										2,
									},
									"wikibase-entityid",
								},
							},
						},
						[]T_PRefString{
							T_PRefString{
								"value",
								"P4",
								"url",
								T_DatavalueRefString{
									"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + idLibrary,
									"string",
								},
							},
						},
						[]T_PRefTime{
							T_PRefTime{
								"value",
								"P7",
								"time",
								T_DatavalueRefTime{
									T_ValueRefTime{
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

	// if it has a publication data, export it to Wikidata
	if len(pubDate) > 0 {
		entity.T_Claim.T_P9 = append(entity.T_Claim.T_P9,
			&T_PTime{
				T_MainsnaksTime{
					"value",
					"P9",
					"time",
					T_DatavalueTime{
						T_ValueTime{
							"+" + pubDate + "-00-00T00:00:00Z",
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
				[]T_References{
					T_References{
						T_Snaks{
							[]T_PRefItem{
								T_PRefItem{
									"value",
									"P3",
									"wikibase-item",
									T_DatavalueRefItem{
										T_ValueRefItem{
											"item",
											"Q2",
											2,
										},
										"wikibase-entityid",
									},
								},
							},
							[]T_PRefString{
								T_PRefString{
									"value",
									"P4",
									"url",
									T_DatavalueRefString{
										"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + idLibrary,
										"string",
									},
								},
							},
							[]T_PRefTime{
								T_PRefTime{
									"value",
									"P7",
									"time",
									T_DatavalueRefTime{
										T_ValueRefTime{
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
	} else {
		entity.T_Claim.T_P9 = nil
	}

	// convert the struct to JSON
	t, _ := json.Marshal(entity)

	// convert the result from []byte to a string
	exportTitleToWikidata := (string(t))

	// replace some characters
	replacerExport := strings.NewReplacer("\\", "", "\"[", "[", "]\"", "]", "\"{", "{", "}\"", "}")
	exportTitleToWikidata = replacerExport.Replace(exportTitleToWikidata)

	// Post the data to Wikidata
	resp, err := client.PostForm("http://127.0.0.1:8181/api.php?", url.Values{
		"action": {"wbeditentity"},
		"new":    {"item"},
		"token":  {tokenCsfr},
		"bot":    {"1"},
		"data":   {exportTitleToWikidata},
		"format": {"json"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// read the response from the Wikidata API
	dataPost, errPost := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(errPost)
	}

	// get the posted Wikidata entity ID in order to write it in the database
	data := T_object{}
	errData := json.Unmarshal(dataPost, &data) // decode the JSON data
	if errData != nil {
		fmt.Println(errData)
	}

	idWiki := data.T_Entity.T_Id

	fmt.Printf(`
	------------------------------
	  NEW WIKIDATA ID: %v					TITLE ID: %v
	  Title: %v
	------------------------------`, data.T_Entity.T_Id, idLibrary, data.T_Entity.T_Label.T_Pt.T_Value)

	// call the function that writes the new entity Wikidata ID in the titles table
	db.WriteTitleId(idWiki, idLibrary)
}
