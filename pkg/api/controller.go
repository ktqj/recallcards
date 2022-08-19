package api

import (
	"encoding/json"
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

	err = c.srv.CreateCard(d.Phrase, d.Translation)
	if err != nil {
		http.Error(w, "Unable to create a card", http.StatusInternalServerError)
		return
	}
}

func NewController(srv cards.CardService) *controller {
	return &controller{srv: srv}
}