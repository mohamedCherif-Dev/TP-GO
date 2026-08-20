package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

var tempDir string

func main() {
	var err error
	tempDir, err = os.MkdirTemp("", "tp4-server-*")
	if err != nil {
		log.Fatalf("Impossible de créer le dossier temporaire: %v", err)
	}
	defer os.RemoveAll(tempDir)

	files := generateSampleFiles(tempDir)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /files", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(files)
	})

	mux.HandleFunc("GET /files/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")

		cleanPath := filepath.Clean(name)
		if strings.Contains(cleanPath, "..") || filepath.IsAbs(cleanPath) {
			http.Error(w, "Accès interdit", http.StatusForbidden)
			return
		}

		filePath := filepath.Join(tempDir, cleanPath)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "Fichier non trouvé", http.StatusNotFound)
			return
		}

		time.Sleep(500 * time.Millisecond)

		http.ServeFile(w, r, filePath)
	})

	fmt.Println("==================================================")
	fmt.Println("Serveur TP4 démarré sur http://localhost:8080")
	fmt.Printf("Fichiers temporaires générés dans %s\n", tempDir)
	fmt.Println("==================================================")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func generateSampleFiles(dir string) []FileInfo {
	sampleList := []struct {
		name string
		size int64
	}{
		{"rapport.pdf", 1_254_400},
		{"presentation.pdf", 3_670_016},
		{"logo.png", 46_080},
		{"photo-equipe.png", 2_936_012},
		{"notes.txt", 2_048},
		{"backup.zip", 5_242_880},
		{"donnees.csv", 512_000},
		{"config.json", 1_024},
		{"image.jpg", 1_800_000},
		{"archive.tar.gz", 4_100_000},
	}

	var list []FileInfo
	for _, item := range sampleList {
		path := filepath.Join(dir, item.name)
		f, err := os.Create(path)
		if err != nil {
			log.Fatalf("Erreur création fichier test: %v", err)
		}

		io.CopyN(f, rand.Reader, item.size)
		f.Close()

		list = append(list, FileInfo{
			Name: item.name,
			Size: item.size,
		})
	}
	return list
}
