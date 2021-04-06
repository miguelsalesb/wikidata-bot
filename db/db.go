package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// DBCon is the connection handle for the database
	DBCon *sql.DB
)

func check(e error) {
	if e != nil {
		// panic(e)
		fmt.Println(e)
	}
}

func WriteOccupations(idLibrary string, nonCoinOccup []string) {

	var dbName = "occupations"

	for o := 0; o < len(nonCoinOccup); o += 4 {

		stmt, err := DBCon.Exec("INSERT INTO " + dbName + " (id_library, id_occupation, name_occupation, id_instance_of, name_instance_of) VALUES ('" + idLibrary + "' , '" + nonCoinOccup[o] + "', '" + nonCoinOccup[o+1] + "', '" + nonCoinOccup[o+2] + "', '" + nonCoinOccup[o+3] + "')")
		check(err)

		n, err := stmt.RowsAffected()
		check(err)

		if n == 0 {
			// Stop the script when no more lines are written in the database
			os.Exit(0)
		}
	}
}

func CreateTableAuthors() {

	createTableAuthors, _ := DBCon.Exec(`	
	CREATE TABLE authors (
		id int(11) NOT NULL AUTO_INCREMENT,
		id_library varchar(11) NOT NULL,
		id_library_wikidata varchar(11) DEFAULT NULL,
		id_wikidata varchar(50) DEFAULT NULL,
		new_id_wikidata varchar(50) NOT NULL,
		name text NOT NULL,
		surname text DEFAULT NULL,
		initials text DEFAULT NULL,
		author_description_wikidata text DEFAULT NULL,
		birth_date_library varchar(4) DEFAULT NULL,
		death_date_library varchar(4) DEFAULT NULL,
		birth_date_wikidata varchar(4) DEFAULT NULL,
		death_date_wikidata varchar(4) DEFAULT NULL,
		nationality_library text DEFAULT NULL,
		nationality_wikidata text DEFAULT NULL,
		ref varchar(50) NOT NULL,
		field varchar(4) NOT NULL,
		occupations_library text DEFAULT NULL,
		occupations_wikidata text DEFAULT NULL,
		same_occupations text NOT NULL,
		coincidental_occupations text NOT NULL,
		non_coincidental_occupations text NOT NULL,
		notablework_wikidata text DEFAULT NULL,
		same_as_field200 varchar(10) DEFAULT NULL,
		exists_in_wiki int(10) DEFAULT NULL,
		same_birthdate varchar(10) DEFAULT NULL,
		same_deathdate varchar(10) DEFAULT NULL,
		same_nationality varchar(10) DEFAULT NULL,
		library_author_public_domain varchar(5) NOT NULL,
		wikidata_author_public_domain varchar(5) NOT NULL,
		signature text NOT NULL,
		image text NOT NULL,
		retrieved_date varchar(50) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)

	n, err := createTableAuthors.RowsAffected()
	check(err)

	if n == 0 {
		return
	} else {
		// Stop the script when no more lines are written in the database
		os.Exit(0)
	}
}

func CreateTableTitles() {

	createTableTitles, _ := DBCon.Exec(`	
	CREATE TABLE titles (
		id int(11) NOT NULL AUTO_INCREMENT,
		id_library varchar(10) NOT NULL,
		language_of_work varchar(10) NOT NULL,
		original_language_of_work varchar(10) NOT NULL,
		title_id_wiki varchar(50) NOT NULL,
		new_id_wikidata varchar(50) NOT NULL,
		title_mattype_wiki varchar(50) NOT NULL,
		title text NOT NULL,
		title_lowercase text NOT NULL,
		original_title_id_wiki varchar(50) NOT NULL,
		original_title_mattype_wiki text NOT NULL,
		original_title text NOT NULL,
		original_title_lowercase text NOT NULL,
		id_author varchar(10) NOT NULL,
		author text NOT NULL,
		field varchar(5) NOT NULL,
		pub_date varchar(10) NOT NULL,
		bibnac text NOT NULL,
		bnd varchar(255) NOT NULL,
		retrieved_date varchar(50) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)

	n, err := createTableTitles.RowsAffected()
	check(err)

	if n == 0 {
		return
	} else {
		// Stop the script when no more lines are written in the database
		os.Exit(0)
	}
}

func CreateTableOccupations() {

	createTableOccupations, _ := DBCon.Exec(`	
	CREATE TABLE occupations (
		id int(11) NOT NULL AUTO_INCREMENT,
		id_library varchar(50) NOT NULL,
		id_occupation text NOT NULL,
		name_occupation text NOT NULL,
		id_instance_of varchar(50) NOT NULL,
		name_instance_of text NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)

	n, err := createTableOccupations.RowsAffected()
	check(err)

	if n == 0 {
		return
	} else {
		// Stop the script when no more lines are written in the database
		os.Exit(0)
	}
}
