package main

import (
	"log"
	"math/rand"
	"net/http"
)

func runServer(successRate int) {
	handler := createApp(successRate)

	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatal(err)
	}
}

func createApp(successRate int) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", unstableHandler(successRate))
	return mux
}

func unstableHandler(successRate int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x := rand.Intn(100)

		if x < successRate {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
