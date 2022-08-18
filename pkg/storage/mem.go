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


func readJsonFile(path string) (storage, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return storage{}, err
	}
	defer f.Close()

	var data storage
	err = json.NewDecoder(f).Decode(&data)

	if err != nil && err != io.EOF {
		return storage{}, err
	}
	return data, nil
}

func writeJsonFile(data storage, path string) error {
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

type storage struct {
	Cards []cards.Card `json:"cards"`
	Recalls []cards.RecallAttempt `json:"recalls"`
}

type repository struct {
	storage
	filepath string
}

func (rep *repository) Insert(c cards.Card) error {
	rep.storage.Cards = append(rep.storage.Cards, c)
	return rep.persist()
}

func (rep *repository) ListUsedBuckets() ([]cards.BucketId, error) {
	buckets := make(map[cards.BucketId]struct{})
	for i := 0; i < len(rep.storage.Cards); i++ {
		buckets[rep.storage.Cards[i].Bucket] = struct{}{}
	}

	res := make([]cards.BucketId, len(buckets))
	i := 0
	for b := range buckets {
		res[i] = b
		i++
	}
	return res, nil
}

func (rep *repository) countByBucket(b cards.BucketId) int {
	count := 0
	for i := 0; i < len(rep.storage.Cards); i++ {
		if rep.storage.Cards[i].Bucket == b {
			count++
		}
	}
	return count
}

func (rep *repository) RandomByBucket(b cards.BucketId) (cards.Card, error) {
	count := rep.countByBucket(b)
	picked := rand.Intn(count)

	j := 0
	for i := 0; i < len(rep.storage.Cards); i++ {
		if rep.storage.Cards[i].Bucket != b {
			continue
		}
		if j != picked {
			j++
			continue
		}
		return rep.storage.Cards[i], nil
	}
	return cards.Card{}, nil
}

func (rep *repository) persist() error {
	err := writeJsonFile(rep.storage, rep.filepath)
	return err
}

func NewMemoryRepository(filepath string) (*repository, error) {
	data, err := readJsonFile(filepath)
	if err != nil {
		return nil, err
	}
	rep := repository{
		storage: data,
		filepath: filepath,
	}
	return &rep, nil
}


