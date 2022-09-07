package file

import (
	"errors"
	"fmt"
	"math/rand"

	"example.com/recallcards/pkg/cards"
)

type cardStorage struct {
	Cards    []cards.Card `json:"cards"`
	filepath string
}

func (s cardStorage) persist() error {
	return writeJsonFile(s, s.filepath)
}

func (s *cardStorage) InsertCard(c cards.Card) error {
	_, err := s.findCardByPhrase(c.Phrase)
	if err == nil {
		return fmt.Errorf("Card with a phrase \"%s\" already exists", c.Phrase)
	}

	c.ID = s.getNextID()
	s.Cards = append(s.Cards, c)
	return s.persist()
}

func (s cardStorage) getNextID() cards.CardId {
	maxID := cards.CardId(0)
	for i := 0; i < len(s.Cards); i++ {
		id := s.Cards[i].ID
		if id > maxID {
			maxID = id
		}
	}
	maxID++
	return maxID
}

func (s cardStorage) findCardByPhrase(phrase string) (cards.Card, error) {
	for i := range s.Cards {
		if s.Cards[i].Phrase == phrase {
			return s.Cards[i], nil
		}
	}
	return cards.Card{}, errors.New("Not found")
}

func (s cardStorage) CardById(cid cards.CardId) (cards.Card, error) {
	for i := range s.Cards {
		if s.Cards[i].ID == cid {
			return s.Cards[i], nil
		}
	}
	return cards.Card{}, errors.New("Not found")
}

func (s cardStorage) ListCardIds() ([]cards.CardId, error) {
	res := make([]cards.CardId, len(s.Cards))
	for i := range s.Cards {
		res[i] = s.Cards[i].ID
	}
	return res, nil
}

func (s cardStorage) ListUsedBuckets() ([]cards.BucketId, error) {
	buckets := make(map[cards.BucketId]struct{})
	for i := 0; i < len(s.Cards); i++ {
		buckets[s.Cards[i].Bucket] = struct{}{}
	}

	res := make([]cards.BucketId, len(buckets))
	i := 0
	for b := range buckets {
		res[i] = b
		i++
	}
	return res, nil
}

func (s cardStorage) RandomCardByBucket(b cards.BucketId) (cards.Card, error) {
	count := s.countByBucket(b)
	picked := rand.Intn(count)

	j := 0
	for i := 0; i < len(s.Cards); i++ {
		if s.Cards[i].Bucket != b {
			continue
		}
		if j != picked {
			j++
			continue
		}
		return s.Cards[i], nil
	}
	return cards.Card{}, nil
}

func (s cardStorage) countByBucket(b cards.BucketId) int {
	count := 0
	for i := 0; i < len(s.Cards); i++ {
		if s.Cards[i].Bucket == b {
			count++
		}
	}
	return count
}
