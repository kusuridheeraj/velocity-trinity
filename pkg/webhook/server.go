package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/velocity-trinity/core/pkg/logger"
	"github.com/velocity-trinity/core/pkg/scheduler"
)

type Server struct {
	Queue  scheduler.Queue
	Router *mux.Router
	Port   string
}

func NewServer(queue scheduler.Queue, port string) *Server {
	s := &Server{
		Queue:  queue,
		Router: mux.NewRouter(),
		Port:   port,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.Router.HandleFunc("/webhook", s.handleGitHubWebhook).Methods("POST")
}

// Minimal struct for GitHub PR payload
type PullRequestEvent struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	PullRequest struct {
		Head struct {
			Ref string `json:"ref"`
			Sha string `json:"sha"`
		} `json:"head"`
		Base struct {
			Ref string `json:"ref"`
		} `json:"base"`
	} `json:"pull_request"`
}

func (s *Server) handleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	// In production, verify X-Hub-Signature-256 here!
	
	var event PullRequestEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		logger.Log.Error("Failed to decode webhook: " + err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	logger.Log.Info(fmt.Sprintf("Received webhook: PR #%d %s", event.Number, event.Action))

	if event.Action == "opened" || event.Action == "synchronize" {
		job := &scheduler.Job{
			ID:       fmt.Sprintf("pr-%d-%s", event.Number, event.PullRequest.Head.Sha[:7]),
			PRNumber: event.Number,
			// Simplified: We assume base PR based on branch name convention or other logic
			// For MVP, we treat it as a regular job
			BasePR: 0, 
		}
		
		if err := s.Queue.Enqueue(job); err != nil {
			logger.Log.Error("Enqueue failed: " + err.Error())
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}

func (s *Server) ListenAndServe() error {
	logger.Log.Info("Quantum Merge Webhook Server listening on port " + s.Port)
	return http.ListenAndServe(":"+s.Port, s.Router)
}
