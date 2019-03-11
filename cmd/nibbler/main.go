package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/krak3n/nibbler/internal/config"
	"github.com/krak3n/nibbler/internal/server"
	"github.com/krak3n/nibbler/internal/storage"
	"github.com/krak3n/nibbler/internal/storage/psql"
)

func main() {
	log.Println("starting nibbler")

	dsn := storage.DSN{
		Name:    config.DBName,
		User:    config.DBUser,
		Pass:    config.DBPassword,
		Host:    config.DBHost,
		SSLMode: config.DBSSLMode,
	}

	db, err := psql.New(dsn)
	if err != nil {
		log.Fatalf("error connecting to DB: %s", err)
	}

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	srv := server.New(db, server.WithAddress(":3000"))
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(fmt.Sprintf("error starting server: %s", err))
			sigC <- syscall.SIGQUIT
		}
	}()

	log.Println("stopping nibbler:", <-sigC)

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Println(fmt.Sprintf("error stopping server: %s", err))
	}

	log.Println("stopped nibbler")
}
