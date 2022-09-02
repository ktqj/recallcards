package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/recallcards/pkg/cards"
)

type Controller interface {
	CreateCardHandler(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	srv cards.CardService
}

func (c *controller) CreateCardHandler(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Phrase      string `json:"phrase"`
		Translation string `json:"translation"`
	}

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	err = c.srv.CreateCard(d.Phrase, d.Translation)
	if err != nil {
		msg := fmt.Sprintf("Unable to create a card: %s", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}

func NewController(srv cards.CardService) Controller {
	return &controller{srv: srv}
}
