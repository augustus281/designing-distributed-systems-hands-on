package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I'm a sidecar ! This is a single node pattern.")
	})
	http.ListenAndServe(":8080", nil)
}
