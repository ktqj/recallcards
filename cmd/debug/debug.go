package main

import (
	"fmt"

	"math/rand"
	"os"

	"example.com/recallcards/pkg/cards"
	"example.com/recallcards/pkg/storage/file"
	"github.com/rs/zerolog/log"
)

func initFileRepository() cards.CardRepository {
	memFilePath := os.Getenv("MEM_STORAGE_JSON_FILE_PATH")
	if memFilePath == "" {
		log.Fatal().Msgf("No file path provided under MEM_STORAGE_JSON_FILE_PATH var")
	}

	rep, err := file.NewRepository(memFilePath)
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
