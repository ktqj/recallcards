package file

import (
	"math/rand"
	"example.com/recallcards/pkg/cards"
)

type repository struct {
	cardStorage
	recallStorage
}

func (rep *repository) ListUsedBuckets() ([]cards.BucketId, error) {
	buckets := make(map[cards.BucketId]struct{})
	for i := 0; i < len(rep.Cards); i++ {
		buckets[rep.Cards[i].Bucket] = struct{}{}
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
	for i := 0; i < len(rep.Cards); i++ {
		if rep.Cards[i].Bucket != b {
			continue
		}
		if j != picked {
			j++
			continue
		}
		return rep.Cards[i], nil
	}
	return cards.Card{}, nil
}

func (rep *repository) countByBucket(b cards.BucketId) int {
	count := 0
	for i := 0; i < len(rep.Cards); i++ {
		if rep.Cards[i].Bucket == b {
			count++
		}
	}
	return count
}

func NewRepository(dir string) (cards.Repository, error) {
	cards, err := readJsonFile[cardStorage](dir + "cards.json")
	cards.filepath = dir + "cards.json"
	if err != nil {
		return nil, err
	}

	recalls, err := readJsonFile[recallStorage](dir + "recalls.json")
	recalls.filepath = dir + "recalls.json"
	if err != nil {
		return nil, err
	}

	rep := repository{
		cardStorage:  cards,
		recallStorage: recalls,
	}
	return &rep, nil
}
