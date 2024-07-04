package directory

import (
	"os"
	"path"
	"path/filepath"

	"jy.org/thumbgen/src/config"
	"jy.org/thumbgen/src/files"
	"jy.org/thumbgen/src/logging"
	"jy.org/thumbgen/src/process/ffmpeg"
)

var logger = logging.Logger
var cfg = config.Config

type DirProcessor struct{
    exp *files.Explorer
    ff *ffmpeg.Ffmpeg
    TargetDir string
}

func NewDirProcessor(dir string, ff *ffmpeg.Ffmpeg, tdir string) (*DirProcessor, error) {
    // scan source dir
    exp, err := files.NewExplorer(dir)
    if err != nil {
        return nil, err
    }
    return &DirProcessor{
        exp: exp,
        ff: ff,
        TargetDir: tdir,
    }, nil
}

func (dp *DirProcessor) GenPreviewGif() error {
    vidCnt := dp.exp.GetFileCount(files.Video)
    imgCnt := dp.exp.GetFileCount(files.Image)

    if vidCnt == 0 && imgCnt == 0 {
        // TODO default thumbnail
        logger.INFO.Printf("No video or image files found in %v\n", dp.exp.Dir)
        return nil
    }

    if vidCnt > 0 {
        return dp.genGifFromVideos()
    } 
    return dp.genGifFromImages()
}

func (dp *DirProcessor) GenPreviewImg() error {
    // TODO if at least one image, resize and use it as thumbnail
    // if not, use middle frame from a random video
    return nil
}

func (dp *DirProcessor) genGifFromVideos() error {
    // make tmp dir
    tmpDir, err := files.MkTmpDir(dp.exp.Dir)
    if err != nil {
        return err
    }
    defer os.RemoveAll(tmpDir)

    // generate gif for each video
    // process at most MaxCuts videos
    vids := dp.exp.VidFiles
    gifs := make([]string, len(vids))
    i := 0
    for i<cfg.Ffmpeg.MaxCuts && i<len(vids) {
        vid := dp.ff.NewFfVideo(vids[i], tmpDir)
        dur, err := vid.GetDuration()
        if err != nil {
            continue
        }
        gif := path.Join(tmpDir, files.GetBaseName(vid.Path, false) + ".gif")
        err = vid.GenGif(gif, dur/2) // TODO case: duration too short
        if err != nil {
            continue
        }
        gifs[i] = gif
        i ++
    }

    // combine gifs
    outFile := path.Join(dp.TargetDir, filepath.Base(dp.exp.Dir) + ".gif")
    return ffmpeg.CombineGifs(gifs, outFile)
}

func (dp *DirProcessor) genGifFromImages() error {
    // TODO images -> gif
    return nil
}
