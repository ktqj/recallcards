package cards

import (
	"math/rand"
	"time"
)

type CardService interface {
  CreateCard(phrase string, translation string) error
  // ListBuckets() []BucketId
  // GetRandomWeighted() Card
  RandomCard() (Card, error)
  // GetRandomByBucket(bid BucketId) Card
  RecordRecallAttempt(cid CardId, result bool) error
  CountRecallAttempts(cid CardId) int
}

type cardService struct {
  repo CardRepository  
}

func NewCardService(repo CardRepository) *cardService {
  return &cardService{repo: repo}
}

func (cs *cardService) CreateCard(phrase string, translation string) error {
  c := Card{
    Phrase: phrase,
    Translation: translation,
    CreatedAt: time.Now(),
    Bucket: DefaultBucket,
  }
  return cs.repo.InsertCard(c)
}

func (cs *cardService) RandomCard() (Card, error) {
  buckets, err := cs.repo.ListUsedBuckets()
  if err != nil {
    return Card{}, err
  }
  randomBucket := buckets[rand.Intn(len(buckets))]
  card, err := cs.repo.RandomCardByBucket(randomBucket)
  if err != nil {
    return Card{}, err
  }
  return card, nil
}

func (cs *cardService) RecordRecallAttempt(cid CardId, success bool) error {
  r := RecallAttempt{
    RecordedAt: time.Now(),
    Success: success,
    CardId: cid,
  }
  err := cs.repo.InsertRecallAttempt(r)
  if err != nil {
    return err
  }
  // TODO: count successfull attempts and move to another bucket if necessary
  return nil
}

func (cs *cardService) CountRecallAttempts(cid CardId) int {
  return cs.repo.CountRecallAttempts(cid)
}
