package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"example.com/recallcards/pkg/cards"
	"example.com/recallcards/pkg/storage/file"
	"github.com/rs/zerolog/log"

	"net/http"
	_ "net/http/pprof"
)

func readlineAfter(prompt string) (string, error) {
	fmt.Fprint(os.Stdout, prompt)

	line, err := bufio.NewReader(os.Stdout).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

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

func shouldBeShown(confidence int) bool {
	if confidence <= 50 {
		return true
	}

	w := (100 - confidence) / 5
	bias := 2
	return rand.Intn(w+bias) < w
}

func main() {
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			log.Fatal().Err(err).Msg("httpServer exited")
		}
	}()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	repository := initFileRepository()
	cardService := cards.NewCardService(repository)

	generator, done := cardService.RandomCardGenerator()

	go func() {
		<-sigChannel
		log.Debug().Msg("sigterm received")
		done()
		time.Sleep(100 * time.Microsecond)
		os.Exit(0)
	}()

	i := 1
	for card := range generator {

		recalls := cardService.CountRecallAttempts(card.ID)
		confidence := cardService.EstimateCardConfidence(recalls)
		if !shouldBeShown(confidence) {
			fmt.Printf("Skipping \"%s\"\n", card.Phrase)
			continue
		}

		fmt.Fprintf(os.Stdout, "Recall #%d, card [#%d|%d%%|%+v]\n", i, card.ID, confidence, recalls)

		_, err := readlineAfter(card.Translation)
		if err != nil {
			log.Fatal().Err(err).Msgf("Could not display card's translation")
		}

		fmt.Print(card.Phrase + "\n")
		for {
			answer, _ := readlineAfter("Got it right? [y/n]\n")
			if answer == "y" {
				err := cardService.RecordRecallAttempt(card.ID, true)
				if err != nil {
					log.Fatal().Err(err).Msgf("Could not record a recall attempt")
				}
				break
			} else if answer == "n" {
				err := cardService.RecordRecallAttempt(card.ID, false)
				if err != nil {
					log.Fatal().Err(err).Msgf("Could not record a recall attempt")
				}
				break
			}
		}
		i++
	}
}
