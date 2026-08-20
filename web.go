package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"time"
)

//go:embed web/*
var webFiles embed.FS

// StartDashboardServer boots up the internal web configuration node
func StartDashboardServer(port string, localPath string) {
	mux := http.NewServeMux()

	publicFiles, err := fs.Sub(webFiles, "web")
		if (err != nil) {
			fmt.Printf("Error mapping embedded asset folders: %v\n", err)
			return
		}

	fileServer := http.FileServer(http.FS(publicFiles))

	mux.Handle("/", fileServer)

	// Live Data API Endpoint - Returns system metadata states to the UI
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		statusPayload := map[string]interface{}{
			"status":         "idle",
			"node_version":   "1.0.0",
			"managed_path":   localPath,
			"last_heartbeat": time.Now().Format(time.RFC3339),
		}

		json.NewEncoder(w).Encode(statusPayload)
	})

	mux.HandleFunc("/api/sync", func(w http.ResponseWriter, r *http.Request) {
		if (r.Method != http.MethodPost) {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
			return
		}

		files, err := ScanDirectory(localPath)
			if (err != nil) {
				fmt.Printf("Error scanning directory. %v\n", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		
		RunClient("127.0.0.1:9090", files, localPath)

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{"message": "Sync Complete"})
		
	}) 

	fmt.Printf("🌐 SyncMe Web UI listening at http://localhost:%s\n", port)
	
	// Start the blocking web loop interceptor
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("❌ Web Server Crash: %v\n", err)
	}
}
