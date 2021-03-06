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

var Replacer = strings.NewReplacer(" ", "%20", "??", "%C3%80", "??", "%C3%81", "??", "%C3%82", "??", "%C3%83", "??", "%C3%84", "??", "%C3%87", "??", "%C3%88",
	"??", "%C3%89", "??", "%C3%8A", "??", "%C3%8B", "??", "%C3%8C", "??", "%C3%8D", "??", "%C3%8E", "??", "%C3%8F", "??", "%C3%92", "??", "%C3%93", "??", "%C3%94",
	"??", "%C3%95", "??", "%C3%96", "??", "%C3%99", "??", "%C3%9A", "??", "%C3%9B", "??", "%C3%9D", "??", "%C3%A0", "??", "%C3%A1", "??", "%C3%A2", "??", "%C3%A3",
	"??", "%C3%A4", "??", "%C3%A7", "??", "%C3%A8", "??", "%C3%A9", "??", "%C3%AA", "??", "%C3%AB", "??", "%C3%AC", "??", "%C3%AD", "??", "%C3%AE", "??", "C3%AF",
	"??", "%C3%B1", "??", "%C3%B2", "??", "%C3%B3", "??", "%C3%B4", "??", "%C3%B5", "??", "%C3%B6", "??", "%C3%B9", "??", "%C3%BA", "??", "%C3%BB", "??", "%C3%BC",
	"??", "%C3%BD", "\"", "'", "??", "%C2%BA", "??", "%C2%AA", "&", "%26", ",", "%2C", "!", "%21", "#", "%23", "$", "%24", "%", "%25", "'", "%27", "(", "%28",
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
