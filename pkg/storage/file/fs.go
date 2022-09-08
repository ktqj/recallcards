package file

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"

	"example.com/recallcards/pkg/cards"
)

type coll interface {
	cards.Cards | cards.Recalls
}

func readJsonFile[S coll](p string) (S, error) {
	var data S

	f, err := os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0755)
	if errors.Is(err, fs.ErrNotExist) {
		pwd, _ := os.Getwd()
		absPath := path.Join(pwd, p)
		f, err = os.OpenFile(absPath, os.O_RDONLY|os.O_CREATE, 0755)
	}
	if err != nil {
		return data, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&data)

	if err != nil && err != io.EOF {
		return data, err
	}
	return data, nil
}

func writeJsonFile[S coll](data S, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func persistCollection[C coll](objects C, filepath string) error {
	return writeJsonFile(objects, filepath)
}
