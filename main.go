package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	modeFlag := flag.String("mode", "server", "Confirm device mode.")
	modeTarget := flag.String("target", "127.0.0.1:8080", "Target IP address.")
	modePath := flag.String("path", "", "Path to save file location.")

	flag.Parse()
	if (*modeFlag == "server") {
		if (*modePath == "") {
			fmt.Println("❌ Error: The -path flag is required to specify a save location.")
			flag.Usage()
			os.Exit(1)
		}
		var serverPath = *modePath
		fmt.Println("SyncMe Server listening on port 9090...")
		RunServer("9090", serverPath)
	} else if (*modeFlag == "client") {
		if *modePath == "" {
			fmt.Println("❌ Error: The -path flag is required to specify a save location.")
			flag.Usage()
			os.Exit(1)
		}
		var clientPath = *modePath
		clientFiles, err := ScanDirectory(clientPath)
			if (err != nil) {
				fmt.Printf("Error scanning provided directory. %v\n", err)
				flag.Usage()
				os.Exit(1)
			}

		RunClient(*modeTarget, clientFiles, clientPath)
	} else if (*modeFlag == "web") {
		if *modePath == "" {
			fmt.Println("❌ Error: The -path flag is required to specify a save location.")
			flag.Usage()
			os.Exit(1)
		}

		StartDashboardServer("8080", *modePath)
	}
}
