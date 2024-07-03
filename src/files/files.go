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

func WalkAndDo(root string, doForFile, doForLeafDir func(string), doForDir func(string) error) error {
    maxDepth := cfg.Dirs.MaxDepth
    ignore := cfg.Dirs.IgnoreMap
    dots := cfg.Files.DotFiles

    ignoreDot := func(path string) bool {
        return !dots && strings.HasPrefix(filepath.Base(path), ".")
    }

    var walk func(string, int)
    walk = func(dir string, depth int) {
        if ignoreDot(dir) {
            logger.INFO.Printf("Skipping dot directory: %v\n", dir)
            return
        }
        if ignore[filepath.Base(dir)] {
            logger.INFO.Printf("Skipping ignored directory: %v\n", dir)
            return
        }
        if depth > maxDepth {
            doForLeafDir(dir)
            return
        }

        // process this dir
        err := doForDir(dir)
        if err != nil {
            return
        }

        files, err := os.ReadDir(dir)
        if err != nil {
            logger.ERROR.Printf("Error when reading directory: %v\n", err)
            return
        }

        // directory contents
        for _, file := range files {
            // descend into subdirs
            dir := filepath.Join(dir, file.Name())
            if file.IsDir() {
                walk(dir, depth + 1)
                continue
            }

            // process files
            ext := filepath.Ext(dir)
            if !ignoreDot(file.Name()) && len(ext) > 0 && cfg.Files.VideoExtMap[ext[1:]] {
                doForFile(dir)
            }
        }

        return
    }
    walk(root, 0)
    return nil
}

func GetTargetDir(dir string, strip bool) (string, error) {
    // get relative path from input dir
    origBase := cfg.Dirs.Input
    rel, err := filepath.Rel(origBase, dir)
    if err != nil {
        logger.ERROR.Printf("Error when getting relative path for %v: %v\n", dir, err)
        return "", err
    }

    // strip last filename/dirname
    if strip {
        rel = filepath.Dir(rel)
    }

    // join with output dir
    targetBase := cfg.Dirs.Output
    targetDir := filepath.Join(targetBase, rel)
    return targetDir, nil
}

func MkTargetDir(dir string) (string, error) {
    targetDir, err := GetTargetDir(dir, false)
    if err != nil {
        return targetDir, err
    }

    // check if dir exists
    if _, err := os.Stat(targetDir); err == nil {
        logger.INFO.Printf("Target dir already exists: %v\n", targetDir)
        return targetDir, nil
    }

    // create dir
    logger.INFO.Printf("Creating target dir: %v\n", targetDir)
    err = os.MkdirAll(targetDir, os.ModePerm)
    if err != nil {
        logger.ERROR.Printf("Error when creating target dir %v: %v\n", targetDir, err)
        return targetDir, err
    }

    return targetDir, nil
}

func GetBaseName(path string, ext bool) string {
    base := filepath.Base(path)
    if !ext {
        return strings.TrimSuffix(base, filepath.Ext(base))
    }
    return base
}

