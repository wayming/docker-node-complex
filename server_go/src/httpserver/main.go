package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "redis host %s", os.Getenv("REDIS_HOST"))
	})

	http.ListenAndServe(":5000", nil)
}
