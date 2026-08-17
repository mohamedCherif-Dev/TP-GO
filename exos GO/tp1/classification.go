package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run classification.go <note>")
		os.Exit(1)
	}

	if note, err := strconv.ParseFloat(os.Args[1], 64); err != nil {
		fmt.Println("Erreur : la note doit être un nombre")
		os.Exit(1)
	} else {
		if note < 0 || note > 20 {
			fmt.Println("Erreur : la note doit être comprise entre 0 et 20")
			os.Exit(1)
		}

		var mention string
		switch {
		case note < 10:
			mention = "Ajourné"
		case note < 12:
			mention = "Passable"
		case note < 14:
			mention = "Assez bien"
		case note < 16:
			mention = "Bien"
		default:
			mention = "Très bien"
		}

		fmt.Printf("Note : %g/20 — Mention %s\n", note, mention)
	}
}
