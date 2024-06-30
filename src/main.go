package main

import (
	"fmt"

	"jy.org/videop/src/config"
)

func main() {
    defer fmt.Println("End of the program")
    fmt.Println("Video preprocessor")
    cfg := config.Config
    fmt.Println("Input dir:", cfg.Dirs.Input)
}

