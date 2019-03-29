package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=" + os.Getenv("PGUSER") +
		" dbname   = " + os.Getenv("PGDATABASE") +
		" password = " + os.Getenv("PGPASSWORD") +
		" host     = " + os.Getenv("PGHOST") +
		" port     = " + os.Getenv("PGPORT")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS values (number integer)")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "redis host %s", os.Getenv("REDIS_HOST"))
		fmt.Fprintf(response, "redis port %s", os.Getenv("REDIS_PORT"))
		fmt.Fprintf(response, "pg    host %s", os.Getenv("PG_HOST"))
		fmt.Fprintf(response, "pg    host %s", os.Getenv("PG_HOST"))
	})

	http.HandleFunc("/values/all", func(response http.ResponseWriter, request *http.Request) {
		rows, err := db.Query("SELECT * FROM values")
		if err != nil {
			log.Fatal(err)
		}
		for rows.Next() {
			var value int64
			if err := rows.Scan(&value); err != nil {
				log.Fatal(err)
			}
			fmt.Fprint(response, value)
		}
	})

	http.HandleFunc("/values", func(response http.ResponseWriter, request *http.Request) {
		if err := request.ParseForm(); err != nil {
			log.Fatal(err)
		}
		index := request.FormValue("index")
		_, err = db.Exec("INSERT INTO values (" + index + ")")
		if err != nil {
			log.Fatal(err)
		}
	})

	http.ListenAndServe(":5000", nil)
}
