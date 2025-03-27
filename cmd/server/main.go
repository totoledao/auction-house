package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi"))
	})

	http.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Subroute"))
	})

	host := ":3000"
	fmt.Printf("Server running at http://localhost%s/ 🌐", host)
	http.ListenAndServe(host, nil)
}
