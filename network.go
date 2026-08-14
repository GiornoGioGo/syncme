package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

func RunServer(port string, localPath string) {
	listener, err := net.Listen("tcp", ":"+port)
		if (err != nil) {
			fmt.Printf("Could not listen in on port %s. %v\n", port, err)
			return
		}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
			if (err != nil) {
				fmt.Printf("Error accepting connection: %v\n", err)
				continue
			}

		handleConnection(conn, localPath)
	}
}

func RunClient(targetIP string, files []FileInfo) {
	jsonData, err := json.Marshal(files)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	client, err := net.Dial("tcp", targetIP)
		if (err != nil) {
			fmt.Printf("Error: %v\n", err)
			return
		}

	defer client.Close()

	client.SetWriteDeadline(time.Now().Add(5 * time.Second))

	bytesWritten, err := client.Write(jsonData)
		if (err != nil) {
			fmt.Printf("Error writing json to client: %v\n", err)
			return
		}

	fmt.Printf("Successfully sent %d bytes.\n", bytesWritten)

	client.Close()
}

func handleConnection(conn net.Conn, localPath string) {
	defer conn.Close()
	fmt.Printf("Client connected from: %s\n", conn.RemoteAddr().String())

	const maxSteamSaveSize = 10 * 1024 * 1024 
	limitedReader := io.LimitReader(conn, maxSteamSaveSize)

	data, err := io.ReadAll(limitedReader)
		if (err != nil) {
			fmt.Printf("Security alert or read failure: %v", err)
			return
		}

	var recievedFiles[]FileInfo

	//client save folder
	err = json.Unmarshal(data, &recievedFiles)
		if (err != nil) {
			fmt.Printf("Failed to parse incoming JSON: %v\n", err)
    		return
		} 
	
	//server save folder
	serverFiles, err := ScanDirectory(localPath)
			if (err != nil) {
				fmt.Printf("Error parsing server files. %v\n", err)
				return
			}
	
	serverMap := make(map[string]FileInfo)
	for _, sFile := range serverFiles {
		serverMap[sFile.FilePath] = sFile
	}

	var filesToRequest []string

	for _, cFile := range recievedFiles {
		sFile, exists := serverMap[cFile.FilePath]
		
		if !exists {
			filesToRequest = append(filesToRequest, cFile.FilePath)
			continue 
		}
		
		if sFile.Hash != cFile.Hash {
			if cFile.UpdatedAt.After(sFile.UpdatedAt) {
				filesToRequest = append(filesToRequest, cFile.FilePath)
			}
		}
	}


	fmt.Printf("Received manifest containing %d files:\n", len(recievedFiles))
	for _, f := range recievedFiles {
		fmt.Printf(" - %s (%d bytes) [Hash: %s]\n", f.FilePath, f.FileSize, f.Hash)
	}

}