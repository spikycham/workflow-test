package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

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
		action := e.GetAction()
		if action == "completed" && e.WorkflowRun.Conclusion != nil {
			conclusion := e.WorkflowRun.GetConclusion()
			if conclusion == "success" {
				run := e.WorkflowRun
				sha := run.GetHeadSHA()[:4]
				runDocker("ci-test-web", fmt.Sprintf("450139220/ci-test-web:%s", sha))
			}
		} else {
			http.Error(w, "not required", http.StatusBadRequest)
			return
		}
	default:
		log.Print("[INFO] default event")
		http.Error(w, "not workflow run event", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func runDocker(containerName, imageName string) {
	ctx, cacenl := context.WithTimeout(context.Background(), 30*time.Second)
	defer cacenl()

	stop := exec.CommandContext(ctx, "docker", "stop", containerName)
	if err := stop.Run(); err != nil {
		log.Printf("[ERROR] failed to stop: %v", err)
		return
	}
	rm := exec.CommandContext(ctx, "docker", "rm", containerName)
	if err := rm.Run(); err != nil {
		log.Printf("[ERROR] failed to remove: %v", err)
		return
	}

	pull := exec.CommandContext(ctx, "docker", "pull", imageName)
	if err := pull.Run(); err != nil {
		log.Printf("[ERROR] failed to pull: %v", err)
		return
	}

	cmd := exec.CommandContext(ctx, "docker", "run", "--name", containerName, "-p", "20004:80", "-d", imageName)
	if err := cmd.Run(); err != nil {
		log.Printf("[ERROR] failed to run docker: %v", err)
		return
	}

	log.Printf("[INFO] success to run docker")
}
