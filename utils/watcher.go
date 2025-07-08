package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func WatchAndUpload(dir string, streamKey string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Watcher error:", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(dir)
	if err != nil {
		log.Println("Failed to watch:", err)
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&(fsnotify.Create|fsnotify.Write) > 0 {
				if strings.HasSuffix(event.Name, ".ts") || strings.HasSuffix(event.Name, ".m3u8") {
					fileName := filepath.Base(event.Name)
					s3Key := fmt.Sprintf("%s/%s", streamKey, fileName)
					go func(p string, k string) {
						err := UploadFileToS3(p, k)
						if err != nil {
							log.Println("S3 Upload failed:", err)
						} else {
							log.Println("Uploaded:", k)
						}
					}(event.Name, s3Key)
				}
			}
		case err := <-watcher.Errors:
			log.Println("Watcher error:", err)
		}
	}
}
