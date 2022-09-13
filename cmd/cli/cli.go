package main

import (
	"bufio"
	"context"
	"fmt"
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

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	generator, err := cardService.FilteredRandomCardGenerator(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not initialize cards stream")
		return
	}

	go func() {
		<-sigChannel
		log.Debug().Msg("sigterm received")
		cancel()
		time.Sleep(100 * time.Microsecond)
		os.Exit(0)
	}()

	i := 1
	for card := range generator {

		recalls := cardService.RecallSummary(card.ID)

		fmt.Fprintf(os.Stdout, "Recall #%d, card [#%d|%d%%|%+v]\n", i, card.ID, recalls.EstimateConfidence(), recalls)

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
