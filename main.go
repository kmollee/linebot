package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := "80"
	if v := os.Getenv("PORT"); len(v) > 0 {
		port = v
	}
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "Hello world")
}
