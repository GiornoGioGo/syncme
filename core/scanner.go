package core

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"encoding/json"
	"sync"
)

func ScanDirectory(root string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
			if err != nil {
				fmt.Printf("Warning: Could not read metadata for %s: %v\n", path, err)
				return nil
			}
		
			if (info.Name() == "prefs.prop") {
				fmt.Println("Ignoring prefs.prop file")
				return nil
			}

			rel, err := filepath.Rel(root, path)
				if err != nil {
					fmt.Println("Could not parse files relative path")
					return nil
				}

			//open file
			file, err := os.Open(path)
			if err != nil {
				fmt.Printf("Failed to open file: %s\n", err)
				return nil
			}

			hasher := sha256.New()
				if _, err := io.Copy(hasher, file); err != nil {
					fmt.Println("Could not extract hsh from file!")
					file.Close()
					return nil
				}

			file.Close()
			f := FileInfo{
					FilePath:  rel,
					Hash:      fmt.Sprintf("%x", hasher.Sum(nil)),
					FileSize:  info.Size(),
					UpdatedAt: info.ModTime(),
				}
			files = append(files, f)
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory! %v\n", err)
		return nil, err // Added missing return values
	}

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(f FileInfo) {
			defer wg.Done() // Moved to top of goroutine to ensure execution
			fmt.Printf("Relative dir: %s\n", f.FilePath)
			fmt.Printf("File size: %d\n", f.FileSize)
			fmt.Printf("Hash: %s\n", f.Hash)
			fmt.Printf("Modification date: %s\n", f.UpdatedAt.String())
		}(file)
	}

	wg.Wait()

	output, err := json.MarshalIndent(files, "", " ")
	if err != nil {
		fmt.Printf("Problem converting to json. %v\n", err)
		return nil, err
	}

	fmt.Println(string(output))
	return files, nil
}