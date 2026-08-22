package scan

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

type ExtStat struct {
	Count int
	Size  int64
}

type Report struct {
	TotalFiles int64
	TotalSize  int64
	Stats      map[string]*ExtStat
}

func Analyze(root string) (*Report, error) {
	report := &Report{
		Stats: make(map[string]*ExtStat),
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == "" {
			ext = "(sans extension)"
		}

		stat, exists := report.Stats[ext]
		if !exists {
			stat = &ExtStat{}
			report.Stats[ext] = stat
		}

		stat.Count++
		stat.Size += info.Size()
		report.TotalFiles++
		report.TotalSize += info.Size()

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("erreur d'analyse du dossier %s: %w", root, err)
	}

	return report, nil
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func PrintReport(report *Report) {
	type kv struct {
		Ext  string
		Stat *ExtStat
	}

	var sorted []kv
	for k, v := range report.Stats {
		sorted = append(sorted, kv{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Stat.Size > sorted[j].Stat.Size
	})

	fmt.Printf("%-20s %-12s %-15s\n", "EXTENSION", "FICHIERS", "TAILLE CUMULÉE")
	fmt.Println(strings.Repeat("-", 50))
	for _, item := range sorted {
		fmt.Printf("%-20s %-12d %-15s\n", item.Ext, item.Stat.Count, formatSize(item.Stat.Size))
	}
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("%-20s %-12d %-15s\n", "TOTAL", report.TotalFiles, formatSize(report.TotalSize))
}

func Run(dir string) error {
	report, err := Analyze(dir)
	if err != nil {
		return err
	}
	PrintReport(report)
	return nil
}
