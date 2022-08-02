package storage

import (
	"encoding/json"
	// "fmt"
	"io"
	"os"

	"example.com/recallcards/pkg/cards"
	// "github.com/rs/zerolog"
	// "github.com/rs/zerolog/log"
)


func readJsonFile[K comparable] (path string) (map[K]cards.Card, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0755)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	data := make(map[K]cards.Card)
	err = json.NewDecoder(f).Decode(&data)

	if err != nil && err != io.EOF {
		return nil, err
	}
	return data, nil
}

func writeJsonFile[K comparable] (data map[K]cards.Card, path string) error {
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
	st map[string]cards.Card
	filepath string
}

func (cr *mem) Insert(c cards.Card) error {
	cr.st[c.Phrase] = c
	return cr.Persist()
}

func (cr *mem) Persist() error {
	err := writeJsonFile(cr.st, cr.filepath)
	return err
}

func NewMemoryRepository(filepath string) (*mem, error) {
	data, err := readJsonFile[string](filepath)
	if err != nil {
		return nil, err
	}
	return &mem{st: data, filepath: filepath}, nil
}


