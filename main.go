package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", generate)
	http.ListenAndServe(":"+os.Getenv("HTTP_PLATFORM_PORT"), nil)
}

func generate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Nothing to see here.")
}
