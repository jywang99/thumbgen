package config

import (
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type directories struct {
    Input string
    Output string
    Temp string
}

type config struct {
    Dirs directories
}

var basePath = "/soft/video-prep/config/"

func ReadConfig(cfg *config) {
    f, err := os.Open(path.Join(basePath, "config.yml"))
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    s, _ := f.Stat()
    if s.Size() == 0 {
        return
    }

    decoder := yaml.NewDecoder(f)
    err = decoder.Decode(&cfg)
    if err != nil {
        log.Fatal(err)
    }
}

func SetupConfig() config {
    err := godotenv.Load(path.Join(basePath, ".env"))
    if err != nil {
        log.Fatal(err)
    }

    var cfg config
    ReadConfig(&cfg)

    cfg.Dirs = directories{
        Input: os.Getenv("INPUT_DIR"),
        Output: os.Getenv("OUTPUT_DIR"),
        Temp: os.Getenv("TEMP_DIR"),
    }

    return cfg
}
var Config = SetupConfig()

