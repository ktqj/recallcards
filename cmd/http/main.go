package main

import (
	// "context"
	// "context"
	"net/http"
	"os"
	"time"

	// "fmt"

	// "os"
	// "os/signal"

	"example.com/recallcards/pkg/cards"
	"example.com/recallcards/pkg/storage/file"
	"example.com/recallcards/pkg/web"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type server struct {
	c      web.Controller
	router *mux.Router
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) routes() {
	s.router.HandleFunc("/", s.c.IndexHandler)

	s.router.HandleFunc("/api/cards/form_create/", s.c.CreateCardForm)
	s.router.HandleFunc("/api/cards/json_create/", s.c.CreateCardJson)
}

func NewServer(c web.Controller) *server {
	s := &server{
		c:      c,
		router: mux.NewRouter(),
	}
	s.routes()
	return s
}

func initInMemRepository() cards.Repository {
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
	repository := initInMemRepository()
	cardService := cards.NewCardService(repository)
	controller := web.NewController(cardService)

	srv := NewServer(controller)

	httpServer := &http.Server{
		Addr:         "127.0.0.1:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      srv, // Pass our instance of gorilla/mux in.
	}
	httpServer.RegisterOnShutdown(func() {
		log.Debug().Msg("shutdown callback")
		time.Sleep(10 * time.Second)
	})

	// idleConnsClosed := make(chan struct{})
	// go func() {
	// 	c := make(chan os.Signal, 1)
	// 	signal.Notify(c, os.Kill, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	// 	// signal.Notify(c, os.Interrupt)
	// 	// signal.Notify(c, scall.SIGTERM)
	// 	// signal.Notify(c, scall.SIGQUIT)
	// 	s := <-c
	// 	log.Debug().Msgf("received signal %v", s)
	// 	log.Error().Msgf("HTTP server Shutdown:")
	// 	// We received an interrupt signal, shut down.
	// 	if err := httpServer.Shutdown(context.Background()); err != nil {
	// 		// Error from closing listeners, or context timeout:
	// 		log.Error().Msgf("HTTP server Shutdown: %v", err)
	// 	}
	// 	log.Debug().Msg("shutdown unblocked")
	// 	close(idleConnsClosed)
	// }()

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("httpServer exited")
	}
	// time.Sleep(10 * time.Second)
	// log.Debug().Msg("listenandserve unblocked")
	// <-idleConnsClosed
}
