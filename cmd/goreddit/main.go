package main

import (
	"log"
	"net/http"

	"github.com/momomo0206/goreddit/postgres"
	"github.com/momomo0206/goreddit/web"
)

func main() {
	store, err := postgres.NewStore("postgres://postgres:secret@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	h := web.NewHandler(store)
	http.ListenAndServe(":3000", h)
}
