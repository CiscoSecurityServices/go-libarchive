package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	ar "github.com/CiscoSecurityServices/go-libarchive"
)

func extractContents(srcFilename, destFilename string) error {
	fmt.Printf("from: %s to: %s\n", srcFilename, destFilename)

	file, err := os.Open(srcFilename)
	if err != nil {
		return fmt.Errorf("error opening source file %s: %w", srcFilename, err)
	}
	reader, err := ar.NewReader(file)
	if err != nil {
		return fmt.Errorf("error opening libarchive reader %s: %w", srcFilename, err)
	}
	defer reader.Close()

	for {
		entry, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading next entry from archive %s: %w", srcFilename, err)
		}

		if err := save(filepath.Join(destFilename, entry.PathName()), entry.Stat().Mode(), reader); err != nil {
			return fmt.Errorf("error saving entry %s: %w", entry.PathName(), err)
		}
	}
	return nil
}

func save(path string, mode os.FileMode, r *ar.Reader) error {
	switch {
	case mode.IsDir():
		return os.MkdirAll(path, 0755)
	case mode.IsRegular():
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err = io.Copy(file, r); err != nil {
			return fmt.Errorf("error copying file %s: %w", path, err)
		}
	default:
		fmt.Printf("Skipping unsupported file type: %s %s\n", path, mode)
	}
	return nil
}

func main() {
	fmt.Printf("libarchive version: %s\n", ar.LibArchiveVersion)

	if len(os.Args) < 3 {
		fmt.Println("Usage: simple-extract-copy <file1> <file2>")
		return
	}

	if err := extractContents(os.Args[1], os.Args[2]); err != nil {
		fmt.Printf("Extraction error: %v\n", err)
		os.Exit(1)
	}
}
