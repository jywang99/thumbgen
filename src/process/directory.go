package process

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"jy.org/thumbgen/src/files"
	"jy.org/thumbgen/src/process/cli"
)

type DirProcessor struct{
    exp *files.Explorer
    TargetDir string
}

func NewDirProcessor(dir string, tdir string) (*DirProcessor, error) {
    // scan source dir
    exp, err := files.NewExplorer(dir)
    if err != nil {
        return nil, err
    }
    return &DirProcessor{
        exp: exp,
        TargetDir: tdir,
    }, nil
}

func (dp *DirProcessor) GenPreviewGif() error {
    outFile := dp.getTargetFile("gif")
    if _, err := os.Stat(outFile); err == nil {
        return nil
    }

    vidCnt := dp.exp.GetFileCount(files.Video)
    imgCnt := dp.exp.GetFileCount(files.Image)

    if vidCnt == 0 && imgCnt == 0 {
        // TODO default gif
        logger.INFO.Printf("No video or image files found in %v\n", dp.exp.Dir)
        return nil
    }

    // make tmp dir
    tmpDir, err := files.MkTmpDir(dp.exp.Dir)
    if err != nil {
        return err
    }
    defer os.RemoveAll(tmpDir)

    if vidCnt > 0 {
        return dp.genGifFromVideos(tmpDir, outFile)
    } 
    return dp.genGifFromImages(tmpDir, outFile)
}

func (dp *DirProcessor) GenPreviewImg() error {
    outFile := dp.getTargetFile("png")
    if _, err := os.Stat(outFile); err == nil {
        return nil
    }

    vidCnt := dp.exp.GetFileCount(files.Video)
    imgCnt := dp.exp.GetFileCount(files.Image)

    if vidCnt == 0 && imgCnt == 0 {
        // TODO default png
        logger.INFO.Printf("No video or image files found in %v\n", dp.exp.Dir)
        return nil
    }

    // if at least one image, resize and use it as thumbnail
    if imgCnt > 0 {
        ipath := dp.exp.ImgFiles[0]
        return cli.ResizeImgTo(ipath, outFile)
    }
    // if not, use middle frame from a random video
    vpath := dp.exp.VidFiles[0]
    dur, err := cli.GetVidDuration(vpath)
    if err != nil {
        return err
    }
    return cli.GetVidFrame(vpath, outFile, dur/2)
}

func (dp *DirProcessor) getTargetFile(ext string) string {
    return path.Join(dp.TargetDir, fmt.Sprintf("%v.%v", filepath.Base(dp.exp.Dir), ext))
}

func (dp *DirProcessor) genGifFromVideos(tmpDir, outFile string) error {
    // generate gif for each video
    // process at most MaxCuts videos
    vids := dp.exp.VidFiles
    i := 0
    for i<cfg.Ffmpeg.MaxCuts && i<len(vids) { // TODO multiple cuts from a single video, if not enough videos
        vpath := vids[i]
        dur, err := cli.GetVidDuration(vpath)
        if err != nil {
            continue
        }
        gif := path.Join(tmpDir, fmt.Sprintf("range%v.gif", i))
        err = cli.GenGif(vpath, gif, dur/2) // TODO case: duration too short
        if err != nil {
            continue
        }
        i ++
    }

    // combine gifs
    return cli.CombineGifs(tmpDir, outFile)
}

func (dp *DirProcessor) genGifFromImages(tmpDir, outFile string) error {
    // resize images to tmp dir
    imgs := dp.exp.ImgFiles
    i := 0
    for i<cfg.Ffmpeg.MaxCuts*2 && i<len(imgs) {
        ipath := imgs[i]
        opath := path.Join(tmpDir, fmt.Sprintf("range%v.png", i))
        err := cli.ResizeImgTo(ipath, opath)
        if err != nil {
            continue
        }
        i ++
    }

    // combine images to gif
    return cli.ImagesToGif(tmpDir, outFile)
}

