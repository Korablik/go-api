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

	"github.com/sirupsen/logrus"
)

// Handler function for the "/hello" endpoint
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Log the request
	log.Println("Request:", r.Method, r.URL.Path)

	// Respond with "Hello, World!"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World!")
	log.Println("Response: 'Hello, World'")
}

func main() {
	port := 5000

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	log.Info("Starting the application...")

	// Create a new HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: http.HandlerFunc(helloHandler),
	}

	// Channel to notify the main routine of server shutdown
	done := make(chan bool, 1)

	go signHandler(server, done, log)

	// Start the server
	log.Infof("Starting server on port: %v", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Errorf("Could not start server: %v", err)
	}

	// Wait for the server to gracefully shutdown
	<-done
	log.Info("Server STOPPED")
	os.Exit(0)
}

// Goroutine to handle signal
func signHandler(server *http.Server, done chan bool, log *logrus.Logger) {
	// Channel to listen for interrupt or terminate signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for an interrupt or terminate signal
	sig := <-sigChan

	log.Infof("Received signal: %s, shutting down server...", sig)
	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	close(done)
}
