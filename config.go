// config.go
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type AppConfig struct {
	LogPrefix string
	Verbose   bool
	Origin    string
	Dest      string
}

func ParseFlags() AppConfig {
	var cfg AppConfig
	flag.StringVar(&cfg.LogPrefix, "log-prefix", "backup", "Préfixe du fichier log")
	flag.BoolVar(&cfg.Verbose, "v", true, "Affiche aussi les logs sur stdout")
	flag.StringVar(&cfg.Origin, "origin", "", "Dossier à sauvegarder (obligatoire)")
	flag.StringVar(&cfg.Dest, "dest", "", "Dossier de dépôt du ZIP (obligatoire)")
	flag.Parse()
	return cfg
}

func ValidateConfig(cfg AppConfig) error {
	if cfg.Origin == "" || cfg.Dest == "" {
		return errors.New("les flags --origin et --dest sont obligatoires")
	}

	// Création dest si manquant (avant tests d'accessibilité)
	if err := os.MkdirAll(cfg.Dest, 0o755); err != nil {
		return fmt.Errorf("impossible de créer dest: %w", err)
	}

	// Vérifs génériques + réseau: origin (lecture), dest (écriture)
	if err := EnsureMounted(cfg.Origin, false); err != nil {
		return fmt.Errorf("origin invalide: %w", err)
	}
	if err := EnsureMounted(cfg.Dest, true); err != nil {
		return fmt.Errorf("dest invalide: %w", err)
	}

	// Bonus: normaliser chemins
	if abs, err := filepath.Abs(cfg.Origin); err == nil {
		cfg.Origin = abs
	}
	if abs, err := filepath.Abs(cfg.Dest); err == nil {
		cfg.Dest = abs
	}
	return nil
}
