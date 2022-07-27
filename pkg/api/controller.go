package api

import (
	"encoding/json"
	"net/http"

	"example.com/recallcards/pkg/cards"
)

type Controller struct {
	srv cards.CardService
}

func (c *Controller) CreateCard(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Phrase string `json:"phrase"`
		Translation string `json:"translation"`
	}
	
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	if d.Phrase == "" {
		http.Error(w, "No phrase is provided", http.StatusBadRequest)
		return
	}

	if d.Translation == "" {
		http.Error(w, "No translation is provided", http.StatusBadRequest)
		return
	}

	err = c.srv.Create(d.Phrase, d.Translation)
	if err != nil {
		http.Error(w, "Unable to create a card", http.StatusInternalServerError)
		return
	}
}

func NewController(srv cards.CardService) *Controller {
	return &Controller{srv: srv}
}