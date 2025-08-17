package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strings"
)

type Entry struct {
	Quadrant int    `json:"quadrant"`
	Ring     int    `json:"ring"`
	Label    string `json:"label"`
	Active   bool   `json:"active"`
	Moved    int    `json:"moved"`
}

func main() {
	f, err := os.Open("radar.csv")
	if err != nil {
		log.Fatalf("open csv: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	var out []Entry
	for {
		rec, err := r.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("read csv: %v", err)
		}
		out = append(out, Entry{
			Quadrant: getQ(rec[2]),
			Ring:     getR(rec[1]),
			Label:    rec[0],
			Active:   true,
			Moved:    0,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(out); err != nil {
		log.Fatalf("encode json: %v", err)
	}
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
