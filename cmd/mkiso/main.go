package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kdomanski/iso9660"
)

func main() {
	writer, err := iso9660.NewWriter()
	if err != nil {
		log.Fatalf("failed to create writer: %s", err)
	}
	defer writer.Cleanup()

	isoFile, err := os.OpenFile(os.Args[1], os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}
	defer isoFile.Close()

	for _, folderPath := range os.Args[2:] {
		walk_err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatalf("walk: %s", err)
				return err
			}
			if info.IsDir() {
				return nil
			}
			outputPath := strings.TrimPrefix(path, folderPath) // remove the source drive name
			outputPath = strings.TrimPrefix(outputPath, string(os.PathSeparator))
			fmt.Printf("Adding file: %s\n", outputPath)

			fileToAdd, err := os.Open(path)
			if err != nil {
				log.Fatalf("failed to open file: %s", err)
			}
			defer fileToAdd.Close()

			err = writer.AddFile(fileToAdd, outputPath)
			if err != nil {
				log.Fatalf("failed to add file: %s", err)
			}
			return nil
		})
		if walk_err != nil {
			log.Fatalf("%s", walk_err)
		}
	}

	err = writer.WriteTo(isoFile, "Test")
	if err != nil {
		log.Fatalf("failed to write ISO image: %s", err)
	}
}
