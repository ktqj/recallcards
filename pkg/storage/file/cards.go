package file

import (
	"errors"
	"fmt"
	"math/rand"

	"example.com/recallcards/pkg/cards"
)

type cardStorage struct {
	cache cards.Cards `filename:"cards.json"`
	filepath string
}

func (s *cardStorage) InsertCard(c cards.Card) error {
	_, err := s.findCardByPhrase(c.Phrase)
	if err == nil {
		return fmt.Errorf("Card with a phrase \"%s\" already exists", c.Phrase)
	}

	c.ID = s.getNextID()
	s.cache = append(s.cache, c)
	return persistCollection(s.cache, s.filepath)
}

func (s cardStorage) getNextID() cards.CardId {
	maxID := cards.CardId(0)
	for i := 0; i < len(s.cache); i++ {
		id := s.cache[i].ID
		if id > maxID {
			maxID = id
		}
	}
	maxID++
	return maxID
}

func (s cardStorage) findCardByPhrase(phrase string) (cards.Card, error) {
	for i := range s.cache {
		if s.cache[i].Phrase == phrase {
			return s.cache[i], nil
		}
	}
	return cards.Card{}, errors.New("Not found")
}

func (s cardStorage) CardById(cid cards.CardId) (cards.Card, error) {
	for i := range s.cache {
		if s.cache[i].ID == cid {
			return s.cache[i], nil
		}
	}
	return cards.Card{}, errors.New("Not found")
}

func (s cardStorage) ListCardIds() ([]cards.CardId, error) {
	res := make([]cards.CardId, len(s.cache))
	for i := range s.cache {
		res[i] = s.cache[i].ID
	}
	return res, nil
}

func (s cardStorage) ListUsedBuckets() ([]cards.BucketId, error) {
	buckets := make(map[cards.BucketId]struct{})
	for i := 0; i < len(s.cache); i++ {
		buckets[s.cache[i].Bucket] = struct{}{}
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
	for i := 0; i < len(s.cache); i++ {
		if s.cache[i].Bucket != b {
			continue
		}
		if j != picked {
			j++
			continue
		}
		return s.cache[i], nil
	}
	return cards.Card{}, nil
}

func (s cardStorage) countByBucket(b cards.BucketId) int {
	count := 0
	for i := 0; i < len(s.cache); i++ {
		if s.cache[i].Bucket == b {
			count++
		}
	}
	return count
}
