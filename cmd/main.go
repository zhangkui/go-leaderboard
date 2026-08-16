package main

import (
	"log"
	"net/http"

	"go-leaderboard/internal/api"
	"go-leaderboard/internal/leaderboard"
)

func main() {
	lb := leaderboard.New()
	mux := api.NewMux(lb)
	log.Println("leaderboard service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
