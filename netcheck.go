// netcheck.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// EnsureMounted vérifie qu'un dossier est accessible.
// - Detecte si c'est probablement un FS réseau (Unix: /proc/mounts ; Windows: UNC).
// - Vérifie existence + droits (lecture pour origin, écriture pour dest).
// - Renvoie une erreur claire si le partage n'est pas monté ou inaccessible.
func EnsureMounted(path string, wantWrite bool) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("chemin invalide: %w", err)
	}

	isNet, fsType := isLikelyNetworkPath(abs)

	// Existence + type dossier
	info, err := os.Stat(abs)
	if err != nil {
		if isNet {
			return fmt.Errorf("chemin réseau introuvable (non monté ?): %s (fs=%s) : %w", abs, fsType, err)
		}
		return fmt.Errorf("chemin introuvable: %s : %w", abs, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("doit être un dossier: %s", abs)
	}

	// Test lecture (ReadDir)
	if _, err := os.ReadDir(abs); err != nil {
		if isNet {
			return fmt.Errorf("dossier réseau illisible (montage/permissions ?) %s (fs=%s) : %w", abs, fsType, err)
		}
		return fmt.Errorf("dossier illisible: %s : %w", abs, err)
	}

	// Test écriture si demandé
	if wantWrite {
		tmp := filepath.Join(abs, ".write_test.tmp")
		f, err := os.Create(tmp)
		if err != nil {
			if isNet {
				return fmt.Errorf("destination réseau non inscriptible (montage/credentials ?) %s (fs=%s) : %w", abs, fsType, err)
			}
			return fmt.Errorf("destination non inscriptible: %s : %w", abs, err)
		}
		f.Close()
		_ = os.Remove(tmp)
	}

	return nil
}

// isLikelyNetworkPath retourne (true, fstype) si le chemin semble être un montage réseau.
// Unix: lit /proc/mounts et regarde le fstype (nfs, cifs, smbfs, fuse.sshfs, afpfs…).
// Windows: UNC (\\serveur\partage) = réseau. Les lecteurs mappés sont détectés via tests d'accès.
func isLikelyNetworkPath(path string) (bool, string) {
	if runtime.GOOS == "windows" {
		// UNC path → \\server\share\...
		clean := filepath.Clean(path)
		if strings.HasPrefix(clean, `\\`) {
			return true, "UNC"
		}
		// Pour un lecteur mappé (Z:\), on ne sait pas trivially en stdlib → on s'appuie sur les tests d'accès.
		return false, ""
	}

	// Unix-like : /proc/mounts
	const mounts = "/proc/mounts"
	f, err := os.Open(mounts)
	if err != nil {
		// Pas de /proc/mounts (ex: macOS sans compat) → fallback heuristique simple
		// Heuristique: si ça ressemble à /Volumes/NAS/... (macOS) ou contient "smb"/"nfs" dans le chemin
		lower := strings.ToLower(path)
		if strings.Contains(lower, "/volumes/") || strings.Contains(lower, "smb") || strings.Contains(lower, "nfs") {
			return true, "heuristic"
		}
		return false, ""
	}
	defer f.Close()

	type mountEntry struct{ target, fstype string }

	var best mountEntry
	var have bool
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		// /proc/mounts format: src target fstype options ...
		fields := strings.Fields(sc.Text())
		if len(fields) < 3 {
			continue
		}
		target := fields[1]
		fstype := fields[2]

		// On cherche le point de montage le plus spécifique (préfixe le plus long)
		if isPathPrefix(path, target) {
			if !have || len(target) > len(best.target) {
				best = mountEntry{target: target, fstype: fstype}
				have = true
			}
		}
	}
	if have {
		netTypes := map[string]bool{
			"nfs": true, "nfs4": true,
			"cifs": true, "smbfs": true, "smb3": true, "smb2": true,
			"fuse.sshfs": true, "sshfs": true,
			"afpfs": true,
		}
		if netTypes[strings.ToLower(best.fstype)] {
			return true, best.fstype
		}
	}
	return false, ""
}

func isPathPrefix(p, prefix string) bool {
	p = filepath.Clean(p)
	prefix = filepath.Clean(prefix)
	if p == prefix {
		return true
	}
	rel, err := filepath.Rel(prefix, p)
	return err == nil && !strings.HasPrefix(rel, "..")
}
