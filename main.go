package main

import (
	"e_healthy/models"
	"e_healthy/router"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func init() {
	models.Setup()
}

func main() {
	serverBasic := &http.Server{
		Addr:         ":7001",
		Handler:      router.RouterBasic(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serverAuth := &http.Server{
		Addr:         ":7002",
		Handler:      router.RouterAuth(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serverClient := &http.Server{
		Addr:         ":7003",
		Handler:      router.RouterClient(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return serverBasic.ListenAndServe()
	})

	g.Go(func() error {
		return serverAuth.ListenAndServe()
	})

	g.Go(func() error {
		return serverClient.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
