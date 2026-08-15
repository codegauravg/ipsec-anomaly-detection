package main

// # Run the script directly
// go run replace.go

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Directories and file types to safely ignore during traversal
var ignoreDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"dist":         true,
	"bin":          true,
}

func main() {
	// Define command-line flags
	configPath := flag.String("config", "replacements.json", "Path to the JSON configuration file")
	targetDir := flag.String("dir", ".", "Target directory to run the replacements")
	flag.Parse()

	// 1. Read and parse the JSON file
	configBytes, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config file '%s': %v", *configPath, err)
	}

	var replacements map[string]string
	if err := json.Unmarshal(configBytes, &replacements); err != nil {
		log.Fatalf("Failed to parse JSON config: %v", err)
	}

	if len(replacements) == 0 {
		log.Println("No replacements found in the JSON file. Exiting.")
		return
	}

	fmt.Printf("Loaded %d replacement rules. Scanning directory: %s\n\n", len(replacements), *targetDir)

	// 2. Walk through the directory tree recursively
	err = filepath.WalkDir(*targetDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err // Skip files we don't have permissions to read
		}

		// Skip ignored directories to save time and prevent corruption
		if d.IsDir() {
			if ignoreDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip the script itself and the config file
		if filepath.Base(path) == "replace.go" || filepath.Base(path) == filepath.Base(*configPath) {
			return nil
		}

		// 3. Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip unreadable files silently
		}

		// 4. Perform replacements
		originalText := string(content)
		newText := originalText
		modified := false

		for oldStr, newStr := range replacements {
			if strings.Contains(newText, oldStr) {
				newText = strings.ReplaceAll(newText, oldStr, newStr)
				modified = true
			}
		}

		// 5. Write the modified content back to the file
		if modified {
			// Get original file permissions
			info, err := d.Info()
			if err != nil {
				log.Printf("Could not get file info for %s: %v", path, err)
				return nil
			}

			if err := os.WriteFile(path, []byte(newText), info.Mode()); err != nil {
				log.Printf("Failed to write updates to %s: %v", path, err)
			} else {
				fmt.Printf("Updated -> %s\n", path)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error traversing directory: %v", err)
	}

	fmt.Println("\nPlaceholder replacement completed successfully.")
}
