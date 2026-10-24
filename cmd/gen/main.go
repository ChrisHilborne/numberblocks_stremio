// Command gen builds a static Stremio addon that maps Numberblocks episodes
// (Cinemeta id tt7978538) to videos on the official Numberblocks YouTube channel.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const seriesID = "tt7978538"

type episode struct {
	YouTubeID string `json:"yt"`
	Title     string `json:"title"`
}

type stream struct {
	ExternalURL string `json:"externalUrl"`
	Name        string `json:"name"`
	Title       string `json:"title"`
}

type streamResponse struct {
	Streams []stream `json:"streams"`
}

type resource struct {
	Name       string   `json:"name"`
	Types      []string `json:"types"`
	IDPrefixes []string `json:"idPrefixes"`
}

type manifest struct {
	ID          string     `json:"id"`
	Version     string     `json:"version"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Resources   []resource `json:"resources"`
	Types       []string   `json:"types"`
	Catalogs    []struct{} `json:"catalogs"`
}

func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func run(dataPath, outDir string) error {
	raw, err := os.ReadFile(dataPath)
	if err != nil {
		return err
	}
	var eps map[string]episode // key "season:episode"
	if err := json.Unmarshal(raw, &eps); err != nil {
		return err
	}
	m := manifest{
		ID:          "com.chrishilborne.numberblocks-youtube",
		Version:     "1.2.0",
		Name:        "Numberblocks (YouTube)",
		Description: "Plays Numberblocks episodes from the official Numberblocks YouTube channel.",
		Resources:   []resource{{Name: "stream", Types: []string{"series"}, IDPrefixes: []string{seriesID}}},
		Types:       []string{"series"},
		Catalogs:    []struct{}{},
	}
	if err := writeJSON(filepath.Join(outDir, "manifest.json"), m); err != nil {
		return err
	}
	for key, ep := range eps {
		resp := streamResponse{Streams: []stream{{
			// Hands off to the native YouTube app; in-app ytId playback lags on TVs.
			ExternalURL: "https://www.youtube.com/watch?v=" + ep.YouTubeID,
			Name:        "YouTube app",
			Title:       ep.Title + "\nOpen in the YouTube app",
		}}}
		path := filepath.Join(outDir, "stream", "series", fmt.Sprintf("%s:%s.json", seriesID, key))
		if err := writeJSON(path, resp); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(outDir, ".nojekyll"), nil, 0o644); err != nil {
		return err
	}
	log.Printf("wrote manifest + %d episode streams to %s", len(eps), outDir)
	return nil
}

func main() {
	if err := run("data/episodes.json", "docs"); err != nil {
		log.Fatal(err)
	}
}
