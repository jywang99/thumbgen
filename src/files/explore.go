package files

import (
	"os"
	"path/filepath"
)

type Explorer struct {
    Dir string
    VidFiles []string
    ImgFiles []string
}

func NewExplorer(dir string) (*Explorer, error) {
    e := &Explorer{
        Dir: dir, 
    }

    err := e.getFiles()
    if err != nil {
        logger.ERROR.Printf("Error when reading directory: %v\n", err)
        return nil, err
    }

    return e, nil
}

func (e *Explorer) getFiles() error {
    logger.INFO.Printf("Exploring directory: %v\n", e.Dir)

    vidExts := cfg.Files.VideoExtMap
    imgExts := cfg.Files.ImageExtMap

	err := filepath.Walk(e.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
        if info.IsDir() || ignoreEntry(path) {
            return nil
        }

        // no extension
        ext := getExt(path)
        if ext == "" {
            return nil
        }

        // add to list depending on file type
        if vidExts[ext] {
            e.VidFiles = append(e.VidFiles, path)
        } else if imgExts[ext] {
            e.ImgFiles = append(e.ImgFiles, path)
        }

        return nil
	})
    if err != nil {
        logger.ERROR.Printf("Error when reading directory")
        return err
    }

    logger.INFO.Printf("Found %v video files and %v image files\n", len(e.VidFiles), len(e.ImgFiles))
	return nil
}

type FileType int
const (
    Video FileType = iota
    Image
)

func (e *Explorer) GetFileCount(ft FileType) int {
    switch ft {
    case Video:
        return len(e.VidFiles)
    case Image:
        return len(e.ImgFiles)
    }
    return 0
}

