package wiki

import (
	"strconv"
	"strings"
)

type PItem struct {
	Mainsnak   MainsnaksItem `json:"mainsnak,omitempty"`
	Type       string        `json:"type,omitempty"`
	Rank       string        `json:"rank,omitempty"`
	References []References  `json:"references,omitempty"`
}

type PTime struct {
	Mainsnak   MainsnaksTime `json:"mainsnak,omitempty"`
	Type       string        `json:"type,omitempty"`
	Rank       string        `json:"rank,omitempty"`
	References []References  `json:"references,omitempty"`
}

type PString struct {
	Mainsnak   MainsnaksString `json:"mainsnak,omitempty"`
	Type       string          `json:"type,omitempty"`
	Rank       string          `json:"rank,omitempty"`
	References []References    `json:"references,omitempty"`
}

type MainsnaksItem struct {
	Snaktype  string        `json:"snaktype,omitempty"`
	Property  string        `json:"property,omitempty"`
	Datatype  string        `json:"datatype,omitempty"`
	Datavalue DatavalueItem `json:"datavalue,omitempty"`
}

type MainsnaksTime struct {
	Snaktype  string        `json:"snaktype,omitempty"`
	Property  string        `json:"property,omitempty"`
	Datatype  string        `json:"datatype,omitempty"`
	Datavalue DatavalueTime `json:"datavalue,omitempty"`
}

type MainsnaksString struct {
	Snaktype  string          `json:"snaktype,omitempty"`
	Property  string          `json:"property,omitempty"`
	Datatype  string          `json:"datatype,omitempty"`
	Datavalue DatavalueString `json:"datavalue,omitempty"`
}

type DatavalueString struct {
	Value string `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
}

type DatavalueItem struct {
	Value ValueItem `json:"value,omitempty"`
	Type  string    `json:"type,omitempty"`
}

type DatavalueTime struct {
	Value ValueTime `json:"value,omitempty"`
	Type  string    `json:"type,omitempty"`
}

type ValueItem struct {
	EntityType string `json:"entity-type,omitempty"`
	ID         string `json:"id,omitempty"`
	NumericID  int    `json:"numeric-id,omitempty"`
}

type ValueTime struct {
	Time          string `json:"time,omitempty"`
	Timezone      int    `json:"timezone"`
	Before        int    `json:"before"`
	After         int    `json:"after"`
	Precision     int    `json:"precision,omitempty"`
	CalendarModel string `json:"calendarmodel,omitempty"`
}

type References struct {
	Snaks      Snaks    `json:"snaks,omitempty"`
	SnaksOrder []string `json:"snaks-order,omitempty"`
}

type Snaks struct {
	P3 []PRefItem   `json:"P3,omitempty"` // Item type
	P4 []PRefString `json:"P8,omitempty"` // String type
	P7 []PRefTime   `json:"P7,omitempty"` // Time type
}

type PRefString struct {
	Snaktype  string             `json:"snaktype,omitempty"`
	Property  string             `json:"property,omitempty"`
	Datatype  string             `json:"datatype,omitempty"`
	Datavalue DatavalueRefString `json:"datavalue,omitempty"`
}

type PRefItem struct {
	Snaktype  string           `json:"snaktype,omitempty"`
	Property  string           `json:"property,omitempty"`
	Datatype  string           `json:"datatype,omitempty"`
	Datavalue DatavalueRefItem `json:"datavalue,omitempty"`
}

type PRefTime struct {
	Snaktype  string           `json:"snaktype,omitempty"`
	Property  string           `json:"property,omitempty"`
	Datatype  string           `json:"datatype,omitempty"`
	Datavalue DatavalueRefTime `json:"datavalue,omitempty"`
}

type DatavalueRefItem struct {
	Value ValueRefItem `json:"value,omitempty"`
	Type  string       `json:"type,omitempty"`
}

type DatavalueRefString struct {
	Value string `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
}

type ValueRefItem struct {
	EntityType string `json:"entity-type,omitempty"`
	Item       string `json:"item,omitempty"`
	NumericId  int    `json:"numeric-id,omitempty"`
}

type DatavalueRefTime struct {
	Value ValueRefTime `json:"value,omitempty"`
	Type  string       `json:"type,omitempty"`
}

type ValueRefTime struct {
	Time          string `json:"time,omitempty"`
	Timezone      int    `json:"timezone"`
	Before        int    `json:"before"`
	After         int    `json:"after"`
	Precision     int    `json:"precision,omitempty"`
	CalendarModel string `json:"calendarmodel,omitempty"`
}

var Replacer = strings.NewReplacer(" ", "%20", "À", "%C3%80", "Á", "%C3%81", "Â", "%C3%82", "Ã", "%C3%83", "Ä", "%C3%84", "Ç", "%C3%87", "È", "%C3%88",
	"É", "%C3%89", "Ê", "%C3%8A", "Ë", "%C3%8B", "Ì", "%C3%8C", "Í", "%C3%8D", "Î", "%C3%8E", "Ï", "%C3%8F", "Ò", "%C3%92", "Ó", "%C3%93", "Ô", "%C3%94",
	"Õ", "%C3%95", "Ö", "%C3%96", "Ù", "%C3%99", "Ó", "%C3%9A", "Û", "%C3%9B", "Ý", "%C3%9D", "à", "%C3%A0", "á", "%C3%A1", "â", "%C3%A2", "ã", "%C3%A3",
	"ä", "%C3%A4", "ç", "%C3%A7", "è", "%C3%A8", "é", "%C3%A9", "ê", "%C3%AA", "ë", "%C3%AB", "ì", "%C3%AC", "í", "%C3%AD", "î", "%C3%AE", "ï", "C3%AF",
	"ñ", "%C3%B1", "ò", "%C3%B2", "ó", "%C3%B3", "ô", "%C3%B4", "õ", "%C3%B5", "ö", "%C3%B6", "ù", "%C3%B9", "ú", "%C3%BA", "û", "%C3%BB", "ü", "%C3%BC",
	"ý", "%C3%BD", "\"", "'", "º", "%C2%BA", "ª", "%C2%AA", "&", "%26", ",", "%2C", "!", "%21", "#", "%23", "$", "%24", "%", "%25", "'", "%27", "(", "%28",
	")", "%29", "-", "%2D", "[", "%5B", "]", "%5D", "^", "%5E", "_", "%5F", "_", "%60", "{", "%7B", "{", "%7C", "}", "%7D")

func ReturnItemProperty(prop string, QItem int, retrieved_date string, id_library string) PItem {
	data :=
		PItem{
			MainsnaksItem{
				"value",
				prop,
				"wikibase-item",
				DatavalueItem{
					ValueItem{
						"item",
						"Q" + strconv.Itoa(QItem),
						QItem,
					},
					"wikibase-entityid",
				},
			},
			"statement",
			"normal",
			[]References{
				References{
					Snaks{
						[]PRefItem{
							PRefItem{
								"value",
								"P3",
								"wikibase-item",
								DatavalueRefItem{
									ValueRefItem{
										"item",
										"Q2",
										2,
									},
									"wikibase-entityid",
								},
							},
						},
						[]PRefString{
							PRefString{
								"value",
								"P4",
								"url",
								DatavalueRefString{
									"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + id_library,
									"string",
								},
							},
						},
						[]PRefTime{
							PRefTime{
								"value",
								"P7",
								"time",
								DatavalueRefTime{
									ValueRefTime{
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
		}
	return data
}

func ReturnTimeProperty(prop string, date string, retrieved_date string, id_library string) *PTime {

	data :=
		&PTime{
			MainsnaksTime{
				"value",
				prop,
				"time",
				DatavalueTime{
					ValueTime{
						"+" + date + "-00-00T00:00:00Z",
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
			[]References{
				References{
					Snaks{
						[]PRefItem{
							PRefItem{
								"value",
								"P3",
								"wikibase-item",
								DatavalueRefItem{
									ValueRefItem{
										"item",
										"Q2",
										2,
									},
									"wikibase-entityid",
								},
							},
						},
						[]PRefString{
							PRefString{
								"value",
								"P4",
								"url",
								DatavalueRefString{
									"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + id_library,
									"string",
								},
							},
						},
						[]PRefTime{
							PRefTime{
								"value",
								"P7",
								"time",
								DatavalueRefTime{
									ValueRefTime{
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
		}
	return data
}

func ReturnStringProperty(prop string, retrieved_date string, id_library string) PString {

	data := PString{
		MainsnaksString{
			"value",
			prop,
			"external-id",
			DatavalueString{
				"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + id_library,
				"string",
			},
		},

		"statement",
		"normal",
		[]References{
			References{
				Snaks{
					[]PRefItem{
						PRefItem{
							"value",
							"P3",
							"wikibase-item",
							DatavalueRefItem{
								ValueRefItem{
									"item",
									"Q2",
									2,
								},
								"wikibase-entityid",
							},
						},
					},
					[]PRefString{
						PRefString{
							"value",
							"P4",
							"url",
							DatavalueRefString{
								"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + id_library,
								"string",
							},
						},
					},
					[]PRefTime{
						PRefTime{
							"value",
							"P7",
							"time",
							DatavalueRefTime{
								ValueRefTime{
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
	}
	return data
}

func ReturnIdentifier(prop string, retrieved_date string, id_library string) PString {

	data := PString{
		MainsnaksString{
			"value",
			prop,
			"external-id",
			DatavalueString{
				id_library,
				"string",
			},
		},
		"statement",
		"normal",
		[]References{
			References{
				Snaks{
					[]PRefItem{
						PRefItem{
							"value",
							"P3",
							"wikibase-item",
							DatavalueRefItem{
								ValueRefItem{
									"item",
									"Q2",
									2,
								},
								"wikibase-entityid",
							},
						},
					},
					[]PRefString{
						PRefString{
							"value",
							"P4",
							"url",
							DatavalueRefString{
								"http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=" + id_library,
								"string",
							},
						},
					},
					[]PRefTime{
						PRefTime{
							"value",
							"P7",
							"time",
							DatavalueRefTime{
								ValueRefTime{
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
	}
	return data
}
