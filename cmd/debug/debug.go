package main

import (
	"fmt"
	"math/rand"
	"os"

	"example.com/recallcards/pkg/cards"
	"example.com/recallcards/pkg/storage/file"
	"github.com/rs/zerolog/log"

	"net/http"
	_ "net/http/pprof"
)

func initFileRepository() cards.Repository {
	dir := os.Getenv("JSON_STORAGE_DIR")
	if dir == "" {
		log.Fatal().Msgf("No file path provided under JSON_STORAGE_DIR var")
	}

	rep, err := file.NewRepository(dir)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not initialize repository")
	}
	return rep
}

func shouldCardBeDisplayed(confidence int) bool {
	if confidence <= 50 {
		return true
	}

	w := (100 - confidence) / 5
	bias := 2
	n := rand.Intn(w + bias)
	return n < w
}

func main() {
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			log.Fatal().Err(err).Msg("httpServer exited")
		}
	}()

	repository := initFileRepository()
	cardService := cards.NewCardService(repository)

	generator, _ := cardService.RandomCardGenerator()

	i := 0
	j := 0
	limit := 65
	for card := range generator {

		recalls := cardService.CountRecallAttempts(card.ID)
		confidence := cardService.EstimateCardConfidence(recalls)
		if confidence >= limit {
			i++
		}

		if shouldCardBeDisplayed(confidence) {
			j++
		}
	}
	ids, _ := repository.ListCardIds()
	fmt.Printf("%d/%d cards are over %d%% of confidence\n", i, len(ids), limit)
	fmt.Printf("%d/%d cards would be shown\n", j, len(ids))
}
