// main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {

	// read flags
	cfg := ParseFlags()

	// validate required flags
	if err := ValidateConfig(cfg); err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}

	execDir := MustExecDir()
	logFile := MustInitLogger(execDir, cfg.LogPrefix, cfg.Verbose)
	defer logFile.Close()

	logger.Printf("🚀 Backup start | origin=%s dest=%s", cfg.Origin, cfg.Dest)

	ts := time.Now().Format("20060102-150405")
	// zipName := fmt.Sprintf("backup-%s.zip", ts)
	// zipPath := filepath.Join(cfg.Dest, zipName)
	tempDir := localTempDir()
	localZip := filepath.Join(tempDir, "backup-"+ts+".zip")

	logger.Printf("📦 Création du ZIP local: %s", localZip)
	if err := ZipFolder(cfg.Origin, localZip); err != nil {
		logger.Fatalf("❌ Erreur backup local: %v", err)
		os.Exit(1)
	}

	// 2) Copier vers destination
	zipDest := filepath.Join(cfg.Dest, filepath.Base(localZip))
	logger.Printf("📤 Copie vers destination: %s", zipDest)
	if err := copyFile(localZip, zipDest); err != nil {
		logger.Fatalf("❌ Erreur copie vers destination: %v", err)
		os.Exit(1)
	}

	logger.Printf("✅ Backup ok → %s", zipDest)

	_ = os.Remove(localZip)

	// 🆕 Vider le contenu du dossier origin
	if err := ClearFolder(cfg.Origin); err != nil {
		logger.Printf("⚠️ Impossible de vider origin: %v", err)
	} else {
		logger.Printf("🗑 Origin vidé: %s", cfg.Origin)
	}
}
