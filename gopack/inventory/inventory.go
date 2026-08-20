package inventory

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Erreurs sentinelles
var (
	ErrExtensionInvalide = errors.New("extension invalide (attendu \".ext\")")
	ErrTailleNegative    = errors.New("taille minimale négative")
	ErrAucunResultat     = errors.New("aucun fichier ne correspond")
)

type File struct {
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Extension string    `json:"extension"`
	Modified  time.Time `json:"modified"`
	Tag       string    `json:"tag,omitempty"`
}

type Stats struct {
	Nombre       int
	TailleTotale int64
}

func (f File) TailleLisible() string {
	return TailleLisible(f.Size)
}

func (f *File) Renommer(nouveauNom string) {
	f.Name = nouveauNom
	f.Extension = filepath.Ext(nouveauNom)
}

func (f *File) Marquer(tag string) {
	f.Tag = tag
}

func TailleLisible(octets int64) string {
	switch {
	case octets >= 1024*1024:
		return fmt.Sprintf("%.1f Mo", float64(octets)/(1024*1024))
	case octets >= 1024:
		return fmt.Sprintf("%.1f Ko", float64(octets)/1024)
	default:
		return fmt.Sprintf("%d o", octets)
	}
}

func FiltrerParExtension(fichiers []File, ext string) ([]File, error) {
	if ext == "" || !strings.HasPrefix(ext, ".") {
		return nil, fmt.Errorf("filtre extension %q: %w", ext, ErrExtensionInvalide)
	}
	var res []File
	for _, f := range fichiers {
		if f.Extension == ext {
			res = append(res, f)
		}
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("filtre extension %q: %w", ext, ErrAucunResultat)
	}
	return res, nil
}

func FiltrerParTailleMin(fichiers []File, tailleMin int64) ([]File, error) {
	if tailleMin < 0 {
		return nil, fmt.Errorf("filtre taille %d: %w", tailleMin, ErrTailleNegative)
	}
	var res []File
	for _, f := range fichiers {
		if f.Size >= tailleMin {
			res = append(res, f)
		}
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("filtre taille %d: %w", tailleMin, ErrAucunResultat)
	}
	return res, nil
}

func StatistiquesParExtension(fichiers []File) map[string]Stats {
	stats := make(map[string]Stats)
	for _, f := range fichiers {
		s := stats[f.Extension]
		s.Nombre++
		s.TailleTotale += f.Size
		stats[f.Extension] = s
	}
	return stats
}

func TrierParTaille(fichiers []File) []File {
	tries := make([]File, len(fichiers))
	copy(tries, fichiers)
	sort.Slice(tries, func(i, j int) bool {
		return tries[i].Size > tries[j].Size
	})
	return tries
}

func MarquerLesLourds(fichiers []File, tailleMin int64) int {
	compteur := 0
	for i := range fichiers {
		if fichiers[i].Size >= tailleMin {
			fichiers[i].Marquer("a-archiver")
			compteur++
		}
	}
	return compteur
}
