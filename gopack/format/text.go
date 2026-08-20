package format

import (
	"fmt"
	"gopack/inventory"
	"strings"

	"github.com/fatih/color"
)

type TextFormatter struct {
	Colored bool
}

func (t TextFormatter) Format(fichiers []inventory.File) (string, error) {
	if len(fichiers) == 0 {
		return "", ErrRienAFormater
	}

	var b strings.Builder
	header := fmt.Sprintf("%-25s %-10s %-12s %-10s", "NOM", "TAILLE", "MODIFIÉ", "TAG")
	if t.Colored {
		header = color.New(color.FgCyan, color.Bold).Sprint(header)
	}
	b.WriteString(header + "\n")

	for _, f := range fichiers {
		tagStr := f.Tag
		if t.Colored && tagStr != "" {
			tagStr = color.New(color.FgYellow).Sprint(tagStr)
		}
		fmt.Fprintf(&b, "%-25s %-10s %-12s %-10s\n",
			f.Name,
			f.TailleLisible(),
			f.Modified.Format("2006-01-02"),
			tagStr,
		)
	}
	return b.String(), nil
}
