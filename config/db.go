package config

import (
	"krayshoping/utils"
	"log"

	"database/sql"

	_ "github.com/lib/pq"
)

func LoadDatabase() *sql.DB {

	dbDSN := utils.DotEnv("DB_DSN")

	db, err := sql.Open("postgres", dbDSN)
	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	} else {
		log.Println("Database Successfully Connected")
	}
	return db
}
