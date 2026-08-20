package format

import (
	"errors"
	"gopack/inventory"
)

var ErrRienAFormater = errors.New("inventaire vide, rien à formater")

type Formatter interface {
	Format(fichiers []inventory.File) (string, error)
}
