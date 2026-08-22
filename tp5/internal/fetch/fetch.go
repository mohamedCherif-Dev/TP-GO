package fetch

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type result struct {
	url string
	err error
}

func downloadFile(url, outDir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statut HTTP %d", resp.StatusCode)
	}

	filename := filepath.Base(url)
	if filename == "" || filename == "." || filename == "/" {
		filename = "downloaded_file"
	}
	targetPath := filepath.Join(outDir, filename)

	out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func Run(urls []string, outDir string, workers int) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("création du dossier %s: %w", outDir, err)
	}

	total := len(urls)
	sem := make(chan struct{}, workers)
	results := make(chan result, total)
	var wg sync.WaitGroup

	var completed int32

	for _, u := range urls {
		wg.Add(1)
		go func(targetURL string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			err := downloadFile(targetURL, outDir)
			results <- result{url: targetURL, err: err}

			done := atomic.AddInt32(&completed, 1)
			fmt.Printf("\rProgrès: [%d/%d]", done, total)
		}(u)
	}

	wg.Wait()
	close(results)
	fmt.Println()

	var hasError bool
	var successCount int

	for res := range results {
		if res.err != nil {
			fmt.Fprintf(os.Stderr, "Échec pour %s : %v\n", res.url, res.err)
			hasError = true
		} else {
			successCount++
		}
	}

	fmt.Printf("Bilan : %d/%d fichier(s) téléchargé(s) avec succès.\n", successCount, total)

	if hasError {
		return errors.New("un ou plusieurs téléchargements ont échoué")
	}

	return nil
}
