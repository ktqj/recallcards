package main

import (
	"fmt"
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
	limit := 65
	for card := range generator {

		recalls := cardService.CountRecallAttempts(card.ID)
		confidence := cardService.EstimateCardConfidence(card.ID, recalls)
		if confidence >= limit {
			i++
		}
	}
	ids, _ := repository.ListCardIds()
	fmt.Printf("%d/%d cards are over %d%% of confidence", i, len(ids), limit)
}
