package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", generate)
	http.ListenAndServe(":8080", nil)
}

func generate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Nothing to see here.")
}
