package file

import (
	"example.com/recallcards/pkg/cards"
	"path/filepath"
	"reflect"
)

type repository struct {
	cardStorage
	recallStorage
	dir string
}

func filenameFromStorage[C cardStorage | recallStorage]() string {
	f, ok := reflect.TypeOf((*C)(nil)).Elem().FieldByName("cache")
	if !ok {
		panic("Unable to get filename from storage type definition")
	}
	return f.Tag.Get("filename")
}

func NewRepository(dir string) (cards.Repository, error) {
	cardsFilePath := filepath.Join(dir, filenameFromStorage[cardStorage]())
	cardsList, err := readJsonFile[cards.Cards](cardsFilePath)
	if err != nil {
		return nil, err
	}

	recallsFilePath := filepath.Join(dir, filenameFromStorage[recallStorage]())
	recallsList, err := readJsonFile[cards.Recalls](recallsFilePath)
	if err != nil {
		return nil, err
	}

	rep := repository{
		cardStorage: cardStorage{
			cache:    cardsList,
			filepath: cardsFilePath,
		},
		recallStorage: recallStorage{
			cache:  recallsList,
			filepath: recallsFilePath,
		},
		dir: dir,
	}
	return &rep, nil
}
