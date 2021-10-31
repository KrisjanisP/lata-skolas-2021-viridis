package main

import (
	"log"
	"net/http"

	"github.com/KrisjanisP/viridis/database"
	"github.com/KrisjanisP/viridis/platform/authenticator"
	"github.com/KrisjanisP/viridis/platform/router"
	"github.com/KrisjanisP/viridis/utils"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbapi, err := database.NewDB()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := router.New(auth)
	go utils.StartWorker(dbapi)

	if err := http.ListenAndServe(":8080", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
