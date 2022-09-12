package cards

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

type CardService interface {
	CreateCard(phrase string, translation string) error
	RandomCard() (Card, error)
	RandomCardGenerator() (<-chan Card, func())
	RecordRecallAttempt(cid CardId, result bool) error
	CountRecallAttempts(cid CardId) RecallSummary
	EstimateCardConfidence(cid CardId, recalls RecallSummary) int
}

type cardService struct {
	repo Repository
}

func NewCardService(repo Repository) CardService {
	return &cardService{repo: repo}
}

// EstimateCardConfidence estimates confidence in knowing card's information
// based on its recalls history. Returns a percentage value, between 0 and 100.
func (cs *cardService) EstimateCardConfidence(cid CardId, recalls RecallSummary) int {
	// card was never recalled correctly yet
	if recalls.Ok == 0 {
		return 0
	}

	// weighted difference between correct and failed recalls
	failWeight := 2
	diff := recalls.Ok - failWeight * recalls.Fail

	confidentDifference := 10
	if diff >= confidentDifference {
		return 100
	}

	inconfidentDifference := -10
	if diff <= inconfidentDifference {
		return 0
	}

	// inconfidentDifference < diff < confidentDifference
	base := 50.0
	if diff >= 0 {
		progress := (100.0 - base) / float64(confidentDifference) * float64(diff)
		return int(base + progress)
	}

	regress := base / float64(inconfidentDifference) * float64(diff)
	return int(base - regress)
}

func (cs *cardService) CreateCard(phrase string, translation string) error {
	phrase = strings.TrimSpace(phrase)
	if phrase == "" {
		return fmt.Errorf("No phrase is provided")
	}

	translation = strings.TrimSpace(translation)
	if translation == "" {
		return fmt.Errorf("No translation is provided")
	}

	c := Card{
		Phrase:      phrase,
		Translation: translation,
		CreatedAt:   time.Now(),
		Bucket:      DefaultBucket,
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

func (cs *cardService) RandomCardGenerator() (<-chan Card, func()) {
	cids, err := cs.repo.ListCardIds()
	if err != nil {
		return nil, nil
	}

	rand.Shuffle(len(cids), func(i, j int) { cids[i], cids[j] = cids[j], cids[i] })

	g := make(chan Card)
	done := make(chan struct{})

	go func() {
		defer close(g)
		for _, cid := range cids {
			card, err := cs.repo.CardById(cid)
			if err != nil {
				continue
			}
			select {
			case g <- card:
			case <-done:
				log.Debug().Msg("Closing random card generator")
				return
			}
		}
	}()

	return g, func() { close(done) }
}

func (cs *cardService) RecordRecallAttempt(cid CardId, success bool) error {
	r := RecallAttempt{
		RecordedAt: time.Now(),
		Success:    success,
		CardId:     cid,
	}
	err := cs.repo.InsertRecallAttempt(r)
	if err != nil {
		return err
	}
	// TODO: count successfull attempts and move to another bucket if necessary
	return nil
}

func (cs *cardService) CountRecallAttempts(cid CardId) RecallSummary {
	return cs.repo.CountRecallAttempts(cid)
}
