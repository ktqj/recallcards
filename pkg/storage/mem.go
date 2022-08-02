package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"example.com/recallcards/pkg/cards"
	// "github.com/rs/zerolog"
	// "github.com/rs/zerolog/log"
)

type CardsByPhrase map[string]cards.Card


func loadData(path string) (CardsByPhrase, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0755)
	defer f.Close()

	if err != nil {
		return nil, fmt.Errorf("Cannot create data.json: %v", err)
	}

	data := make(CardsByPhrase)
	err = json.NewDecoder(f).Decode(&data)

	if err != nil {
		return nil, fmt.Errorf("Error decoding data.json: %v", err)
	}
	return data, nil
}

func saveData(data CardsByPhrase, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()

	if err != nil {
		return err
	}

	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		return err
	}
	return nil
}

type mem struct {
	st CardsByPhrase
	filepath string
}

func (cr *mem) Insert(c cards.Card) error {
	cr.st[c.Phrase] = c
	err := cr.Persist()
	return err
}

func (cr *mem) Persist() error {
	err := saveData(cr.st, cr.filepath)
	return err
}

func NewMemoryRepository(filepath string) (*mem, error) {
	data, err := loadData(filepath)
	if err != nil {
		return nil, err
	}
	return &mem{st: data, filepath: filepath}, nil
}


