package main

import (
    "fmt"
	"net/http"
	
	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello there!")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)

    http.ListenAndServe(":3000", r)
}