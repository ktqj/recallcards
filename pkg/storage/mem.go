package storage

import (
	"encoding/json"
	"math/rand"

	// "fmt"
	"io"
	"os"

	"example.com/recallcards/pkg/cards"
	// "github.com/rs/zerolog"
	// "github.com/rs/zerolog/log"
)


func readJsonFile[K comparable] (path string) (map[K]cards.Card, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := make(map[K]cards.Card)
	err = json.NewDecoder(f).Decode(&data)

	if err != nil && err != io.EOF {
		return nil, err
	}
	return data, nil
}

func writeJsonFile[K comparable] (data map[K]cards.Card, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
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

type mem struct {
	st map[string]cards.Card
	filepath string
}

func (cr *mem) Insert(c cards.Card) error {
	cr.st[c.Phrase] = c
	return cr.persist()
}

func (cr *mem) ListUsedBuckets() ([]cards.BucketId, error) {
	buckets := make(map[cards.BucketId]struct{})
	for ph := range cr.st {
		c := cr.st[ph]
		buckets[c.Bucket] = struct{}{}
	}

	res := make([]cards.BucketId, len(buckets))
	i := 0
	for b := range buckets {
		res[i] = b
		i++
	}
	return res, nil
}

func (cr *mem) countByBucket(b cards.BucketId) int {
	count := 0
	for ph := range cr.st {
		c := cr.st[ph]
		if c.Bucket == b {
			count++
		}
	}
	return count
}

func (cr *mem) RandomByBucket(b cards.BucketId) (cards.Card, error) {
	count := cr.countByBucket(b)
	picked := rand.Intn(count)

	i := 0
	for ph := range cr.st {
		c := cr.st[ph]
		if c.Bucket != b {
			continue
		}
		if i != picked {
			i++
			continue
		}
		return c, nil
	}
	return cards.Card{}, nil
}

func (cr *mem) persist() error {
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


