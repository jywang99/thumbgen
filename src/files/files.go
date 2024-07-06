package files

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"jy.org/thumbgen/src/config"
)

func WalkAndDo(root string, process func(string, bool), doForDir func(string) error) error {
    maxDepth := cfg.Dirs.MaxDepth

    var walk func(string, int)
    walk = func(dir string, depth int) {
        if ignoreEntry(dir) {
            return
        }
        if depth > maxDepth {
            process(dir, true)
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
            pth := filepath.Join(dir, file.Name())
            if file.IsDir() {
                walk(pth, depth + 1)
                continue
            }

            // process video files
            ext := filepath.Ext(pth)
            if !ignoreEntry(file.Name()) && len(ext) > 0 && cfg.Files.VideoExtMap[ext[1:]] {
                process(pth, false)
            }
        }
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

func MkTmpDir(fileNm string) (string, error) {
    tmpDir := path.Join(config.Config.Dirs.Temp, path.Base(fileNm))

    err := os.RemoveAll(tmpDir)
    if err != nil {
        return "", err
    }

    err = os.Mkdir(tmpDir, 0755)
    if err != nil {
        return "", err
    }
    return tmpDir, nil
}

func GetBaseName(path string, ext bool) string {
    base := filepath.Base(path)
    if !ext {
        return strings.TrimSuffix(base, filepath.Ext(base))
    }
    return base
}

func CheckFileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

