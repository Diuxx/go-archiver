// backup.go
package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// ZipFolder compresse tout le dossier source dans targetZip (créé/écrasé).
func ZipFolder(source, targetZip string) error {
	// Créer le fichier ZIP
	zf, err := os.Create(targetZip)
	if err != nil {
		return err
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	srcAbs, err := filepath.Abs(source)
	if err != nil {
		return err
	}
	parent := filepath.Dir(srcAbs) // pour inclure le nom du dossier racine dans le zip

	// Walk
	return filepath.Walk(srcAbs, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Exclusions simples (ex: fichiers temporaires)
		base := filepath.Base(path)
		if base == ".DS_Store" || base == "Thumbs.db" {
			return nil
		}

		// Chemin relatif pour le zip, incluant le dossier racine
		rel, err := filepath.Rel(parent, path)
		if err != nil {
			return err
		}

		// Dossier → on force une entrée se terminant par '/'
		if info.IsDir() {
			_, err := zw.Create(rel + "/")
			return err
		}

		// Fichier → header + copie du contenu
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		hdr.Name = rel
		hdr.Method = zip.Deflate
		// Timestamp stable
		hdr.SetModTime(info.ModTime().Round(time.Second))

		w, err := zw.CreateHeader(hdr)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(w, f); err != nil {
			return err
		}

		logger.Printf("➕ %s", rel)
		return nil
	})
}

// ClearFolder supprime tous les fichiers et sous-dossiers d'un dossier, mais garde le dossier lui-même.
func ClearFolder(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, e := range entries {
		full := filepath.Join(path, e.Name())
		err := os.RemoveAll(full)
		if err != nil {
			return fmt.Errorf("échec suppression %s: %w", full, err)
		}
	}
	return nil
}

func localTempDir() string {
	return os.TempDir() // cross-platform
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
