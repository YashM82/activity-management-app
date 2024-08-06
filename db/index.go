package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var MySqlSession *sql.DB

func ConnectMySqlDb() {

	var err error

	dsn := ("hr_common_trng_rw:rWcWVjAOheuS127@tcp(10.26.25.26:3306)/hr_common_trng")
	MySqlSession, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error in intializing connection with MySql Db :", err)
	}

	// Verify the connection
	err = MySqlSession.Ping()
	if err != nil {
		log.Println("Error in Connecting MySql Db :", err)
	}

	log.Println("MySql Db Connected Successfully")
}
