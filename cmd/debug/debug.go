package main

import (
	"fmt"

	"math/rand"
	"os"

	"example.com/recallcards/pkg/cards"
	"example.com/recallcards/pkg/storage/file"
	"github.com/rs/zerolog/log"
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
	rep := initFileRepository()
	buckets, _ := rep.ListUsedBuckets()
	fmt.Printf("%v\n", buckets)

	card, _ := rep.RandomCardByBucket(buckets[0])
	if rand.Intn(2) == 0 {
		fmt.Printf("%v\n", card.Phrase)
	} else {
		fmt.Printf("%v\n", card.Translation)
	}
}
