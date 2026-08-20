package format

import (
	"encoding/json"
	"fmt"
	"gopack/inventory"
)

type JSONFormatter struct{}

func (j JSONFormatter) Format(fichiers []inventory.File) (string, error) {
	if len(fichiers) == 0 {
		return "", ErrRienAFormater
	}
	bytes, err := json.MarshalIndent(fichiers, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encodage JSON %w", err)
	}
	return string(bytes), nil
}
