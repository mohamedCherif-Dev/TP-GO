package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	serverURL = "http://localhost:8080"
	outDir    = "downloads"
)

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type Result struct {
	FileName string
	Size     int64
	Duration time.Duration
	Err      error
}

func fetchList() ([]FileInfo, error) {
	resp, err := http.Get(serverURL + "/files")
	if err != nil {
		return nil, fmt.Errorf("requête HTTP impossible: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("statut HTTP inattendu: %s", resp.Status)
	}

	var files []FileInfo
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("décodage JSON impossible: %w", err)
	}

	return files, nil
}

func downloadFile(name string) (int64, error) {
	resp, err := http.Get(serverURL + "/files/" + name)
	if err != nil {
		return 0, fmt.Errorf("erreur HTTP GET: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("statut HTTP %d (%s)", resp.StatusCode, resp.Status)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return 0, fmt.Errorf("impossible de créer le dossier %s: %w", outDir, err)
	}

	filePath := filepath.Join(outDir, name)
	out, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return 0, fmt.Errorf("impossible de créer le fichier %s: %w", filePath, err)
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de l'écriture: %w", err)
	}

	return written, nil
}

func main() {
	maxConcurrent := flag.Int("n", 4, "Nombre maximum de téléchargements simultanés")
	flag.Parse()

	fmt.Println("=== Client de Téléchargement Concurrent ===")

	files, err := fetchList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de la récupération de la liste: %v\n", err)
		os.Exit(1)
	}

	totalFiles := len(files)
	fmt.Printf("Liste récupérée : %d fichiers trouvés (Concurrence max: %d)\n\n", totalFiles, *maxConcurrent)

	start := time.Now()

	results := make(chan Result, totalFiles)
	sem := make(chan struct{}, *maxConcurrent)
	var wg sync.WaitGroup

	for _, f := range files {
		wg.Add(1)
		go func(file FileInfo) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			t0 := time.Now()
			written, err := downloadFile(file.Name)
			duration := time.Since(t0)

			results <- Result{
				FileName: file.Name,
				Size:     written,
				Duration: duration,
				Err:      err,
			}
		}(f)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var successCount, failureCount int
	completed := 0

	for res := range results {
		completed++
		if res.Err != nil {
			failureCount++
			fmt.Fprintf(os.Stderr, "[%d/%d]  Erreur sur %s : %v\n", completed, totalFiles, res.FileName, res.Err)
		} else {
			successCount++
			fmt.Printf("[%d/%d]  %-20s (%d octets) téléchargé en %v\n", completed, totalFiles, res.FileName, res.Size, res.Duration.Round(time.Millisecond))
		}
	}

	totalDuration := time.Since(start)

	fmt.Printf("\n=== Bilan du téléchargement ===\n")
	fmt.Printf("Temps total    : %v\n", totalDuration.Round(time.Millisecond))
	fmt.Printf("Succès         : %d\n", successCount)
	fmt.Printf("Échecs         : %d\n", failureCount)

	if failureCount > 0 {
		os.Exit(1)
	}
}
