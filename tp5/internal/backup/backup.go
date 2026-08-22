package backup

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Run(src, dst string) error {
	absSrc, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("erreur résolution chemin source: %w", err)
	}

	absDst, err := filepath.Abs(dst)
	if err != nil {
		return fmt.Errorf("erreur résolution chemin destination: %w", err)
	}

	if strings.HasPrefix(absDst, absSrc) {
		return fmt.Errorf("la destination %q ne peut pas être située dans la source %q", absDst, absSrc)
	}

	copiedFiles := 0

	err = filepath.WalkDir(absSrc, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(absSrc, path)
		if err != nil {
			return fmt.Errorf("erreur rel path: %w", err)
		}
		targetPath := filepath.Join(absDst, relPath)

		info, err := d.Info()
		if err != nil {
			return err
		}

		if d.IsDir() {
			return os.MkdirAll(targetPath, info.Mode().Perm())
		}

		if d.Type().IsRegular() {
			if err := copyFile(path, targetPath, info.Mode().Perm()); err != nil {
				return err
			}
			copiedFiles++
			fmt.Printf("\rCopie en cours... %d fichier(s)", copiedFiles)
		}

		return nil
	})

	fmt.Println()
	if err != nil {
		return fmt.Errorf("erreur lors de la sauvegarde: %w", err)
	}

	fmt.Printf("Sauvegarde terminée avec succès : %d fichier(s) copiés.\n", copiedFiles)
	return nil
}

func copyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("ouverture %s: %w", src, err)
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("création %s: %w", dst, err)
	}

	_, err = io.Copy(out, in)
	closeErr := out.Close()

	if err != nil {
		return fmt.Errorf("écriture %s: %w", dst, err)
	}
	if closeErr != nil {
		return fmt.Errorf("fermeture %s: %w", dst, closeErr)
	}

	if err := os.Chmod(dst, perm); err != nil {
		return fmt.Errorf("application permissions %s: %w", dst, err)
	}

	return nil
}
