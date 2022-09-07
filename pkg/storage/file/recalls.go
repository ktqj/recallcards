package file

import (
	"example.com/recallcards/pkg/cards"
)

type Recalls []cards.RecallAttempt

func (objects Recalls) append(r cards.RecallAttempt) Recalls {
	return append(objects, r)
}

type recallStorage struct {
	Recalls Recalls `json:"recalls"`
	filepath string
}

func (s recallStorage) persist() error {
	return writeJsonFile(s, s.filepath)
}

func (s *recallStorage) InsertRecallAttempt(r cards.RecallAttempt) error {
	s.Recalls = s.Recalls.append(r)
	return s.persist()
}

func (s recallStorage) CountRecallAttempts(cid cards.CardId) int {
	count := 0
	for i := 0; i < len(s.Recalls); i++ {
		r := s.Recalls[i]
		if r.CardId == cid && r.Success {
			count++
		}
	}
	return count
}