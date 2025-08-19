// logger.go
package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logger *log.Logger

// MustExecDir : dossier où se trouve le binaire (pas le cwd).
func MustExecDir() string {
	exe, err := os.Executable()
	if err != nil {
		wd, _ := os.Getwd()
		return wd
	}
	return filepath.Dir(exe)
}

// MustInitLogger crée logs/<prefix>-YYYYMMDD.log à côté du binaire.
// Et (optionnel) duplique sur stdout si verbose.
func MustInitLogger(execDir, prefix string, verbose bool) *os.File {
	logsDir := filepath.Join(execDir, "logs")
	_ = os.MkdirAll(logsDir, 0o755)

	name := prefix + "-" + time.Now().Format("20060102") + ".log"
	full := filepath.Join(logsDir, name)

	f, err := os.OpenFile(full, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}

	var out io.Writer = f
	if verbose {
		out = io.MultiWriter(os.Stdout, f)
	}

	logger = log.New(out, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Printf("📝 log file: %s", full)
	return f
}
