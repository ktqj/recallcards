package main

import (
	"bufio"
	"fmt"
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
	fmt.Fprint(os.Stdout, prompt)
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

	repository := initFileRepository()
	cardService := cards.NewCardService(repository)

	generator, done := cardService.RandomCardGenerator()

	go func() {
		<-sigChannel
		fmt.Println("\nsigterm received")
		done()
		time.Sleep(100 * time.Microsecond)
		os.Exit(0)
	}()

	i := 1
	for card := range generator {

		fmt.Fprintf(os.Stdout, "Recall #%d, card [ID: %d, attempts: %d]\n", i, card.ID, cardService.CountRecallAttempts(card.ID))

		_, err := readline(card.Translation)
		if err != nil {
			panic(err)
		}

		fmt.Print(card.Phrase + "\n")
		for {
			answer, _ := readline("Got it right? [y/n]\n")
			if answer == "y" {
				err := cardService.RecordRecallAttempt(card.ID, true)
				if err != nil {
					panic(err)
				}
				break
			} else if answer == "n" {
				err := cardService.RecordRecallAttempt(card.ID, false)
				if err != nil {
					panic(err)
				}
				break
			}
		}
		i++
	}
}
