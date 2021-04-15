package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"wikidata/authorities"
	"wikidata/db"
	"wikidata/logs"
	"wikidata/titles"
	"wikidata/wiki"
)

// This app gets the data of the Library authoritoes and bibliographic repositories
// searches the author's names and titles in Wikidata in order to find out if that data already exists, and:
// if the Library author's and titles don't exist in Wikidata, creates new Wikidata entries with that Library data
// if the Library author's or titles already exists, but some of the Library data doesn't, adds that new data

func main() {

	var (
		errLog, errDB, errPing                                                   error
		repTitlesFirst, repTitlesLast, repAuthorsFirst, repAuthorsLast, populate string
		authorFirst, authorLast, titleFirst, titleLast                           int
	)

	var tables = []string{"authors", "titles", "occupations"}

	fmt.Print("Do you want to populate Wikidata with some initial properties and entities (yes or no) \n ")
	fmt.Scanln(&populate)

	if populate != "yes" && populate != "no" {
		fmt.Print("\nPlease type yes or no \n ")
		fmt.Scanln(&populate)
	}

	fmt.Print("Insert the number from the authors repository where you want to start and then press Enter \n ")
	fmt.Scanln(&repAuthorsFirst)
	authorFirst, _ = strconv.Atoi(repAuthorsFirst)

	fmt.Print("Insert the number from the authors repository where you want to finish and then press Enter \n ")
	fmt.Scanln(&repAuthorsLast)
	authorLast, _ = strconv.Atoi(repAuthorsLast)

	if authorLast < authorFirst {
		fmt.Print("\nThe finish number has to be greater than the starting number. Please insert the finish number and press Enter \n ")
		fmt.Scanln(&repAuthorsLast)
		authorLast, _ = strconv.Atoi(repAuthorsLast)
	}

	fmt.Print("Insert the number from the titles repository where you want to start and then press Enter \n ")
	fmt.Scanln(&repTitlesFirst)
	titleFirst, _ = strconv.Atoi(repTitlesFirst)

	fmt.Print("Insert the number from the titles repository where you want to finish and then press Enter \n ")
	fmt.Scanln(&repTitlesLast)
	titleLast, _ = strconv.Atoi(repTitlesLast)

	if repTitlesLast < repTitlesFirst {
		fmt.Print("\nThe finish number has to be greater than the starting number. Please insert the finish number and press Enter \n ")
		fmt.Scanln(&repTitlesLast)
		titleLast, _ = strconv.Atoi(repTitlesLast)
	}

	// write the logs in the logs.txt file
	logs.File, errLog = os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if errLog != nil {
		log.Println("LOG FILE ERROR: ", errLog)
	}
	defer logs.File.Close()

	// check if database exists
	//if it doesn't exist, create it
	// from: https://play.golang.org/p/jxza3pbqq9
	const TEST_ROOT_URI = "root:@tcp(localhost:3306)/?charset=utf8mb4&autocommit=true"

	dba, err := sql.Open("mysql", TEST_ROOT_URI)
	if err != nil {
		log.Fatal(errDB)
	}

	_, err = dba.Exec("CREATE DATABASE IF NOT EXISTS wikidata")
	if err != nil {
		log.Fatal(errDB)
	}

	dba.Close()

	// open db
	db.DBCon, errDB = sql.Open("mysql", "root:@tcp(localhost:3306)/wikidata")
	if errDB != nil {
		log.Fatal(errDB)
	}
	defer db.DBCon.Close()

	if errPing = db.DBCon.Ping(); errPing != nil {
		log.SetOutput(logs.File)
		log.Println("DATABASE CONNECTION FAILED: ", errPing)
	}

	if populate == "yes" {
		wiki.CreateFirst()
	}

	// check if the tables exist
	//if not, create them
	for _, v := range tables {

		_, table_check := db.DBCon.Query("select * from " + v)

		if table_check != nil && v == "authors" {
			db.CreateTableAuthors()
		}

		if table_check != nil && v == "titles" {
			db.CreateTableTitles()
		}

		if table_check != nil && v == "occupations" {
			db.CreateTableOccupations()
		}
	}

	doneAuthorities := make(chan bool)
	doneTitles := make(chan bool)
	// https://medium.com/@ishagirdhar/import-cycles-in-golang-b467f9f0c5a0
	go authorities.GetAuthors(doneAuthorities, authorFirst, authorLast)
	go titles.GetTitles(doneTitles, titleFirst, titleLast)

	<-doneAuthorities
	<-doneTitles
	// time.Sleep(60 * time.Second)

	t := time.Now().Format("02-01-2006")
	log.SetOutput(logs.File)
	log.Printf("%v - The title repository was scrapped from record: %v to record %v", t, strings.Trim(repAuthorsFirst, " "), strings.Trim(repAuthorsLast, " "))
	log.Printf("%v - The title repository was scrapped from record: %v to record %v", t, strings.Trim(repTitlesFirst, " "), strings.Trim(repTitlesLast, " "))
}
