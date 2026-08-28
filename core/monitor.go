package core

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func RunProcessMonitorWorker(share ShareProfile, targetIp string) {
	var isGameRunning bool = false

	for {
		running := isProcessRunning(share.ProcessName)

		if (running && !isGameRunning) {
			fmt.Printf("Detected game startup: %s", share.Name)
			isGameRunning = true
		} else if (!running && isGameRunning) {
			fmt.Printf("💾 Game closed: %s. Initiating SyncMe data transfer...\n", share.Name)
    		isGameRunning = false

			files, err := ScanDirectory(share.LocalPath)
				if (err != nil) {
					fmt.Printf("Error scanning local path. %v\n", err)
					continue
				}
			RunClient(targetIp, files, share.LocalPath)
		} 

		time.Sleep(3 * time.Second)
	}
}

func RunTimerWorker(share ShareProfile, targetIp string) {
	ticker := time.NewTicker(time.Duration(share.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Printf("Timer triggered: Running background sync for %s\n", share.Name)
		files, err := ScanDirectory(share.LocalPath)
			if (err != nil) {
				fmt.Printf("Error scanning local path. %v\n", err)
				continue
			}

		RunClient(targetIp, files, share.LocalPath)
	}
}

func isProcessRunning(name string) bool {
	var cmd *exec.Cmd

	if (runtime.GOOS == "windows") {
		cmd = exec.Command("tasklist", "/FI", "IMAGENAME eq "+name)
	} else {
		cmd = exec.Command("pgrep", "-f", name)
	}

	output, err := cmd.Output()
		if (err != nil) {
			return false
		}

	cleanOutput := strings.ToLower(string(output))
	cleanName := strings.ToLower(name)

	return strings.Contains(string(cleanOutput), cleanName)
}