package main

import (
	"flag"

	"jy.org/thumbgen/src/config"
)

type args struct {
    config config.ConfigOverride
}

func parseArgs() args {
    var a args
    flag.StringVar(&a.config.YmlPath, "f", "", "Path to the configuration file")
    flag.StringVar(&a.config.InputDir, "i", "", "Path to the input directory")
    flag.StringVar(&a.config.OutputDir, "o", "", "Path to the output directory")
    flag.Parse()
    return a
}

