package ffmpeg

import (
	"log"
	"os"
	"os/exec"
	"path"

	"jy.org/videop/src/config"
	"jy.org/videop/src/logging"
)

var logger = logging.Logger

// list of cut start points
// evenly distribute cuts
func getCuts(duration, cutlen float64, maxCuts int) []float64 {
    cuts := make([]float64, 0)
    for i := 0; i < maxCuts; i++ {
        start := float64(i) * duration / float64(maxCuts)
        if len(cuts) > 0 && cuts[len(cuts) - 1] + cutlen > start {
            continue
        }
        cuts = append(cuts, start)
    }
    return cuts
}

func mkTmpDir(fileNm string) (string, error) {
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

func execCmd(cmd *exec.Cmd) ([]byte, error) {
    out, err := cmd.CombinedOutput()
    if err != nil {
        log.Printf("Error when executing command: %v\n", cmd.String())
        log.Printf("stderr: %v\n", err)
        log.Printf("stdout: %v\n", string(out))
        return nil, err
    }
    return out, nil
}
