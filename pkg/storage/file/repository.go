package file

import (
	"example.com/recallcards/pkg/cards"
)

type repository struct {
	cardStorage
	recallStorage
}

func NewRepository(dir string) (cards.Repository, error) {
	cardsFilePath := dir + "cards.json"
	cardsList, err := readJsonFile[cards.Cards](cardsFilePath)
	if err != nil {
		return nil, err
	}

	recallsFilePath := dir + "recalls.json"
	recallsList, err := readJsonFile[cards.Recalls](recallsFilePath)
	if err != nil {
		return nil, err
	}

	rep := repository{
		cardStorage:   cardStorage{
			Cards: cardsList,
			filepath: cardsFilePath,
		},
		recallStorage: recallStorage{
			Recalls: recallsList,
			filepath: recallsFilePath,
		},
	}
	return &rep, nil
}
