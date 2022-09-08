package file

import (
	"example.com/recallcards/pkg/cards"
)

type recallStorage struct {
	objects  cards.Recalls `filename:"recalls.json"`
	filepath string
}

func (s *recallStorage) InsertRecallAttempt(r cards.RecallAttempt) error {
	// this section is not thread-safe
	s.objects = append(s.objects, r)
	return persistCollection(s.objects, s.filepath)
}

func (s recallStorage) CountRecallAttempts(cid cards.CardId) cards.RecallSummary {
	var res cards.RecallSummary
	for i := 0; i < len(s.objects); i++ {
		r := s.objects[i]
		if r.CardId != cid {
			continue
		}
		if r.Success {
			res.Ok++
		} else {
			res.Fail++
		}
	}
	return res
}
