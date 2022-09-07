package file

import (
	"errors"
	"fmt"

	"example.com/recallcards/pkg/cards"
)

type Cards []cards.Card

func (objects Cards) getNextID() cards.CardId {
	maxID := cards.CardId(0)
	for i := 0; i < len(objects); i++ {
		id := objects[i].ID
		if id > maxID {
			maxID = id
		}
	}
	maxID++
	return maxID
}

func (objects Cards) append(c cards.Card) Cards {
	c.ID = objects.getNextID()
	return append(objects, c)
}

type cardStorage struct {
	Cards   Cards   `json:"cards"`
	filepath string
}

func (s cardStorage) persist() error {
	return writeJsonFile(s, s.filepath)
}

func (s cardStorage) InsertCard(c cards.Card) error {
	_, err := s.findCardByPhrase(c.Phrase)
	if err == nil {
		return fmt.Errorf("Card with a phrase \"%s\" already exists", c.Phrase)
	}
	s.Cards = s.Cards.append(c)
	return s.persist()
}

func (s cardStorage) findCardByPhrase(phrase string) (cards.Card, error) {
	for i := range s.Cards {
		if s.Cards[i].Phrase == phrase {
			return s.Cards[i], nil
		}
	}
	return cards.Card{}, errors.New("Not found")
}

func (s cardStorage) CardById(cid cards.CardId) (cards.Card, error) {
	for i := range s.Cards {
		if s.Cards[i].ID == cid {
			return s.Cards[i], nil
		}
	}
	return cards.Card{}, errors.New("Not found")
}

func (s cardStorage) ListCardIds() ([]cards.CardId, error) {
	res := make([]cards.CardId, len(s.Cards))
	for i := range s.Cards {
		res[i] = s.Cards[i].ID
	}
	return res, nil
}