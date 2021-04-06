package authorities

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"wikidata/db"
	"wikidata/functions"
	"wikidata/wiki"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

// Authors type with data retrieved from the author's data repository and the associated Wikidata data
type Authors struct {
	Names                 []string
	Dates                 []string
	Nationality           string
	Occupations           string
	WikiAuthors           []string
	WikiAuthorDescription string
	WikiChecked           []string
}

const dbName = "authors"

func check(e error) {
	if e != nil {
		// panic(e)
		fmt.Println(e)
	}
}

func GetAuthors(doneAuthorities chan bool, repAuthorsFirst int, repAuthorsLast int) {

	var (
		names, dates               []string
		nationality, occupations   string
		wikiresults1, wikiresults2 []string // auhor info: id's dates, nationality, occupations,

	)

	// Get data from the author's repository
	for n := repAuthorsFirst; n <= repAuthorsLast; n++ {
		fmt.Println("\nAuthorities - ", n)
		time.Sleep(350 * time.Millisecond)
		if n%500 == 0 {
			time.Sleep(120 * time.Second)
		}

		urlMarcxchange := fmt.Sprintf("%s%d", "http://urn.bn.pt/nca/unimarc-authorities/marcxchange?id=", n) // web address of the authorities repository in the marcxchange format

		if IsAuthor(urlMarcxchange) { // check first if it is a personal name authority
			names = GetAuthorNames(urlMarcxchange)
			dates = GetDates(urlMarcxchange)
			nationality = GetNacionality(urlMarcxchange)
			occupations = GetOccupations(urlMarcxchange)
			wikiresults1, wikiresults2 = wiki.GetWikiContent(names)

			author := Authors{
				Names:       names,
				Dates:       dates,
				Nationality: nationality,
				Occupations: occupations,
				WikiAuthors: wikiresults1,
				WikiChecked: wikiresults2,
			}
			WriteAuthors(author, db.DBCon)
		}
	}
	// Finishes only afterwards all the authorities are scrapped and wrote in the database
	doneAuthorities <- true
}

// IsAuthor - checks if it is a personal name author (instead of corporate, geographic, etc,)
func IsAuthor(urlMarcxchange string) bool {

	var (
		leader           string
		isRecord, has200 bool
	)

	res, err := http.Get(urlMarcxchange)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// get body page content
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	leader = doc.Find("leader").Contents().Text() // get leader field data

	// There may be some cases where it has the correct leadeer code and there insn't field 200
	doc.Find("datafield").Each(func(i int, s *goquery.Selection) {
		tag, _ := s.Attr("tag")
		if tag == "200" { // does it have field 200
			has200 = true
		}
	})
	if (strings.Contains(leader, "cx  a") || strings.Contains(leader, "cx a") || strings.Contains(leader, "nx  a") || strings.Contains(leader, "nx a")) && has200 == true {
		isRecord = true // it is personal name author's info
	} else if strings.Contains(leader, "cx  a") || strings.Contains(leader, "nx  a") || has200 == false {
		isRecord = false
	}
	return isRecord
}

// GetAuthorNames - get author's name
func GetAuthorNames(urlMarcxchange string) []string {

	const empty = ""
	var (
		id, name, surname, initials, ref string
		replacer                         = strings.NewReplacer(",", "", "\\", "", "'", "\\'", "<", "", ">", "")
	)
	var authInfo = make([]string, 0, 6)
	res, err := http.Get(urlMarcxchange)
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
		// To get all authors, put tag[:1]
		if tag == "200" { // get data from field 200 (primary responsability)
			controlfield := doc.Find("controlfield")
			if tag, ok := controlfield.Attr("tag"); tag == "001" {
				if ok {
					id = controlfield.First().Text()
					authInfo = append(authInfo, id)
				} else {
					authInfo = append(authInfo, empty)
				}
			}
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "a" {
					nm := e.Text()
					name = replacer.Replace(nm)
				}
				if attr, _ := e.Attr("code"); attr == "b" {
					sn := e.Text()
					surname = replacer.Replace(sn)
				}
				if attr, ok := e.Attr("code"); attr == "c" {
					if ok {
						in := e.Text()
						initials = replacer.Replace(in)
					}
				}
			})
			if name == "" {
				authInfo = append(authInfo, empty)
			} else {
				authInfo = append(authInfo, name)
			}
			if surname == "" {
				authInfo = append(authInfo, empty)
			} else {
				authInfo = append(authInfo, surname)
			}
			if initials == "" {
				authInfo = append(authInfo, empty)
			} else {
				authInfo = append(authInfo, initials)
			}
			authInfo = append(authInfo, empty)
			authInfo = append(authInfo, "200")
		}

		if tag == "400" { // get data from the 400 field (variant names)
			controlfield := doc.Find("controlfield")
			if tag, ok := controlfield.Attr("tag"); tag == "001" {
				if ok {
					id = controlfield.First().Text()
					authInfo = append(authInfo, id)
				} else {
					authInfo = append(authInfo, empty)
				}
			}
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "a" {
					nm := e.Text()
					name = replacer.Replace(nm)
				}
				if attr, _ := e.Attr("code"); attr == "b" {
					sn := e.Text()
					surname = replacer.Replace(sn)
				}
				if attr, ok := e.Attr("code"); attr == "c" {
					if ok {
						in := e.Text()
						initials = replacer.Replace(in)
					}
				}
				if attr, ok := e.Attr("code"); attr == "3" {
					if ok {
						r := e.Text()
						ref = replacer.Replace(r)
					}
				}
			})
			if name == "" {
				authInfo = append(authInfo, empty)
			} else {
				authInfo = append(authInfo, name)
			}
			if surname == "" {
				authInfo = append(authInfo, empty)
			} else {
				authInfo = append(authInfo, surname)
			}
			if initials == "" {
				authInfo = append(authInfo, empty)
			} else {
				authInfo = append(authInfo, initials)
			}
			if ref == "" {
				authInfo = append(authInfo, empty)
			} else {
				authInfo = append(authInfo, ref)
			}
			authInfo = append(authInfo, "400")
		}
	})
	return authInfo
}

// GetDates - get author birth and death dates
func GetDates(urlMarcxchange string) []string {
	var (
		dates                                                  []string
		birthDate200, deathDate200, birthDate400, deathDate400 string
		countBD, countBD400                                    int
		replacer                                               = strings.NewReplacer("ca ", "")
	)
	res, err := http.Get(urlMarcxchange)
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
		// To get all authors, put tag[:1], "2"
		if strings.Contains(tag[:3], "200") {
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "f" {
					if countBD == 0 {
						// For when you have the birth date and the death date
						r, _ := regexp.Compile("[0-9]{4}-[0-9]{4}")
						d := e.Text()
						dd := replacer.Replace(d)
						dt := r.FindString(dd)
						s, _ := regexp.Compile("[0-9]{4}-")
						e := e.Text()
						dts := s.FindString(e)
						hyphens400 := strings.Contains(e, "--")
						hyphens200 := strings.Contains(d, "--")

						if len(dt) == 9 && hyphens200 != true {
							i := strings.Index(dt, "-")
							birthDate200 = dt[i-4 : i]
							deathDate200 = dt[i+1:]
						} else if len(dts) == 5 && len(dts) != 9 && hyphens400 != true {
							i := strings.Index(dts, "-")
							birthDate200 = dts[i-4 : i]
							deathDate200 = ""
						}
					}
					countBD++
				}
			})
		}
	})

	doc.Find("datafield").Each(func(i int, s *goquery.Selection) {

		tag, _ := s.Attr("tag")
		if strings.Contains(tag[:3], "400") {
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "f" {
					if countBD400 == 0 {
						r, _ := regexp.Compile("[0-9]{4}-[0-9]{4}")
						d := e.Text()
						dt := r.FindString(d)
						hyphens200 := strings.Contains(d, "--")
						s, _ := regexp.Compile("[0-9]{4}-")
						e := e.Text()
						dts := s.FindString(e)
						hyphens400 := strings.Contains(e, "--")

						if len(dt) == 9 && hyphens200 != true {
							i := strings.Index(dt, "-")
							birthDate400 = dt[i-4 : i]
							deathDate400 = dt[i+1:]
						} else if len(dts) == 5 && len(dts) != 9 && hyphens400 != true {
							i := strings.Index(dts, "-")
							birthDate400 = dts[i-4 : i]
							deathDate400 = ""
						}
					}
					countBD400++
				}
			})
		}
	})
	if birthDate200 != "" {
		dates = append(dates, birthDate200, deathDate200)
	} else {
		dates = append(dates, birthDate400, deathDate400)
	}
	return dates
}

// GetNacionality - get author nationality
func GetNacionality(urlMarcxchange string) string {
	var nationality string
	var nation102, nation830 string
	var count102, count830 int

	res, err := http.Get(urlMarcxchange)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	doc.Find("datafield").Each(func(i int, s *goquery.Selection) {
		tag, _ := s.Attr("tag")
		if tag == "102" {
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "a" {
					if count102 == 0 {
						nation102 = e.Text()
					}
					count102++
				}
				if count102 == 0 {
					nation102 = ""
				}
			})
		}
	})

	doc.Find("datafield").Each(func(i int, s *goquery.Selection) {
		tag, _ := s.Attr("tag")
		if tag == "830" {
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "9" {
					if count830 == 0 {
						nation830 = e.Text()
					}
					count830++
				}
				if count830 == 0 {
					nation830 = ""
				}
			})
		}
	})
	if nation102 != "" {
		nationality = nation102
	} else {
		nationality = nation830
	}
	return nationality
}

func GetOccupations(urlMarcxchange string) string {
	var occupations, occupationsFinal string
	var occupationsArray []string
	var replacer = strings.NewReplacer("\\", "", "'", "\\'")

	res, err := http.Get(urlMarcxchange)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	doc.Find("datafield").Each(func(i int, s *goquery.Selection) {
		tag, _ := s.Attr("tag")
		if tag == "830" {
			s.Find("subfield").Each(func(i int, e *goquery.Selection) {
				if attr, _ := e.Attr("code"); attr == "a" {
					nm := e.Text()
					occupations = replacer.Replace(nm)
					occupationsArray = append(occupationsArray, occupations)
				}
			})
		}
	})
	occupationsFinal = strings.Join(occupationsArray, ", ")
	return occupationsFinal
}

// WriteAuthors - write data in the database
func WriteAuthors(author Authors, dbs *sql.DB) {

	var (
		idLibrary, name, surname, initials, ref, field, sameOccupations, sameBirthdate, sameDeathdate,
		sameNationality, coincidentalOccupations, nonCoincidentalOccupations, nonCoincidentalOccupationsWithoutInstanceOf,
		libraryAuthorPublicDomain, wikidataAuthorPublicDomain, existsInWikiFinal, retrieved_date, nationalityLibrary string
		coincidOccup, nonCOccup, nonCoinOccup, nonCoinOccupWithoutInstanceOf, surnameArray, nameArray, authorWikiArray []string
		existsInWiki                                                                                                   int
	)
	var replacer = strings.NewReplacer(",,", ",", "none", "")
	var counterX, counterZ = 0, 0

	// firstAuthorName, secondAuthorName, authorNameArray - to calculate the probability of the author already being in Wikidata

	for x := 0; x < len(author.WikiAuthors); x += 10 {
		surnameArray = nil
		nameArray = nil
		if x <= 0 {
			idLibrary = author.Names[x]
			name = author.Names[x+1]
			surname = author.Names[x+2]
			initials = author.Names[x+3]
			ref = author.Names[x+4]
			field = author.Names[x+5]
			surnameArray = append(surnameArray, surname)
			nameArray = append(nameArray, name)

		} else if x > 0 {
			counterX = counterX + 4
			counterZ = counterZ + 9
			y := x - counterX
			idLibrary = author.Names[y]
			name = author.Names[y+1]
			surname = author.Names[y+2]
			initials = author.Names[y+3]
			ref = author.Names[y+4]
			field = author.Names[y+5]
			surnameArray = append(surnameArray, surname)
			nameArray = append(nameArray, name)
		}
		// To help calculate the probability of the author already existing in Wikidata
		surnameString := strings.Join(surnameArray, ",")
		nameString := strings.Join(nameArray, ",")
		completeNameString := surnameString + " " + nameString

		numberOfWordsSurname := len(strings.Fields(surnameString))
		numberOfWordsName := len(strings.Fields(nameString))

		checkWordsCompleteName := checkWords(completeNameString)

		idWikidata := author.WikiAuthors[x]
		idLibraryWikidata := author.WikiAuthors[x+1]
		birthDateWikidata := author.WikiAuthors[x+2]
		deathDateWikidata := author.WikiAuthors[x+3]
		nationalityWikidata := author.WikiAuthors[x+4]
		signature := author.WikiAuthors[x+5]
		image := author.WikiAuthors[x+6]
		occupationsWikidata := author.WikiAuthors[x+7]
		notableworkWikidata := author.WikiAuthors[x+8]
		authorDescriptionWikidata := author.WikiAuthors[x+9]

		birthDateLibrary := author.Dates[0]
		deathDateLibrary := author.Dates[1]
		nationalityLibrary = author.Nationality
		occupationsLibrary := author.Occupations

		// To check if the author's information in Wikidata is the same as the Library author's standardized heading
		sameAsField200 := author.WikiChecked[0]

		if birthDateLibrary == birthDateWikidata {
			sameBirthdate = "true"
		} else if birthDateLibrary != birthDateWikidata {
			sameBirthdate = "false"
		} else if birthDateLibrary == "" || birthDateWikidata == "" {
			sameBirthdate = ""
		}

		if deathDateLibrary == deathDateWikidata {
			sameDeathdate = "true"
		} else if deathDateLibrary != deathDateWikidata {
			sameDeathdate = "false"
		} else if deathDateLibrary == "" || deathDateWikidata == "" {
			sameDeathdate = ""
		}

		countryName := functions.CheckCountryName(nationalityLibrary)

		if nationalityLibrary != "" && nationalityWikidata != "" && nationalityWikidata == countryName {
			sameNationality = "true"
		} else if nationalityLibrary != "" && nationalityWikidata != "" && nationalityWikidata != countryName {
			sameNationality = "false"
		} else {
			sameNationality = ""
		}

		// Get the atual year to calculate which authors works have entered in the public domain
		t := time.Now()
		year := t.Year()

		libraryDeathDate, err := strconv.Atoi(deathDateLibrary)
		wikidataDeathDate, err := strconv.Atoi(deathDateWikidata)

		if year-libraryDeathDate > 70 {
			libraryAuthorPublicDomain = "true"
		} else {
			libraryAuthorPublicDomain = "false"
		}

		if year-wikidataDeathDate > 70 {
			wikidataAuthorPublicDomain = "true"
		} else {
			wikidataAuthorPublicDomain = "false"
		}

		if occupationsLibrary != "" {
			sameOccupations, coincidOccup, nonCOccup = compareOccupations(strings.ToLower(occupationsLibrary), strings.ToLower(occupationsWikidata))
			coincidentalOccupations = strings.Join(coincidOccup, ", ")
		} else {
			sameOccupations, coincidOccup, nonCOccup = "", nil, nil
		}

		if len(nonCOccup) > 0 {
			// first I was only passing the different occupations, but the comparison between the different occupations and the same occupations
			// should be done after getting all the different occupations wiki data, because of the
			// for the cases in which the library occupation has more than one word and the Wikidata has one, but have the same meaning, ex.: "político português" and "político", then convert it to an empty string inside the nonCoincidentalOccupations slice

			nonCoinOccup, nonCoinOccupWithoutInstanceOf = wiki.GetOccupationsWiki(nonCOccup) // get an array with the occupations ID and descriptions
			// The occupation and its label

			nonCoincidOccupWithoutInstanceOf := strings.Join(nonCoinOccupWithoutInstanceOf, ",")
			nonCoincidentaldOccupWithoutInstanceOf := strings.Trim(nonCoincidOccupWithoutInstanceOf, ",")
			nonCoincidentalOccupationsWithoutInstanceOf = replacer.Replace(nonCoincidentaldOccupWithoutInstanceOf)
		} else {
			nonCoinOccup, nonCoinOccupWithoutInstanceOf = nil, nil
		}

		// Take away the multiples commas
		replacerNonCoincidental := strings.NewReplacer(",,", "")
		nonCoincidentalOccupationsWithoutInstanceOf = replacerNonCoincidental.Replace(nonCoincidentalOccupationsWithoutInstanceOf)

		// Calculation of the probability of the author already existing in Wikidata

		// If it doesn't have the words: "de", "da", "das", "do" and "dos", and the total number of words of the name are 2
		if !checkWordsCompleteName && (numberOfWordsSurname+numberOfWordsName) == 2 {
			existsInWiki = 1

			// If it has the words: "de", "da", "das", "do" and "dos", and the total number of words of the name are 2
		} else if checkWordsCompleteName && (numberOfWordsSurname+numberOfWordsName) == 2 {
			existsInWiki = 2

			// If it has the words: "de", "da", "das", "do" and "dos", and the total number of words of the name greater than 2 (but not 3)
		} else if checkWordsCompleteName && (numberOfWordsSurname+numberOfWordsName) > 2 && (numberOfWordsSurname+numberOfWordsName) != 3 {
			existsInWiki = 0

			// If it has the words: "de", "da", "das", "do" and "dos", and the total number of words are 3
		} else if checkWordsCompleteName && (numberOfWordsSurname+numberOfWordsName) == 3 {
			existsInWiki = 1

			// If it doesn't have the words: "de", "da", "das", "do" and "dos", and the total number of words of the name are greater than 2
		} else if !checkWordsCompleteName && (numberOfWordsSurname+numberOfWordsName) > 2 {
			existsInWiki = 0

			// If it doesn't have the words: "de", "da", "das", "do" and "dos", and the total number of words of the name are greater than 3
		} else if checkWordsCompleteName && (numberOfWordsSurname+numberOfWordsName) > 3 {
			existsInWiki = 0
		}

		// If there are coincidental occupations in both the library catalogue and Wikidata
		if len(coincidOccup) > 0 {
			existsInWiki += 1
		}

		// If the author already exists in Wikidata
		if idWikidata != "" {
			existsInWiki += 1
		}

		// Don't want to export the author's that for instance have only a first name, like: "A. C. L.", or "A.D."
		if numberOfWordsName == 0 || numberOfWordsSurname == 0 {
			existsInWiki = 5
		}

		// Don't export if the author already exists im Wikidata
		if len(idWikidata) > 0 {
			existsInWiki = 5
		}

		existsInWikiFinal = strconv.Itoa(existsInWiki)

		// $timestamp must resemble ISO 8601
		retrieved_date = time.Now().Format("+2006-01-02T00:00:00Z")

		authorWikiArray = append(authorWikiArray, idLibrary, surname+" "+name, birthDateLibrary, deathDateLibrary, nationalityWikidata, occupationsLibrary, field, retrieved_date, existsInWikiFinal)

		// write the author data in the authors table
		stmt, err := dbs.Exec("INSERT INTO " + dbName + " (id_library, id_library_wikidata, id_wikidata, name, surname, initials, author_description_wikidata, birth_date_library, death_date_library, birth_date_wikidata, death_date_wikidata, nationality_library, nationality_wikidata, ref, field, occupations_library, occupations_wikidata, same_occupations, coincidental_occupations, non_coincidental_occupations, notablework_wikidata, same_as_field200, exists_in_wiki, same_birthdate, same_deathdate, same_nationality, library_author_public_domain, wikidata_author_public_domain, signature, image, retrieved_date) VALUES (" + idLibrary + " , '" + idLibraryWikidata + "','" + idWikidata + "', '" + name + "', '" + surname + "', '" + initials + "', '" + authorDescriptionWikidata + "', '" + birthDateLibrary + "', '" + deathDateLibrary + "', '" + birthDateWikidata + "', '" + deathDateWikidata + "', '" + nationalityLibrary + "', '" + nationalityWikidata + "', '" + ref + "', '" + field + "', '" + occupationsLibrary + "', '" + occupationsWikidata + "', '" + sameOccupations + "', '" + coincidentalOccupations + "', '" + nonCoincidentalOccupationsWithoutInstanceOf + "', '" + notableworkWikidata + "', '" + sameAsField200 + "', '" + existsInWikiFinal + "', '" + sameBirthdate + "', '" + sameDeathdate + "', '" + sameNationality + "', '" + libraryAuthorPublicDomain + "', '" + wikidataAuthorPublicDomain + "', '" + signature + "', '" + image + "', '" + retrieved_date + "')")
		check(err)

		n, err := stmt.RowsAffected()
		check(err)

		if n == 0 {
			// Stop the script when no more lines are written in the database
			os.Exit(0)
		}

	}

	// just export the portuguese authors
	if nationalityLibrary == "PT" {
		wiki.ExportAuthor(authorWikiArray)
	}

	if len(nonCOccup) > 1 && nonCOccup[0] != "" {
		// The occupation and its label, and the occupation instanceOf and its label (Wikidata info)
		nonCoincidOccup := strings.Join(nonCoinOccup, ",")
		nonCoincidentalOccupations = replacer.Replace(nonCoincidOccup)
		writeOccupations := strings.Split(nonCoincidentalOccupations, ",")
		db.WriteOccupations(idLibrary, writeOccupations)
	}
}

// check if the author's name has any prepositions to help calculate the probability of it existing in Wikidata
func checkWords(check string) bool {

	var name = strings.Split(check, " ")

	var checkInsideNameArray bool
	for _, v := range name {
		if v == "de" || v == "da" || v == "das" || v == "do" || v == "dos" {
			checkInsideNameArray = true
		}
	}
	return checkInsideNameArray
}

// write the occupations in a different table

func compareOccupations(occupLibrary string, occupWiki string) (string, []string, []string) {
	var (
		occupationsLibrary, occupationsWikidata, occupationsL, occupationsW, coincidOccupations,
		coincidentalOccupations, noncoincidOccup, nonCoincidentalOccupations []string
		coincidental string
		regexp       = regexp.MustCompile("[^\\s]+")
	)

	replacer := strings.NewReplacer(".", ",", " e ", ",")
	occupationsLib := replacer.Replace(occupLibrary) // to put into an array the individual occupations. Ex.: Romancista, dramaturgo, historiador, crítico literário e memorialista. Prof. univ., desempenhou cargos públicos

	occupationsLibrary = strings.Split(occupationsLib, ",") // put in an array
	occupationsWikidata = strings.Split(occupWiki, ",")     // to get the terms divided by commas and dots and not all the words separated

	for x := 0; x < len(occupationsLibrary); x++ {
		occupationsL = regexp.FindAllString(strings.Trim(occupationsLibrary[x], " "), -1) // to divide all the words in order to compare individual words

		// add the occupations that have more than one word, afterwards, if there is an occupation with two words that has a Wikidata ID
		// and the two words that compose that more than two words, have also WD ID's, then remove the singular words
		if len(occupationsL) >= 2 {
			noncoincidOccup = append(noncoincidOccup, strings.Trim(occupationsLibrary[x], " "))
		}
		for y := 0; y < len(occupationsWikidata); y++ {
			occupationsW = regexp.FindAllString(strings.Trim(occupationsWikidata[y], " "), -1)

			for t := 0; t < len(occupationsL); t++ {
				occupationL := occupationsL[t]
				occupationSizeL := len(occupationsL[t])
				occupationWithoutLastLetterL := occupationL[:occupationSizeL-1]
				occupationLastLetterL := occupationL[occupationSizeL-1:]

				noncoincidOccup = append(noncoincidOccup, strings.Trim(occupationsL[t], " "))

				// in order to verify the occupation in Wikidata, I removed the last letter
				if occupationLastLetterL == "a" && occupationWithoutLastLetterL == strings.TrimLeft(occupationsWikidata[y], " ") && len(strings.Trim(occupationL, " ")) > 2 {
					coincidOccupations = append(coincidOccupations, strings.Trim(occupationL, " "))
					noncoincidOccup = append(noncoincidOccup, strings.Trim(occupationL, " ")) // apend the Library occupation to compare and remove it later
				}

				if occupationsL[t] == strings.TrimLeft(occupationsWikidata[y], " ") && len(strings.Trim(occupationsWikidata[y], " ")) > 2 {
					coincidOccupations = append(coincidOccupations, strings.Trim(occupationsWikidata[y], " "))
					noncoincidOccup = append(noncoincidOccup, strings.Trim(occupationsLibrary[x], " ")) // apend the Library occupation to compare and remove it later
				}
				// for the cases - médico vs. médica
				for u := 0; u < len(occupationsW); u++ {
					occupationW := occupationsW[u]
					occupationSizeW := len(occupationsW[u])
					occupationWithoutLastLetterW := occupationW[:occupationSizeW-1]

					if occupationWithoutLastLetterL == occupationWithoutLastLetterW && len(strings.Trim(occupationW, " ")) > 2 {
						coincidOccupations = append(coincidOccupations, strings.Trim(occupationW, " "))
					}
				}

				for u := 0; u < len(occupationsW); u++ {
					if occupationsW[u] == occupationsL[t] && len(strings.Trim(occupationsW[u], " ")) > 2 {
						coincidOccupations = append(coincidOccupations, strings.Trim(occupationsW[u], " "))
						// noncoincidOccup = append(noncoincidOccup, strings.Trim(occupationsLibrary[x], " "))
					}
				}
			}
		}
	}

	if len(coincidOccupations) > 0 {
		coincidental = "true"
	} else {
		coincidental = "false"
	}

	coincidentalOccupations = functions.RemoveDuplicateValues(coincidOccupations)
	nonCoincidentalOccupations = functions.RemoveDuplicateValues(noncoincidOccup)
	nonCoincidentalOccupations = functions.Diff(nonCoincidentalOccupations, coincidentalOccupations)

	// for the cases tin which the library occupation has more than one word and the Wikidata has one, but have the same meaning, ex.: "político português" and "político", then convert it to an empty string inside the nonCoincidentalOccupations slice
	for _, c := range coincidentalOccupations {
		for i, n := range nonCoincidentalOccupations {
			if strings.Contains(c, n) {
				nonCoincidentalOccupations[i] = ""
			}
		}
	}
	return coincidental, coincidentalOccupations, nonCoincidentalOccupations
}
