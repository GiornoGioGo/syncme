package main

import "time"

type FileInfo struct {
	FilePath string `json:"file_path"`
	FileSize int64 `json:"file_size"`
	Hash string `json:"hash"`
	UpdatedAt  time.Time `json:"updated_at"`
}