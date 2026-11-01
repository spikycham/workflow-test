package main

import (
	"log"
	"net/http"

	"github.com/google/go-github/v93/github"
)

func main() {
	http.HandleFunc("POST /webhook/run", handleWebhook)

	log.Print("[INFO] started service at :3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Printf("[ERROR] start error")
	}
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(""))
	if err != nil {
		http.Error(w, "invalid secret", http.StatusBadRequest)
		return
	}

	eventType := github.WebHookType(r)
	event, err := github.ParseWebHook(eventType, payload)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
	}

	log.Printf("[INFO] event: %v", event)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
