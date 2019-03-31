package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"database/sql"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

type RequestStructure struct {
	Index string
}

func dumpHTTPRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func main() {
	connStr := "user=" + os.Getenv("PGUSER") +
		" dbname   = " + os.Getenv("PGDATABASE") +
		" password = " + os.Getenv("PGPASSWORD") +
		" host     = " + os.Getenv("PGHOST") +
		" port     = " + os.Getenv("PGPORT") +
		" sslmode=disable "

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS values (number integer)")
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       0})

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "redis host %s\n", os.Getenv("REDIS_HOST"))
		fmt.Fprintf(response, "redis port %s\n", os.Getenv("REDIS_PORT"))
		fmt.Fprintf(response, "pg    host %s\n", os.Getenv("PGHOST"))
		fmt.Fprintf(response, "pg    port %s\n", os.Getenv("PGPORT"))
		fmt.Fprintf(response, "pg    db   %s\n", os.Getenv("PGDATABASE"))
		fmt.Fprintf(response, "pg    user %s\n", os.Getenv("PGUSER"))
	})

	http.HandleFunc("/values/all", func(response http.ResponseWriter, request *http.Request) {
		rows, err := db.Query("SELECT * FROM values")
		if err != nil {
			log.Print(err)
		}
		cols, err := rows.Columns()
		var results = []map[string]int64{}
		for rows.Next() {
			var number int64
			if err := rows.Scan(&number); err != nil {
				log.Print(err)
			}
			row := make(map[string]int64)
			row[cols[0]] = number
			results = append(results, row)
		}

		//Set Content-Type header so that clients will know how to read response
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(response)
		enc.Encode(results)
		rows.Close()

	})

	http.HandleFunc("/values/current", func(response http.ResponseWriter, request *http.Request) {
		values := redisClient.HGetAll("values")
		if values.Err() != nil {
			log.Print(err)
			return
		}

		//Set Content-Type header so that clients will know how to read response
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(response)
		enc.Encode(values.Val())
	})

	http.HandleFunc("/values", func(response http.ResponseWriter, request *http.Request) {

		request.ParseForm()
		log.Println(request.Form)

		decoder := json.NewDecoder(request.Body)
		var requestStructure RequestStructure
		err := decoder.Decode(&requestStructure)
		if err != nil {
			log.Print("Decoding error " + err.Error())
			return
		}

		log.Print("Request index: " + requestStructure.Index)

		index := requestStructure.Index
		if index = request.FormValue("index"); len(index) == 0 {
			log.Print("invalid index")
			dump, err := httputil.DumpRequestOut(request, true)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(response, "%q", dump)
			return
		}
		_, err = db.Exec("INSERT INTO values(number) VALUES (" + index + ")")
		if err != nil {
			log.Print(err)
			return
		}

		cmd := redisClient.HSet("values", index, "Nothing yet")
		if cmd.Err() != nil {
			log.Print(err)
			return
		}
		redisClient.Publish("insert", index)

		fmt.Fprintf(response, "inserted value "+index)
	})

	http.ListenAndServe(":5000", nil)
}
