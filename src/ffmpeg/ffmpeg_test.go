package ffmpeg_test

import (
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
    err := ff.GenPreviewGif(path.Join(testDir, "res/file_example_MP4_1280_10MG.mp4"), path.Join(os.Getenv("OUTPUT_DIR"), "test.gif"))
    assert.Nil(t, err)
}
