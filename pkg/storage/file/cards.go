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
	for _, c := range s.objects {
		if c.ID > maxID {
			maxID = c.ID
		}
	}
	maxID++
	return maxID
}

func (s cardStorage) findCardByPhrase(phrase string) (cards.Card, error) {
	for _, c := range s.objects {
		if c.Phrase == phrase {
			return c, nil
		}
	}
	return cards.Card{}, errors.New("Not found")
}

func (s cardStorage) CardById(cid cards.CardId) (cards.Card, error) {
	for _, c := range s.objects {
		if c.ID == cid {
			return c, nil
		}
	}
	return cards.Card{}, errors.New("Not found")
}

func (s cardStorage) ListCardIds() ([]cards.CardId, error) {
	res := make([]cards.CardId, len(s.objects))
	for i, c := range s.objects {
		res[i] = c.ID
	}
	return res, nil
}

func (s cardStorage) ListUsedBuckets() ([]cards.BucketId, error) {
	buckets := make(map[cards.BucketId]struct{})
	for _, c := range s.objects {
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

func (s cardStorage) RandomCardByBucket(b cards.BucketId) (cards.Card, error) {
	count := s.countByBucket(b)
	picked := rand.Intn(count)

	j := 0
	for _, c := range s.objects {
		if c.Bucket != b {
			continue
		}
		if j != picked {
			j++
			continue
		}
		return c, nil
	}
	return cards.Card{}, nil
}

func (s cardStorage) countByBucket(b cards.BucketId) int {
	count := 0
	for _, c := range s.objects {
		if c.Bucket == b {
			count++
		}
	}
	return count
}
