package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os/signal"
	"syscall"

	// "runtime"
	"time"

	"os"

	"example.com/recallcards/pkg/cards"
	"example.com/recallcards/pkg/storage/file"
	"github.com/rs/zerolog/log"
)

func readline(prompt string) (string, error) {
	fmt.Fprintf(os.Stdout, prompt)
	r := bufio.NewReader(os.Stdout)
	input, err := r.ReadString('\n')
	return input[:len(input)-1], err
}

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
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	rand.Seed(int64(time.Now().Nanosecond()))
	// rand.Seed(1)

	repository := initFileRepository()
	cardService := cards.NewCardService(repository)

	generator, done := cardService.RandomCardGenerator()

	go func() {
		select {
		case <-sigChannel:
			fmt.Println("\nsigterm received")
			done()
			time.Sleep(100 * time.Microsecond)
			os.Exit(0)
		}
	}()

	i := 1
	for card := range generator {

		fmt.Fprintf(os.Stdout, "Recall #%d, card [ID: %d, attempts: %d]\n", i, card.ID, cardService.CountRecallAttempts(card.ID))

		readline(card.Translation)
		fmt.Print(card.Phrase + "\n")
		for {
			answer, _ := readline("Got it right? [y/n]\n")
			if answer == "y" {
				cardService.RecordRecallAttempt(card.ID, true)
				break
			} else if answer == "n" {
				cardService.RecordRecallAttempt(card.ID, false)
				break
			}
		}
		i++
	}
}
