package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type TechRadar struct {
	Technology string `json:"technology"`
	Status     string `json:"status"`
	Category   string `json:"category"`
}

type Entry struct {
	Quadrant int    `json:"quadrant"`
	Ring     int    `json:"ring"`
	Label    string `json:"label"`
	Active   bool   `json:"active"`
	Moved    int    `json:"moved"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <json_folder_path> [output_file]")
	}

	jsonFolder := os.Args[1]
	outputFile := ""
	if len(os.Args) >= 3 {
		outputFile = os.Args[2]
	}

	// Read all JSON files from the folder
	var entries []Entry

	err := filepath.WalkDir(jsonFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-JSON files
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".json") {
			return nil
		}

		// Read and parse JSON file
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Failed to read file %s: %v", path, err)
			return nil // Continue with other files
		}

		var techRadar TechRadar
		if err := json.Unmarshal(data, &techRadar); err != nil {
			log.Printf("Failed to unmarshal JSON from %s: %v", path, err)
			return nil // Continue with other files
		}

		// Convert to Entry format
		entry := Entry{
			Quadrant: getQ(techRadar.Category),
			Ring:     getR(techRadar.Status),
			Label:    techRadar.Technology,
			Active:   true,
			Moved:    0,
		}

		entries = append(entries, entry)
		fmt.Printf("Processed: %s -> %s (Q:%d, R:%d)\n", path, entry.Label, entry.Quadrant, entry.Ring)

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walk directory: %v", err)
	}

	if len(entries) == 0 {
		log.Fatal("No valid JSON files found in the specified folder")
	}

	// Encode to JSON
	var output []byte
	var encodeErr error

	if outputFile != "" {
		// Write to file with pretty formatting
		output, encodeErr = json.MarshalIndent(entries, "", "  ")
		if encodeErr != nil {
			log.Fatalf("Failed to marshal JSON: %v", encodeErr)
		}

		if err := os.WriteFile(outputFile, output, 0644); err != nil {
			log.Fatalf("Failed to write output file: %v", err)
		}
		fmt.Printf("\nOutput written to: %s\n", outputFile)
	} else {
		// Write to stdout
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(entries); err != nil {
			log.Fatalf("Failed to encode JSON to stdout: %v", err)
		}
	}

	fmt.Printf("\nProcessed %d JSON files successfully\n", len(entries))
}

func getQ(v string) int {
	switch v {
	case "Data management":
		return 1
	case "Infrastructure":
		return 0
	case "Languages & Frameworks":
		return 3
	default:
		return 2
	}
}

func getR(v string) int {
	switch strings.ToLower(v) {
	case "adopt":
		return 0
	case "assess":
		return 2
	case "hold":
		return 3
	default:
		return 1
	}
}
