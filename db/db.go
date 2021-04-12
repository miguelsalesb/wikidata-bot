package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// DBCon is the connection handle for the database
	DBCon *sql.DB
)

var ReplacerDB = strings.NewReplacer("%20", " ", "%2C", ",", "%C3%80", "À", "%C3%81", "Á", "%C3%82", "Â", "%C3%83", "Ã", "%C3%84", "Ä",
	"%C3%87", "Ç", "%C3%88", "È", "%C3%89", "É", "%C3%8A", "Ê", "%C3%8B", "Ë", "%C3%8C", "Ì", "%C3%8D", "Í", "%C3%8E", "Î", "%C3%8F", "Ï",
	"%C3%92", "Ò", "%C3%93", "Ó", "%C3%94", "Ô", "%C3%95", "Õ", "%C3%96", "Ö", "%C3%99", "Ù", "%C3%9A", "Ó", "%C3%9B", "Û", "%C3%9D", "Ý",
	"%C3%A0", "à", "%C3%A1", "á", "%C3%A2", "â", "%C3%A3", "ã", "%C3%A4", "ä", "%C3%A7", "ç", "%C3%A8", "è", "%C3%A9", "é", "%C3%AA", "ê",
	"%C3%AB", "ë", "%C3%AC", "ì", "%C3%AD", "í", "%C3%AE", "î", "C3%AF", "ï", "%C3%B1", "ñ", "%C3%B2", "ò", "%C3%B3", "ó", "%C3%B4", "ô",
	"%C3%B5", "õ", "%C3%B6", "ö", "%C3%B9", "ù", "%C3%BA", "ú", "%C3%BB", "û", "%C3%BC", "ü", "%C3%BD", "ý", "\"", "'", "%C2%BA", "º",
	"%C2%AA", "ª", "%26", "&", "%23", "#", "%24", "$", "%25", "%", "\\%27", "", "%28", "(", "%29", ")", "%2D", "-", "%5B", "[", "%5D", "]",
	"%5E", "^", "%5F", "_", "%60", "`", "%7B", "{", "%7C", "|", "%7D", "}", "none,", "", ",,", "", ",", "")

func check(e error) {
	if e != nil {
		// panic(e)
		fmt.Println(e)
	}
}

func WriteOccupations(idLibrary string, nonCoinOccup []string) {

	var dbName = "occupations"

	for o := 0; o < len(nonCoinOccup); o += 4 {

		if len(nonCoinOccup[o]) > 0 {

			stmt, err := DBCon.Exec("INSERT INTO " + dbName + " (id_library, id_occupation, name_occupation, id_instance_of, name_instance_of) VALUES ('" + idLibrary + "' , '" + nonCoinOccup[o] + "', '" + ReplacerDB.Replace(nonCoinOccup[o+1]) + "', '" + nonCoinOccup[o+2] + "', '" + ReplacerDB.Replace(nonCoinOccup[o+3]) + "')")
			check(err)

			n, err := stmt.RowsAffected()
			check(err)

			if n == 0 {
				// Stop the script when no more lines are written in the database
				os.Exit(0)
			}
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
