// TP2 — Inventaire de fichiers en mémoire
package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"time"
)

// ────────────────────────────────────────────────────────────────────
// Étape 1 — Définition de la struct File
// ────────────────────────────────────────────────────────────────────
type File struct {
	Name      string    // nom du fichier, ex : "rapport.pdf"
	Size      int64     // taille en octets
	Extension string    // extension avec le point, ex : ".pdf"
	Modified  time.Time // date de dernière modification
	Tag       string    // étiquette libre, vide par défaut
}

// ────────────────────────────────────────────────────────────────────
// Étape 2 — Les données en dur : 12 fichiers simulant un dossier projet
// ────────────────────────────────────────────────────────────────────
var inventaire = []File{
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

// ────────────────────────────────────────────────────────────────────
// Étape 3 — Filtrage avec retours multiples (résultat, erreur)
// ────────────────────────────────────────────────────────────────────
func filtrerParExtension(fichiers []File, ext string) ([]File, error) {
	if ext == "" || ext[0] != '.' {
		return nil, errors.New("l'extension doit être non vide et commencer par un point")
	}
	var res []File
	for _, f := range fichiers {
		if f.Extension == ext {
			res = append(res, f)
		}
	}
	return res, nil
}

func filtrerParTailleMin(fichiers []File, tailleMin int64) ([]File, error) {
	if tailleMin < 0 {
		return nil, errors.New("la taille minimale ne peut pas être négative")
	}
	var res []File
	for _, f := range fichiers {
		if f.Size >= tailleMin {
			res = append(res, f)
		}
	}
	return res, nil
}

// ────────────────────────────────────────────────────────────────────
// Étape 4 — Statistiques par extension
// ────────────────────────────────────────────────────────────────────
type Stats struct {
	Nombre       int
	TailleTotale int64
}

func statistiquesParExtension(fichiers []File) map[string]Stats {
	stats := make(map[string]Stats)
	for _, f := range fichiers {
		s := stats[f.Extension] // Zero value {0, 0} si inexistante
		s.Nombre++
		s.TailleTotale += f.Size
		stats[f.Extension] = s
	}
	return stats
}

// ────────────────────────────────────────────────────────────────────
// Étape 5 — Tri par taille décroissante (sans modifier l'original)
// ────────────────────────────────────────────────────────────────────
func trierParTaille(fichiers []File) []File {
	tries := make([]File, len(fichiers))
	copy(tries, fichiers)
	sort.Slice(tries, func(i, j int) bool {
		return tries[i].Size > tries[j].Size
	})
	return tries
}

// ────────────────────────────────────────────────────────────────────
// Étape 6 — Affichage formaté
// ────────────────────────────────────────────────────────────────────
func tailleLisible(octets int64) string {
	switch {
	case octets >= 1024*1024:
		return fmt.Sprintf("%.1f Mo", float64(octets)/(1024*1024))
	case octets >= 1024:
		return fmt.Sprintf("%.1f Ko", float64(octets)/1024)
	default:
		return fmt.Sprintf("%d o", octets)
	}
}

func afficher(fichiers []File) {
	fmt.Printf("%-25s %-10s %-12s %-10s\n", "NOM", "TAILLE", "MODIFIÉ", "TAG")
	for _, f := range fichiers {
		fmt.Printf("%-25s %-10s %-12s %-10s\n",
			f.Name,
			f.TailleLisible(),
			f.Modified.Format("2006-01-02"),
			f.Tag,
		)
	}
}

// ────────────────────────────────────────────────────────────────────
// Étape 7 — Mini-challenge : méthodes & pointer receiver
// ────────────────────────────────────────────────────────────────────
// Value receiver
func (f File) TailleLisible() string {
	return tailleLisible(f.Size)
}

// Pointer receiver
func (f *File) Renommer(nouveauNom string) {
	f.Name = nouveauNom
	f.Extension = filepath.Ext(nouveauNom)
}

// Pointer receiver
func (f *File) Marquer(tag string) {
	f.Tag = tag
}

func marquerLesLourds(fichiers []File, tailleMin int64) int {
	compteur := 0
	for i := range fichiers {
		if fichiers[i].Size >= tailleMin {
			fichiers[i].Marquer("a-archiver")
			compteur++
		}
	}
	return compteur
}

// ────────────────────────────────────────────────────────────────────
// Main — Déroulement du scénario
// ────────────────────────────────────────────────────────────────────
func main() {
	fmt.Println("=== TP2 — inventaire de fichiers ===")

	// --- Étape 1 ---
	fmt.Println("\n--- Étape 1 : Struct File ---")
	fTest := File{Name: "test.txt", Size: 100, Extension: ".txt", Modified: time.Now()}
	fmt.Println("Affichage simple  :", fTest)
	fmt.Printf("Affichage détaillé : %+v\n", fTest)

	// --- Étape 2 ---
	fmt.Println("\n--- Étape 2 : Slice et Append ---")
	fmt.Println("Nombre initial de fichiers :", len(inventaire))
	inventaire = append(inventaire, File{
		Name:      "notes-perso.txt",
		Size:      256,
		Extension: ".txt",
		Modified:  time.Now(),
	})
	fmt.Println("Nombre après append        :", len(inventaire))

	// --- Étape 3 ---
	fmt.Println("\n--- Étape 3 : Filtrage & Gestion des erreurs ---")
	goFiles, err := filtrerParExtension(inventaire, ".go")
	if err != nil {
		fmt.Println("Erreur filtrage .go :", err)
	} else {
		fmt.Printf("Fichiers .go trouvés : %d\n", len(goFiles))
	}

	// Test d'erreur intentionnel
	_, err = filtrerParExtension(inventaire, "go")
	if err != nil {
		fmt.Println("Erreur interceptée (sans point) :", err)
	}

	// Chaînage des filtres (.png >= 1 Mo)
	pngs, _ := filtrerParExtension(inventaire, ".png")
	pngsLourds, _ := filtrerParTailleMin(pngs, 1_000_000)
	fmt.Printf("PNG > 1 Mo trouvés : %d (%s)\n", len(pngsLourds), pngsLourds[0].Name)

	// --- Étape 4 ---
	fmt.Println("\n--- Étape 4 : Statistiques ---")
	stats := statistiquesParExtension(inventaire)
	if s, ok := stats[".zip"]; ok {
		fmt.Printf("Les .zip (%d fichier(s)) pèsent au total %s\n", s.Nombre, tailleLisible(s.TailleTotale))
	}

	// --- Étape 5 ---
	fmt.Println("\n--- Étape 5 : Tri par taille (Top 3) ---")
	top3 := trierParTaille(inventaire)
	afficher(top3[:3])

	// --- Étape 6 & 7 ---
	fmt.Println("\n--- Étape 7 : Renommage, Marquage et Affichage Final ---")
	// Renommage de todo.txt -> todo-2026.md
	for i := range inventaire {
		if inventaire[i].Name == "todo.txt" {
			inventaire[i].Renommer("todo-2026.md")
		}
	}

	// Marquage des fichiers >= 1 Mo
	nbMarques := marquerLesLourds(inventaire, 1_000_000)
	fmt.Printf("Nombre de fichiers marqués 'a-archiver' : %d\n\n", nbMarques)

	// Affichage complet formaté
	afficher(inventaire)
}
