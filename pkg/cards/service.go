package cards

import "time"

type CardService interface {
  Create(phrase string, translation string) error
  // ListBuckets() []BucketId
  // GetRandomWeighted() Card
  // GetRandom() Card
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
