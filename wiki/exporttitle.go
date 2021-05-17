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
	P8  []PItem   `json:"P8,omitempty"`  // Item type
	P10 []PItem   `json:"P10,omitempty"` // Item Type
	P9  []*PTime  `json:"P9,omitempty"`  // Time type
	P6  []PString `json:"P6,omitempty"`  // Entity type
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
	entity.T_Claim.P8 = append(entity.T_Claim.P8, ReturnItemProperty("P8", 3, retrieved_date, idLibrary))

	// country of origin - Portugal - it only exports portuguese works
	entity.T_Claim.P10 = append(entity.T_Claim.P10, ReturnItemProperty("P10", 4, retrieved_date, idLibrary))

	// if it has a publication data, export it to Wikidata

	if len(pubDate) > 0 {
		entity.T_Claim.P9 = append(entity.T_Claim.P9, ReturnTimeProperty("P9", pubDate, retrieved_date, idLibrary))
	} else {
		entity.T_Claim.P9 = nil
	}

	// NOT WORKING YET - should export to the author's page identifiers section
	// but it is exporting to the declarations section
	entity.T_Claim.P6 = append(entity.T_Claim.P6, ReturnIdentifier("P6", retrieved_date, idLibrary))

	// convert the struct to JSON
	t, _ := json.Marshal(entity)

	// convert the result from []byte to a string
	exportTitleToWikidata := (string(t))

	// replace some characters
	replacerExport := strings.NewReplacer("\\", "", "\"[", "[", "]\"", "]")

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
