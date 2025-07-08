package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"streamit/utils"
)

func (s *RTMPSession) StartFFmpegConversion() {
	hlsPath := fmt.Sprintf("/tmp/hls/%s/%s", s.channel, s.key) // Local storage for HLS segments
	os.MkdirAll(hlsPath, os.ModePerm)

	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-i", fmt.Sprintf("rtmp://localhost/%s/%s", s.channel, s.key), // Input RTMP URL
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-b:v", "3000k",
		"-c:a", "aac",
		"-f", "hls",
		"-hls_time", "4",
		"-hls_list_size", "10",
		"-hls_flags", "delete_segments",
		"-hls_segment_filename", filepath.Join(hlsPath, "segment_%03d.ts"),
		filepath.Join(hlsPath, "index.m3u8"),
	)

	ffmpegCmd.Stdout = os.Stdout
	ffmpegCmd.Stderr = os.Stderr

	err := ffmpegCmd.Start()
	go utils.WatchAndUpload(hlsPath, s.key)

	if err != nil {
		fmt.Println("Failed to start FFmpeg:", err)
		return
	}

	s.ffmpegProcess = ffmpegCmd // Store the process to stop later

	// Start S3 sync in a separate goroutine
	// go s.SyncHLSToS3(hlsPath)
}
