package db

import (
	"database/sql"
	"fmt"
	"perpus_api/config"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() error {
	conf := config.GetConfig()

	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		conf.DB_HOST, conf.DB_PORT, conf.DB_USER, conf.DB_PASSWORD, conf.DB_NAME)

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	fmt.Println("Connected to the database")
	return nil

}

func GetDB() *sql.DB {
	return db
}
