package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/idepin/pomodoro/internal/handler"
	"github.com/idepin/pomodoro/internal/repository"
	"github.com/idepin/pomodoro/internal/service"
)

func main() {
	repo := repository.NewMemoryItemRepo()
	svc := service.NewItemService(repo)
	h := handler.NewItemHandler(svc)

	mux := http.NewServeMux()
	h.Register(mux)

	srv := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
