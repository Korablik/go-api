package main

import (
	"fmt"
	"log"
	"net/http"

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
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	log.Info("Starting the application...")

	http.HandleFunc("/hello-world", helloHandler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	log.Info("Application finished.")
}
