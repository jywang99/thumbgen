package files

import (
	"path/filepath"
	"strings"

	"jy.org/thumbgen/src/config"
	"jy.org/thumbgen/src/logging"
)

var logger = logging.Logger
var cfg = config.Config

func ignoreEntry(path string) bool {
    ignore := cfg.Dirs.IgnoreMap

    if !cfg.Files.DotFiles && strings.HasPrefix(filepath.Base(path), ".") {
        logger.INFO.Printf("Skipping dot directory/file: %v\n", path)
        return true
    }

    if ignore[filepath.Base(path)] {
        logger.INFO.Printf("Skipping ignored directory: %v\n", path)
        return true
    }

    return false
}

func getExt(path string) string {
    ext := filepath.Ext(path)
    if ext == "" {
        return ""
    }
    return strings.ToLower(ext[1:])
}

