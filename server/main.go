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
		return
	}

	switch e := event.(type) {
	case *github.WorkflowRunEvent:
		// action := e.GetAction()
		status := e.WorkflowRun.GetStatus()
		log.Printf("[INFO] status=%s", status)
		//
		// if action == "completed" && e.WorkflowRun.Conclusion != nil {
		// 	conclusion := e.WorkflowRun.GetConclusion()
		//
		// 	if conclusion == "failure" {
		// 		http.Error(w, "workflow failed", http.StatusBadRequest)
		// 		log.Println("[ERROR] workflow failed")
		// 		return
		// 	} else if conclusion == "completed" {
		// 		// TODO: deploy docker image
		// 		log.Printf("[SUCCESS] workflow completed")
		//
		// 		w.WriteHeader(http.StatusOK)
		// 		w.Write([]byte("OK"))
		// 		return
		// 	}
		// }
	default:
		log.Print("[INFO] default event")
		http.Error(w, "not workflow run event", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
