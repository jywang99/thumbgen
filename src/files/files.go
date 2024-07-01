package files

import (
	"os"
	"path/filepath"
	"strings"

	"jy.org/videop/src/config"
	"jy.org/videop/src/logging"
)

var logger = logging.Logger
var cfg = config.Config

func WalkAndDo(path string, doForFile func(string), doForDir func(string) error) error {
    ignore := cfg.Dirs.Ignore
    targets := cfg.Files.TargetExts
    err := filepath.Walk(config.Config.Dirs.Input, func(path string, info os.FileInfo, err error) error {
		if err != nil {
            logger.ERROR.Printf("Error when walking through the directory: %v\n", err)
			return nil
		}

        // ignore dot files
        if !cfg.Files.DotFiles && strings.HasPrefix(filepath.Base(path), ".") {
            logger.INFO.Printf("Skipping dot file: %v\n", path)
            if info.IsDir() {
                return filepath.SkipDir
            }
            return nil
        }

        if info.IsDir() {
            if ignore[filepath.Base(path)] {
                logger.INFO.Printf("Skipping directory: %v\n", path)
                return filepath.SkipDir
            }
            return doForDir(path)
        }

        ext := filepath.Ext(path)
        if len(ext) > 0 && targets[ext[1:]] {
            logger.INFO.Printf("Processing file: %v\n", path)
            doForFile(path)
        }

		return nil
	})
    return err
}

func GetTargetFile(file string) (string, error) {
    // get relative path from input dir
    origBase := cfg.Dirs.Input
    rel, err := filepath.Rel(origBase, file)
    if err != nil {
        logger.ERROR.Printf("Error when getting relative path for %v: %v\n", file, err)
        return "", err
    }

    // change extension to .gif
    ext := filepath.Ext(file)
    base := strings.TrimSuffix(rel, ext)
    rel = base + ".gif"

    // join with output dir
    targetBase := cfg.Dirs.Output
    return filepath.Join(targetBase, rel), nil
}

func MkTargetDir(dir string) error {
    // get relative path from input dir
    origBase := cfg.Dirs.Input
    rel, err := filepath.Rel(origBase, dir)
    if err != nil {
        logger.ERROR.Printf("Error when getting relative path for %v: %v\n", dir, err)
        return err
    }

    // join with output dir
    targetBase := cfg.Dirs.Output
    targetDir := filepath.Join(targetBase, rel)
    if _, err := os.Stat(targetDir); err == nil {
        // dir exists, done
        return nil
    }

    // create dir
    err = os.MkdirAll(targetDir, os.ModePerm)
    if err != nil {
        logger.ERROR.Printf("Error when creating target dir %v: %v\n", targetDir, err)
        return err
    }

    return nil
}
