package file

import (
	"example.com/recallcards/pkg/cards"
)

type recallStorage struct {
	cache    cards.Recalls `filename:"recalls.json"`
	filepath string
}

func (s *recallStorage) InsertRecallAttempt(r cards.RecallAttempt) error {
	// this section is not thread-safe
	s.cache = append(s.cache, r)
	return persistCollection(s.cache, s.filepath)
}

func (s recallStorage) CountRecallAttempts(cid cards.CardId) int {
	count := 0
	for i := 0; i < len(s.cache); i++ {
		r := s.cache[i]
		if r.CardId == cid && r.Success {
			count++
		}
	}
	return count
}
