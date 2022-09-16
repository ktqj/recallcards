package file

import (
	"path/filepath"
	"reflect"

	"example.com/recallcards/pkg/cards"
)

type repository struct {
	cardStorage
	recallStorage
	dir string
}

func filenameFromStorage[C cardStorage | recallStorage]() string {
	f, ok := reflect.TypeOf((*C)(nil)).Elem().FieldByName("objects")
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
			objects:  cardsList,
			filepath: cardsFilePath,
		},
		recallStorage: recallStorage{
			objects:  recallsList,
			filepath: recallsFilePath,
		},
		dir: dir,
	}
	return &rep, nil
}
