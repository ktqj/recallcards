package file

import (
	"errors"
	"fmt"
	"math/rand"

	"example.com/recallcards/pkg/cards"
)

type cardStorage struct {
	objects  cards.Cards `filename:"cards.json"`
	filepath string
}

func (s *cardStorage) InsertCard(c cards.Card) error {
	_, err := s.findCardByPhrase(c.Phrase)
	if err == nil {
		return fmt.Errorf("Card with a phrase \"%s\" already exists", c.Phrase)
	}

	// this section is not thread-safe
	c.ID = s.getNextID()
	s.objects = append(s.objects, c)
	return persistCollection(s.objects, s.filepath)
}

func (s cardStorage) getNextID() cards.CardId {
	// not thread-safe
	maxID := cards.CardId(0)
	for i := 0; i < len(s.objects); i++ {
		id := s.objects[i].ID
		if id > maxID {
			maxID = id
		}
	}
	maxID++
	return maxID
}

func (s cardStorage) findCardByPhrase(phrase string) (cards.Card, error) {
	for i := range s.objects {
		if s.objects[i].Phrase == phrase {
			return s.objects[i], nil
		}
	}
	return cards.Card{}, errors.New("Not found")
}

func (s cardStorage) CardById(cid cards.CardId) (cards.Card, error) {
	for i := range s.objects {
		if s.objects[i].ID == cid {
			return s.objects[i], nil
		}
	}
	return cards.Card{}, errors.New("Not found")
}

func (s cardStorage) ListCardIds() ([]cards.CardId, error) {
	res := make([]cards.CardId, len(s.objects))
	for i := range s.objects {
		res[i] = s.objects[i].ID
	}
	return res, nil
}

func (s cardStorage) ListUsedBuckets() ([]cards.BucketId, error) {
	buckets := make(map[cards.BucketId]struct{})
	for i := 0; i < len(s.objects); i++ {
		buckets[s.objects[i].Bucket] = struct{}{}
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
	for i := 0; i < len(s.objects); i++ {
		if s.objects[i].Bucket != b {
			continue
		}
		if j != picked {
			j++
			continue
		}
		return s.objects[i], nil
	}
	return cards.Card{}, nil
}

func (s cardStorage) countByBucket(b cards.BucketId) int {
	count := 0
	for i := 0; i < len(s.objects); i++ {
		if s.objects[i].Bucket == b {
			count++
		}
	}
	return count
}
