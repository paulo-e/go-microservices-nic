package main

import (
	"context"
	"go-microservices-nic/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "[product-api] ", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/products/", ph)

	s := &http.Server{
		Addr:         ":8000",
		Handler:      sm,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	l.Printf("server: Starting server on '%s'...", s.Addr)

	go func() {
		l.Fatalln(s.ListenAndServe())
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Printf("server: Received %s, gracefully shutting down", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	err := s.Shutdown(tc)
	if err != nil {
		l.Fatalln(err)
	}
}
