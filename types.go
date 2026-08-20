package main

import "time"

type FileInfo struct {
	FilePath string `json:"file_path"`
	FileSize int64 `json:"file_size"`
	Hash string `json:"hash"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FileHeader struct {
	RelPath string `json:"rel_path"`
	FileSize int64 `json:"file_size"`
}

type Config struct {
	DeviceName string `json:"device_name"`
	TargetIP string `json:"target_ip"`
	LocalSavePath string `json:"local_save_path"`
	IgnoredFiles []string `json:"ignored_files"`
}