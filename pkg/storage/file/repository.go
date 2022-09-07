package file

import (
	"example.com/recallcards/pkg/cards"
)

type repository struct {
	cardStorage
	recallStorage
}

func NewRepository(dir string) (cards.Repository, error) {
	cards, err := readJsonFile[cardStorage](dir + "cards.json")
	cards.filepath = dir + "cards.json"
	if err != nil {
		return nil, err
	}

	recalls, err := readJsonFile[recallStorage](dir + "recalls.json")
	recalls.filepath = dir + "recalls.json"
	if err != nil {
		return nil, err
	}

	rep := repository{
		cardStorage:   cards,
		recallStorage: recalls,
	}
	return &rep, nil
}
