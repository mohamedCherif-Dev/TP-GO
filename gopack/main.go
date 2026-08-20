package main

import (
	"errors"
	"fmt"
	"gopack/format"
	"gopack/inventory"
	"os"
	"time"

	"github.com/fatih/color"
)

func donnees() []inventory.File {
	return []inventory.File{
		{Name: "main.go", Size: 4_096, Extension: ".go", Modified: time.Date(2026, 7, 1, 10, 15, 0, 0, time.UTC)},
		{Name: "utils.go", Size: 2_048, Extension: ".go", Modified: time.Date(2026, 7, 3, 16, 40, 0, 0, time.UTC)},
		{Name: "README.md", Size: 8_192, Extension: ".md", Modified: time.Date(2026, 6, 20, 9, 0, 0, 0, time.UTC)},
		{Name: "lisezmoi.md", Size: 1_024, Extension: ".md", Modified: time.Date(2026, 6, 21, 11, 30, 0, 0, time.UTC)},
		{Name: "rapport.pdf", Size: 1_254_400, Extension: ".pdf", Modified: time.Date(2026, 6, 12, 9, 30, 0, 0, time.UTC)},
		{Name: "presentation.pdf", Size: 3_670_016, Extension: ".pdf", Modified: time.Date(2026, 5, 30, 14, 0, 0, 0, time.UTC)},
		{Name: "logo.png", Size: 46_080, Extension: ".png", Modified: time.Date(2026, 4, 2, 8, 45, 0, 0, time.UTC)},
		{Name: "photo-equipe.png", Size: 2_936_012, Extension: ".png", Modified: time.Date(2026, 4, 2, 8, 50, 0, 0, time.UTC)},
		{Name: "todo.txt", Size: 512, Extension: ".txt", Modified: time.Date(2026, 7, 18, 18, 5, 0, 0, time.UTC)},
		{Name: "backup.zip", Size: 15_728_640, Extension: ".zip", Modified: time.Date(2026, 1, 15, 3, 0, 0, 0, time.UTC)},
		{Name: "archive-2025.zip", Size: 44_040_192, Extension: ".zip", Modified: time.Date(2025, 12, 31, 23, 59, 0, 0, time.UTC)},
		{Name: "config.txt", Size: 730, Extension: ".txt", Modified: time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)},
	}
}

func main() {
	fichiers := donnees()
	inventory.MarquerLesLourds(fichiers, 1_000_000)

	mode := "texte"
	extFiltre := ""

	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	if len(os.Args) > 2 {
		extFiltre = os.Args[2]
	}

	var err error
	if extFiltre != "" {
		fichiers, err = inventory.FiltrerParExtension(fichiers, extFiltre)
		if err != nil {
			switch {
			case errors.Is(err, inventory.ErrAucunResultat):
				color.Yellow("[Avertissement] %v\n", err)
				os.Exit(0)
			case errors.Is(err, inventory.ErrExtensionInvalide):
				color.Red("[Erreur] %v\n", err)
				fmt.Println("Usage: go run . [texte|json] [.ext]")
				os.Exit(1)
			default:
				color.Red("[Erreur inattendue] %v\n", err)
				os.Exit(1)
			}
		}
	}

	var fmtr format.Formatter
	switch mode {
	case "json":
		fmtr = format.JSONFormatter{}
	case "texte":
		fmtr = format.TextFormatter{Colored: true}
	default:
		color.Red("Format inconnu: %s (utilisez 'texte' ou 'json')\n", mode)
		os.Exit(1)
	}

	sortie, err := fmtr.Format(fichiers)
	if err != nil {
		color.Red("Erreur d'affichage: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(sortie)
}
