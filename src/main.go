package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	log.Info("Starting the application...")

	// Print Hello, World! to the console
	fmt.Println("Hello, World!")

	log.Info("Application finished.")
}
