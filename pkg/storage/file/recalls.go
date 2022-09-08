package file

import (
	"example.com/recallcards/pkg/cards"
)

type recallStorage struct {
	Recalls  cards.Recalls `json:"recalls"`
	filepath string
}

func (s recallStorage) persist() error {
	return writeJsonFile(s.Recalls, s.filepath)
}

func (s *recallStorage) InsertRecallAttempt(r cards.RecallAttempt) error {
	s.Recalls = append(s.Recalls, r)
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
