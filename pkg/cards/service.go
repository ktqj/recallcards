package cards

import (
	"math/rand"
	"time"
)

type CardService interface {
  Create(phrase string, translation string) error
  // ListBuckets() []BucketId
  // GetRandomWeighted() Card
  // Random() (Card, error)
  // GetRandomByBucket(bid BucketId) Card
  // RecordRecallAttempt(cid CardId, result bool) error
}

type cardService struct {
  repo CardRepository  
}

func NewCardService(repo CardRepository) *cardService {
  return &cardService{repo: repo}
}

func (cs *cardService) Create(phrase string, translation string) error {
  c := Card{
    Phrase: phrase,
    Translation: translation,
    Created_at: time.Now(),
    RecallAttempts: []RecallAttempt{},
    Bucket: DefaultBucket,
  }
  return cs.repo.Insert(c)
}

func (cs *cardService) Random() (Card, error) {
  buckets, err := cs.repo.ListUsedBuckets()
  if err != nil {
    return Card{}, err
  }
  randomBucket := buckets[rand.Intn(len(buckets))]
  card, err := cs.repo.RandomByBucket(randomBucket)
  if err != nil {
    return Card{}, err
  }
  return card, nil
}
