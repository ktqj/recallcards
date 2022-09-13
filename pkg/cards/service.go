package cards

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"

	// "runtime"
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
	RandomCardGenerator(ctx context.Context) (<-chan Card, error)
	FilteredRandomCardGenerator(ctx context.Context) (<-chan Card, error)
	RecordRecallAttempt(cid CardId, result bool) error
	RecallSummary(cid CardId) RecallSummary
}

type cardService struct {
	repo Repository
}

func NewCardService(repo Repository) CardService {
	return &cardService{repo: repo}
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
	var card Card

	buckets, err := cs.repo.ListUsedBuckets()
	if err != nil {
		return card, err
	}

	randomBucket := buckets[rand.Intn(len(buckets))]
	card, err = cs.repo.RandomCardByBucket(randomBucket)
	if err != nil {
		return card, err
	}

	return card, nil
}

func (cs *cardService) RandomCardGenerator(ctx context.Context) (<-chan Card, error) {
	cids, err := cs.repo.ListCardIds()
	if err != nil {
		return nil, err
	}

	rand.Shuffle(len(cids), func(i, j int) { cids[i], cids[j] = cids[j], cids[i] })

	gen := make(chan Card)
	go func() {
		defer close(gen)
		for _, cid := range cids {
			card, err := cs.repo.CardById(cid)
			if err != nil {
				continue
			}
			select {
			case gen <- card:
			case <-ctx.Done():
				log.Debug().Msg("Closing random card generator")
				return
			}
		}
	}()

	return gen, nil
}

func (cs *cardService) shouldCardBeDisplayed(cid CardId) bool {
	confidence := cs.RecallSummary(cid).EstimateConfidence()

	if confidence <= 50 {
		return true
	}

	w := (100 - confidence) / 5
	// TODO: should be configurable
	bias := 2
	return rand.Intn(w+bias) < w
}

func (cs *cardService) filterCardsStreamByConfidence(ctx context.Context, in <-chan Card) <-chan Card {
	out := make(chan Card)
	go func() {
		defer close(out)
		defer log.Debug().Msg("Closing a card filter")

		log.Debug().Msg("Starting a card filter")
		for c := range in {
			if !cs.shouldCardBeDisplayed(c.ID) {
				continue
			}
			select {
			case out <- c:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func (cs *cardService) FilteredRandomCardGenerator(ctx context.Context) (<-chan Card, error) {
	stream, err := cs.RandomCardGenerator(ctx)
	if err != nil {
		return nil, fmt.Errorf("Could not initialize random card generator: [%w]", err)
	}

	// fan out cards to filters
	workers := make([]<-chan Card, runtime.NumCPU())
	for i := 0; i < len(workers); i++ {
		workers[i] = cs.filterCardsStreamByConfidence(ctx, stream)
	}

	// fan in filtered cards
	var wg sync.WaitGroup
	res := make(chan Card)

	multiplex := func(ctx context.Context, in <-chan Card) {
		defer wg.Done()
		defer log.Debug().Msg("Closing a multiplexer")

		log.Debug().Msg("Starting a multiplexer")
		for c := range in {
			select {
			case <-ctx.Done():
				return
			case res <- c:
			}
		}
	}

	wg.Add(len(workers))
	for _, w := range workers {
		go multiplex(ctx, w)
	}

	go func() {
		defer close(res)
		wg.Wait()
		log.Debug().Msg("Closing filtered cards generator")
	}()

	return res, nil
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
	return nil
}

func (cs *cardService) RecallSummary(cid CardId) RecallSummary {
	return cs.repo.RecallSummary(cid)
}
