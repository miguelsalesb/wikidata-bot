package titles

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"wikidata/db"
	"wikidata/functions"
	"wikidata/wiki"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

type Data struct {
	LibraryData       libraryData
	wkTitleID         WikiData
	wkOriginalTitleID WikiData
}

type libraryData struct {
	id                     string
	languageOfWork         string
	originalLanguageOfWork string
	title                  string
	titleLowercase         string
	authors                []string
	originalTitle          string
	originalTitleLowercase string
	pubDate                string
	bibnac                 string
	bnd                    string
}

type WikiData struct {
	id      string
	mattype string
}

// Structs to get the titles Wikidata ID
type object struct {
	head
	Results results
}

type head struct {
	Vars vars
}

type vars struct {
	Vars []string `json:"type"`
}

type results struct {
	Bindings []bindings `json:"bindings"`
}

type bindings struct {
	Item    item    `json:"item"`
	Mattype mattype `json:"mattype"`
	Author  author  `json:"occurences"`
}

type author struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type item struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type mattype struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func check(e error) {
	if e != nil {
		// panic(e)
		fmt.Println(e)
	}
}

var dbName = "titles"
var replacerTitle = strings.NewReplacer("<", "", ">", "", "'", "\\'", "«", "", "»", "", "º", "", "[", "", "]", "", "\"", "\\'")

func GetTitles(doneTitles chan bool, repTitlesFirst int, repTitlesLast int) {

	var wkOriginalTitleID WikiData
	const empty = ""

	for n := repTitlesFirst; n <= repTitlesLast; n++ {
		time.Sleep(500 * time.Millisecond)
		if n%500 == 0 {
			time.Sleep(120 * time.Second)
		}
		fmt.Println("\nTitles - ", n)
		url := fmt.Sprintf("%s%d", "http://urn.bn.pt/ncb/unimarc/marcxchange?id=", n)

		libData := getTitles(url)

		// In Wikidata search the title in lowercase
		titleLibrary := libData.titleLowercase
		language := functions.CheckLanguagesCodes(libData.languageOfWork)

		// In Wikidata search the original language title in lowercase
		originalTitleLibrary := libData.originalTitleLowercase
		originalLanguage := functions.CheckLanguagesCodes(libData.originalLanguageOfWork)

		authorsLibrary := libData.authors

		wkTitleID := GetWikiTitleID(titleLibrary, authorsLibrary, language)

		if originalLanguage != "" { // there main not be an original language if the item language io Portuguese
			wkOriginalTitleID = GetWikiTitleID(originalTitleLibrary, authorsLibrary, originalLanguage)
		} else {
			wkOriginalTitleID = WikiData{
				empty,
				empty,
			}
		}
		titles := Data{
			libData,
			wkTitleID,
			wkOriginalTitleID,
		}
		WriteTitles(titles, db.DBCon)
	}
	doneTitles <- true
}

func getTitles(url string) libraryData {

	var (
		idLibrary, titl, titlLowercase, idAuthor, name, surname, languageOfWork, originalLanguageOfWork, field,
		pubDate, origTi, originalTitle, originalTitleLowercase, bibnac, bnd string
		authors = make([]string, 0, 2)
	)
	const empty = ""

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the XML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	doc.Find("datafield").Each(func(i int, s *goquery.Selection) {
		tag, _ := s.Attr("tag")
		ind1, _ := s.Attr("ind1")
		ind2, _ := s.Attr("ind2")

		id := doc.Find("controlfield")
		if tag, ok := id.Attr("tag"); tag == "001" {
			if ok {
				idLibrary = id.First().Text()
			}
		}
		if tag == "101" {
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "a" {
					languageOfWork = e.Text()
				}
				if attr, _ := e.Attr("code"); attr == "c" {
					originalLanguageOfWork = e.Text()
				}
			})
		}
		if tag == "200" {
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "a" {
					ti := e.Text()
					// replacer := strings.NewReplacer("<", "", ">", "", "'", "\\'")
					titl = replacerTitle.Replace(ti)
					titlLowercase = strings.ToLower(titl) // library in lowercase and the sparql query will also search for the lowercased title
				}
			})
		}
		if tag == "304" {
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "a" {
					o := s.Text()
					origT := strings.Index(o, "orig.")
					origTi = strings.TrimLeft(o[origT+6:], ":")
					// replacer := strings.NewReplacer("<", "", ">", "", "'", "\\'", "«", "", "»", "", "º", "")
					originalTitle = replacerTitle.Replace(origTi)
					originalTitleLowercase = strings.TrimLeft(strings.ToLower(originalTitle), " ") // library in lowercase and the sparql query will also search for the lowercased title
				}
			})
		}
		if ind1 == "4" && ind2 == "0" && tag == "856" {
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "u" {
					bibnac = s.Text()
				}
			})
		}
		if tag == "900" {
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "a" {
					b := s.Text()
					if strings.Contains(b, "BIBNAC") {
						bibnac = idLibrary
					}
				}
			})
		}
	})

	doc.Find("datafield").Each(func(i int, s *goquery.Selection) {

		tag, _ := s.Attr("tag")
		// To get all authors, put tag[:1]
		if strings.Contains(tag[:3], "210") {
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "d" {
					r, _ := regexp.Compile("[0-9]{4}")
					e := s.Text()
					dts := r.FindString(e)
					pubDate = dts
				}
			})
		}
	})

	doc.Find("datafield").Each(func(i int, s *goquery.Selection) {

		tag, _ := s.Attr("tag")

		if tag[:2] == "70" { // To get only authors persons and not authors entities

			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "3" {
					id := e.Text()
					// replacer := strings.NewReplacer(",", "", "\\", "", "'", "\\'")
					idAuthor = replacerTitle.Replace(id)
				}
				if attr, _ := e.Attr("code"); attr == "a" {
					nm := e.Text()
					// replacer := strings.NewReplacer(",", "", "\\", "", "'", "\\'")
					name = replacerTitle.Replace(nm)
				}
				if attr, _ := e.Attr("code"); attr == "b" {
					sn := e.Text()
					// replacer := strings.NewReplacer(",", "", "\\", "", "'", "\\'")
					surname = replacerTitle.Replace(sn)
				}
			})
			field = tag
			if idAuthor != "" {
				authors = append(authors, idAuthor)
			} else {
				authors = append(authors, empty)
			}
			if name != "" {
				authors = append(authors, surname+" "+name)
			} else {
				authors = append(authors, empty)
			}
			if field != "" {
				authors = append(authors, field)
			} else {
				authors = append(authors, empty)
			}
		}
	})

	data := libraryData{
		idLibrary,
		languageOfWork,
		originalLanguageOfWork,
		titl,
		titlLowercase,
		authors,
		originalTitle,
		originalTitleLowercase,
		pubDate,
		bibnac,
		bnd,
	}
	// fmt.Println("\nDATA: ", data)
	return data
}

func WriteTitles(data Data, db *sql.DB) {

	// author for wikidata
	// const empty = ""
	var (
		authorWiki, idAuthor, author, field string
	)
	const empty = ""

	idLibrary := data.LibraryData.id
	languageOfWork := data.LibraryData.languageOfWork
	originalLanguageOfWork := data.LibraryData.originalLanguageOfWork
	title := data.LibraryData.title
	// fmt.Println("\nTitle: ", title)
	titleLowercase := data.LibraryData.titleLowercase
	titleIDWiki := data.wkTitleID.id
	titleMattypeWiki := data.wkTitleID.mattype
	originalTitleIDWiki := data.wkOriginalTitleID.id
	originalTitleMattypeWiki := data.wkOriginalTitleID.mattype
	pubDate := data.LibraryData.pubDate
	originalTitle := data.LibraryData.originalTitle
	originalTitleLowercase := data.LibraryData.originalTitleLowercase
	bibnac := data.LibraryData.bibnac
	bnd := data.LibraryData.bnd
	retrieved_date := time.Now().Format("+2006-01-02T00:00:00Z")

	for x := 0; x < len(data.LibraryData.authors); x += 3 {

		idAuthor = data.LibraryData.authors[x]
		author = data.LibraryData.authors[x+1]
		field = data.LibraryData.authors[x+2]

		stmt, err := db.Exec("INSERT INTO " + dbName + " (id_library, language_of_work, original_language_of_work, title_id_wiki, title_mattype_wiki, title, title_lowercase, original_title_id_wiki, original_title_mattype_wiki, original_title, original_title_lowercase, id_author, author, field, pub_date, bibnac, bnd, retrieved_date) VALUES (" + idLibrary + " , '" + languageOfWork + "', '" + originalLanguageOfWork + "', '" + titleIDWiki + "', '" + titleMattypeWiki + "', '" + title + "', '" + titleLowercase + "', '" + originalTitleIDWiki + "', '" + originalTitleMattypeWiki + "', '" + originalTitle + "', '" + originalTitleLowercase + "', '" + idAuthor + "', '" + author + "', '" + field + "', '" + pubDate + "', '" + bibnac + "', '" + bnd + "', '" + retrieved_date + "')")
		check(err)

		n, err := stmt.RowsAffected()
		check(err)

		if n == 0 {
			// Stop the script when no more lines are written in the database
			os.Exit(0)
		}

		if len(data.LibraryData.authors) > 0 && field == "700" {
			authorWiki = author
		}
	}

	// export only the books with published in portuguese, that have no Wikidata ID and that have an author

	if languageOfWork == "por" && len(titleIDWiki) == 0 && len(authorWiki) > 0 {

		wiki.ExportTitle(idLibrary, originalLanguageOfWork, title, originalTitle, authorWiki, pubDate, retrieved_date)
	}
}

func GetWikiTitleID(titleLowercase string, authors []string, language string) WikiData {
	var titleArray WikiData
	var idW, wMat, idWiki, wMaterial, ti, aut, url string
	idW, wMat, idWiki, wMaterial = "", "", "", ""
	const empty = ""

	ti = wiki.Replacer.Replace(titleLowercase)

	if titleLowercase != "" {
		if len(authors) > 0 {
			aut = wiki.Replacer.Replace(authors[1])

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
				fmt.Println("errI: ", errI)
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
