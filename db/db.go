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
