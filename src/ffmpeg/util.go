package ffmpeg

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"

	"jy.org/thumbgen/src/logging"
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

func execPipeCmd(c1, c2 *exec.Cmd) error {
    r, w := io.Pipe() 
    c1.Stdout = w
    c2.Stdin = r

    var b2 bytes.Buffer
    c2.Stdout = &b2

    if err := c1.Start(); err != nil {
        logger.ERROR.Printf("c1.Start() failed: %s\n", err)
        return err
    }
    if err := c2.Start(); err != nil {
        logger.ERROR.Printf("c2.Start() failed: %s\n", err)
        return err
    }

    if err := c1.Wait(); err != nil {
        logger.ERROR.Printf("c1.Wait() failed: %s\n", err)
        return err
    }

    if err := w.Close(); err != nil {
        logger.ERROR.Printf("w.Close() failed: %s\n", err)
        return err
    }

    if err := c2.Wait(); err != nil {
        logger.ERROR.Printf("c2.Wait() failed: %s\n", err)
        return err
    }

    if _, err := io.Copy(os.Stdout, &b2); err != nil {
        logger.ERROR.Printf("io.Copy failed: %s\n", err)
        return err
    }

    return nil
}
