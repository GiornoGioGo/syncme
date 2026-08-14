package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
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

func RunClient(targetIP string, files []FileInfo, srcRoot string) {
	jsonData, err := json.Marshal(files)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	jsonData = append(jsonData, '\n')

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

	client.SetReadDeadline(time.Now().Add(5 * time.Second))

	reader := bufio.NewReader(client)

	data, err := reader.ReadBytes('\n')
		if (err != nil) {
			fmt.Printf("Security alert or read failure: %v", err)
			return
		}

	var filesToUpload[]string

	err = json.Unmarshal(data, &filesToUpload)
		if (err != nil) {
			fmt.Printf("Failed to parse incoming JSON: %v\n", err)
    		return
		}

	for _, f := range filesToUpload {
		fmt.Printf("Server requested these files for upload: %s\n", f)
	}

	
	//upload file(s) to server
	for _, f := range filesToUpload {
		for _, file := range files {
			if (f == file.FilePath) {
				h := FileHeader{
					RelPath:  f,
					FileSize:  file.FileSize,
				}
				
				headerBytes, err := json.Marshal(h)
					if (err != nil) {
						fmt.Printf("Error marshalling file header: %v\n", err)
                		continue
					}
				
				headerBytes = append(headerBytes, '\n')

				_, err = client.Write(headerBytes)
					if (err != nil) {
						fmt.Printf("Failed to send file header: %v\n", err)
						return
					}

				fmt.Printf("Sent header for: %s (%d bytes)\n", h.RelPath, h.FileSize)
				fullPath := filepath.Join(srcRoot, f)

				file, err := os.Open(fullPath)
					if (err != nil) {
						fmt.Printf("Failed to open %s. %v\n", f, err)
						return
					}
				
				written, err := io.CopyN(client, file, h.FileSize)
					if (err != nil) {
						fmt.Printf("Failed to copy %s. %v\n", f, err)
						return
					}

				file.Close()
				fmt.Printf("written: %v\n", written)
			}
		}
	}
}

func handleConnection(conn net.Conn, localPath string) {
	defer conn.Close()
	fmt.Printf("Client connected from: %s\n", conn.RemoteAddr().String())

	const maxSteamSaveSize = 10 * 1024 * 1024 
	limitedReader := io.LimitReader(conn, maxSteamSaveSize)

	reader := bufio.NewReader(limitedReader)

	data, err := reader.ReadBytes('\n')
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

	responseJSON, err := json.Marshal(filesToRequest)
		if (err != nil) {
			fmt.Printf("Error marshalling data: %v\n", err)
			return
		}

	responseJSON = append(responseJSON, '\n')
	
	_, err = conn.Write(responseJSON)
		if (err != nil) {
			fmt.Printf("Error writing response data: %v\n", err)
			return
		}

	fmt.Printf("Received manifest containing %d files:\n", len(recievedFiles))
	for _, f := range recievedFiles {
		fmt.Printf(" - %s (%d bytes) [Hash: %s]\n", f.FilePath, f.FileSize, f.Hash)
	}

	for range filesToRequest {
		data, err := reader.ReadBytes('\n')
			if (err != nil) {
				fmt.Printf("Security alert or read failure: %v", err)
				return
			}
		
		var h FileHeader
		err = json.Unmarshal(data, &h)
			if (err != nil) {
				fmt.Printf("Failed to parse incoming JSON: %v\n", err)
    			return
			}
		
		destPath := filepath.Join(localPath, h.RelPath)
		err = os.MkdirAll(filepath.Dir(destPath), 0755)
			if (err != nil) {
				fmt.Printf("Failed to create subfolders at %s. %v\n", localPath, err)
				return
			}
		
		newFile, err := os.Create(destPath)

		written, err := io.CopyN(newFile, reader, h.FileSize)
			if (err != nil) {
				fmt.Printf("Failed to copy %s. %v\n", newFile, err)
				return
			}
		
		newFile.Close()
		fmt.Printf("written: %v\n", written)
	}
}