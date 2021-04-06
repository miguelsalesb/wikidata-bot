package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"wikidata/authorities"
	"wikidata/db"
	"wikidata/logs"
	"wikidata/titles"
)

func check(e error) {
	if e != nil {
		// panic(e)
		fmt.Println(e)
	}
}

// This app gets the data of the Library authoritoes and bibliographic repositories
// searches the author's names and titles in Wikidata in order to find out if that data already exists, and:
// if the Library author's and titles don't exist in Wikidata, creates new Wikidata entries with that Library data
// if the Library author's or titles already exists, but some of the Library data doesn't, adds that new data

func main() {

	var (
		errLog, errDB, errPing error
	)

	var tables = []string{"authors", "titles", "occupations"}
	// DBCon is the connection handle for the database

	// write the logs in the logs.txt file
	logs.File, errLog = os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if errLog != nil {
		log.Println("LOG FILE ERROR: ", errLog)
	}
	defer logs.File.Close()

	// check if database exists
	//if it doesn't exist, create it
	// https://play.golang.org/p/jxza3pbqq9
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
	go authorities.GetAuthors(doneAuthorities)
	go titles.GetTitles(doneTitles)

	<-doneAuthorities
	<-doneTitles
	// time.Sleep(60 * time.Second)
}
