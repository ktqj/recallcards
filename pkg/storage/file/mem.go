package file

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

type Recalls []cards.RecallAttempt

func (objects Recalls) append(r cards.RecallAttempt) Recalls {
	return append(objects, r)
}

type storage struct {
	Cards Cards `json:"cards"`
	Recalls Recalls `json:"recalls"`
}

type repository struct {
	storage
	filepath string
}

func (rep *repository) InsertCard(c cards.Card) error {
	rep.storage.Cards = rep.storage.Cards.append(c)
	return rep.persist()
}

func (rep *repository) InsertRecallAttempt(r cards.RecallAttempt) error {
	rep.storage.Recalls = rep.storage.Recalls.append(r)
	return rep.persist()
}

func (rep *repository) CountRecallAttempts(cid cards.CardId) int {
	count := 0
	for i := 0; i < len(rep.storage.Recalls); i++ {
		if rep.storage.Recalls[i].Success {
			count++
		}
	}
	return count
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

func (rep *repository) RandomCardByBucket(b cards.BucketId) (cards.Card, error) {
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

func (rep *repository) countByBucket(b cards.BucketId) int {
	count := 0
	for i := 0; i < len(rep.storage.Cards); i++ {
		if rep.storage.Cards[i].Bucket == b {
			count++
		}
	}
	return count
}


func (rep *repository) persist() error {
	err := writeJsonFile(rep.storage, rep.filepath)
	return err
}

func NewRepository(filepath string) (*repository, error) {
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


