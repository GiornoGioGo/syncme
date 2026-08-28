package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"git.giorno.dev/giorno/syncme/core"
)

func main() {
	configBytes, err := os.ReadFile("syncme.json")
		if (err != nil) {
			fmt.Printf("Error reading syncme.json. %v\n", err)
			return
		}
	var myConfig core.Config
	err = json.Unmarshal(configBytes, &myConfig)
		if (err != nil) {
			fmt.Printf("Failed to parse incoming JSON: %v\n", err)
    		return
		} 
	
	modeFlag := flag.String("mode", "server", "Confirm device mode.")
	modeTarget := flag.String("target", myConfig.TargetIP, "Target IP address.")
	modePath := flag.String("path", myConfig.LocalSavePath, "Path to save file location.")

	flag.Parse()
	if (*modeFlag == "server") {
		if (*modePath == "") {
			*modePath = myConfig.LocalSavePath
			if (*modePath == "") {
				fmt.Println("❌ Error: A directory path must be provided via syncme.json or the -path flag.")
    			os.Exit(1)
			}
		}
		var serverPath = *modePath
		fmt.Println("SyncMe Server listening on port 9090...")
		core.RunServer("9090", serverPath)
		
	} else if (*modeFlag == "client") {
		if *modePath == "" {
			*modePath = myConfig.LocalSavePath
			if (*modePath == "") {
				fmt.Println("❌ Error: A directory path must be provided via syncme.json or the -path flag.")
    			os.Exit(1)
			}
		}
		var clientPath = *modePath
			clientFiles, err := core.ScanDirectory(clientPath)
				if (err != nil) {
					fmt.Printf("Error scanning provided directory. %v\n", err)
					flag.Usage()
					os.Exit(1)
				}

			core.RunClient(*modeTarget, clientFiles, clientPath)
		
	} else if (*modeFlag == "web") {
		if *modePath == "" {
			*modePath = myConfig.LocalSavePath
			if (*modePath == "") {
				fmt.Println("❌ Error: A directory path must be provided via syncme.json or the -path flag.")
    			os.Exit(1)
			}
		}
		StartDashboardServer("8080", *modePath)
	} else if (*modeFlag == "daemon") {
		for _, share := range myConfig.Shares {

			currentShare := share

			if (currentShare.SyncStrategy == "timer") {
				go core.RunTimerWorker(share, myConfig.TargetIP)
			} else if (currentShare.SyncStrategy == "process_monitor") {
				go core.RunProcessMonitorWorker(share, myConfig.TargetIP)
			}
		}
		select {}
	} 
}
