package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mudithanda/students-api/internal/config"
	"github.com/mudithanda/students-api/internal/http/handlers/student"
	"github.com/mudithanda/students-api/internal/storage/sqlite"
)

func main() {

	// load config
	cfg := config.MustLoad()

	// database setup

	storage, err := sqlite.New(cfg) // _ => storage
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initialized: ", slog.String("env:", cfg.Env), slog.String("version:", "1.0.0"))

	// setup router

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	router.HandleFunc("PATCH /api/students/{id}", student.UpdateList(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteList(storage))

	
	// setup server

	server := http.Server{
		Addr: cfg.Address,
		Handler: router,
	}

	//fmt.Println("Server started")
	// fmt.Printf("Server started %s", cfg.Address)
	//slog.Info("Server started %s", cfg.Address)
	slog.Info("Server started", slog.String("Address ->", cfg.Address))

	// make channel for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func ()  {
		err := server.ListenAndServe() // blocking
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<- done

	slog.Info("shutting down server...")

	cxt, cancel := context.WithTimeout(context.Background(), 5 *time.Second)
	defer cancel()

	if err := server.Shutdown(cxt); err != nil {
		slog.Error("Failed to shutdown server", slog.Any("error", err))
	}

	slog.Info("Server shutdown successfully")

}