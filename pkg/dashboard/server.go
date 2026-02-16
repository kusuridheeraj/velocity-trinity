package dashboard

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/velocity-trinity/core/pkg/logger"
	"github.com/velocity-trinity/core/pkg/scheduler"
)

//go:embed ui/*
var uiFiles embed.FS

type DashboardServer struct {
	Queue scheduler.Queue
}

func RegisterRoutes(router *mux.Router, queue scheduler.Queue) {
	d := &DashboardServer{Queue: queue}

	// Serve Static Files
	uiFS, _ := fs.Sub(uiFiles, "ui")
	router.PathPrefix("/").Handler(http.FileServer(http.FS(uiFS)))

	// API Routes
	router.HandleFunc("/api/jobs", d.listJobs).Methods("GET")
}

func (d *DashboardServer) listJobs(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, we'd need a method on Queue to ListAll()
	// Since our interface only has Get/Dequeue, let's cheat and cast to MemoryQueue
	// Or better, update the interface. For MVP, we'll return a static list if interface doesn't support it.
	
	if mq, ok := d.Queue.(*scheduler.MemoryQueue); ok {
		jobs := mq.ListAll()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jobs)
	} else {
		// Fallback for other queue types not implementing list
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}
}
