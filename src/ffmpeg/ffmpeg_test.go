package ffmpeg_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"jy.org/videop/src/ffmpeg"
)

var (
    curDir, _ = os.Getwd()
    testDir = path.Join(curDir, "../../test")
)

func TestPath(t *testing.T) {
    _, err := os.Stat(testDir)
    assert.Nil(t, err)
}

func TestGenPreviewGif(t *testing.T) {
    ff := ffmpeg.NewFfmpeg()
    path, err := ff.GenPreviewGif(path.Join(testDir, "res/file_example_MP4_1280_10MG.mp4"), path.Join(testDir, "out/test.gif"))
    fmt.Println(path)
    assert.Nil(t, err)
}
