package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"example.com/recallcards/pkg/cards"
)

type Controller interface {
	IndexHandler(w http.ResponseWriter, r *http.Request)
	CreateCardForm(w http.ResponseWriter, r *http.Request)
	CreateCardJson(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	srv cards.CardService
}

func (c *controller) CreateCardForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	err := c.srv.CreateCard(r.PostForm.Get("phrase"), r.PostForm.Get("translation"))
	if err != nil {
		msg := fmt.Sprintf("Unable to create a card: %s", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func (c *controller) CreateCardJson(w http.ResponseWriter, r *http.Request) {
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

func (c *controller) IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pkg/web/templates/create_card.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, &struct{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewController(srv cards.CardService) Controller {
	return &controller{srv: srv}
}
