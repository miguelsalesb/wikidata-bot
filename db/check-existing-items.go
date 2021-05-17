package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func CheckDatabaseItem(dbName string, id_library string) bool {

	var (
		newID, new_id_wikidata string
		idExists               bool
	)

	row, err := DBCon.Query("SELECT new_id_wikidata FROM " + dbName + " WHERE id_library = '" + id_library + "'")
	if err != nil {
		fmt.Println(err)
	}
	defer row.Close()

	for row.Next() {
		err := row.Scan(&new_id_wikidata)
		if err != nil {
			fmt.Println(err)
		}
		newID = new_id_wikidata
		

	}

	if len(newID) > 0 {
		idExists = true
	} else {
		idExists = false
	}
	return idExists
}
