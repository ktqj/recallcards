package main

import (
	"bufio"
	"fmt"
	"math/rand"

	"os"

	"example.com/recallcards/pkg/cards"
	"example.com/recallcards/pkg/storage"
	"github.com/rs/zerolog/log"
)

func readline(prompt string) (string, error) {
	fmt.Fprintf(os.Stdout, prompt)
	r := bufio.NewReader(os.Stdout)
	input, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return input, nil
}

func initInMemRepository() cards.CardRepository {
	memFilePath := os.Getenv("MEM_STORAGE_JSON_FILE_PATH")
	if memFilePath == "" {
		log.Fatal().Msgf("No file path provided under MEM_STORAGE_JSON_FILE_PATH var")
	}

	rep, err := storage.NewMemoryRepository(memFilePath)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not initialize repository")
	}
	return rep
}

func main() {
	repository := initInMemRepository()
	cardService := cards.NewCardService(repository)

	i := 1
	for {
		fmt.Fprintf(os.Stdout, "Recall #%d\n", i)

		card, _ := cardService.Random()
		if rand.Intn(2) == 0 {
			readline(card.Phrase)
			readline(card.Translation)
		} else {
			readline(card.Translation)
			readline(card.Phrase)
		}
		i++
	}
}