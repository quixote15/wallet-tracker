package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"wallet-tracker/controllers"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get configuration from environment
	infuraProjectID := os.Getenv("INFURA_PROJECT_ID")
	ethereumNetwork := os.Getenv("ETHEREUM_NETWORK")
	if infuraProjectID == "" || ethereumNetwork == "" {
		log.Fatal("Missing required environment variables")
	}

	// Connect to Ethereum client
	infuraURL := fmt.Sprintf("https://%s.infura.io/v3/%s", ethereumNetwork, infuraProjectID)
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize router
	router := chi.NewRouter()

	// Initialize controllers
	controllerList := []controllers.Controller{
		controllers.NewHealthController(),
		controllers.NewWalletController(client),
	}

	// Register routes from all controllers
	for _, controller := range controllerList {
		for _, route := range controller.Routes() {
			router.Method(route.Method, route.Path, route.Handler)
		}
	}

	log.Printf("Server starting on port 3000...")
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}
