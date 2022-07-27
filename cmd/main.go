package main

import (
	"log"
	"net/http"

	"example.com/recallcards/pkg/api"
	"example.com/recallcards/pkg/cards"
	"example.com/recallcards/pkg/storage"
)

func main() {
  rep := storage.NewMemoryRepository()
  srv := cards.NewCardService(rep)
  c := api.NewController(srv)

  http.HandleFunc("/cards/create/", c.CreateCard)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

