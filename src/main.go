package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

const (
	//TODO: move to config file
	port = 5000
)

// Handler function for the "/hello" endpoint
func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request:", r.Method, r.URL.Path)

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Hello, World!")
	log.Println("Response: 'Hello, World'")
}

func main() {
	// init logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Starting the application")

	// Create a new HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: http.HandlerFunc(helloHandler),
	}

	// Channel to notify the main routine of server shutdown
	done := make(chan bool, 1)

	go signHandler(server, done, logger)

	// Start the server
	logger.Info("Starting server", zap.Int("port", port))
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("Could not start server", zap.Error(err))
	}

	// Wait for the server to gracefully shutdown
	<-done

	logger.Info("Server STOPPED")
	os.Exit(0)
}

// Goroutine to handle signal
func signHandler(server *http.Server, done chan bool, logger *zap.Logger) {
	// Channel to listen for interrupt or terminate signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for an interrupt or terminate signal
	sig := <-sigChan

	logger.Info("Received signal. Shutting down server...", zap.String("signal", sig.String()))
	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", zap.Error(err))
	}

	close(done)
}
