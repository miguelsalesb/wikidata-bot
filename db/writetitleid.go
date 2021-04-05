package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func WriteTitleId(idWiki string, id_library string) {

	var dbName = "titles"

	// open db
	DBCon, errDB := sql.Open("mysql", "root:@tcp(localhost:3306)/wikidata")
	if errDB != nil {
		fmt.Println(errDB)
	}
	defer DBCon.Close()

	if errPing := DBCon.Ping(); errPing != nil {
		log.Println("DATABASE CONNECTION FAILED: ", errPing)
	}

	stmt, err := DBCon.Exec("UPDATE " + dbName + " SET new_id_wikidata = ('" + idWiki + "') WHERE id_library = '" + id_library + "'")
	if err != nil {
		fmt.Println(err)
	}

	n, err := stmt.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}

	if n == 0 {
		return
	}

}
