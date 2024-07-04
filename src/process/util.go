package process

import (
	"jy.org/thumbgen/src/config"
	"jy.org/thumbgen/src/logging"
)

var logger = logging.Logger
var cfg = config.Config

// list of cut start points
// evenly distribute cuts
func getCuts(dur float64) []float64 {
    cuts := make([]float64, 0)
    for i := 1; i <= ffcfg.MaxCuts; i++ {
        start := float64(i) * dur / float64(ffcfg.MaxCuts+1)
        if len(cuts) > 0 && cuts[len(cuts) - 1] + ffcfg.CutDuration > start {
            continue
        }
        cuts = append(cuts, start)
    }
    return cuts
}

