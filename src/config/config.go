package config

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type directories struct {
    Input string
    Output string
    Temp string
    Ignore map[string]bool
}

type config struct {
    Dirs directories
    TargetExts map[string]bool
}

var basePath = "/soft/video-prep/config/"

func readYmlConfig(cfg *config) {
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

func initConfig() config {
    err := godotenv.Load(path.Join(basePath, ".env"))
    if err != nil {
        log.Fatal(err)
    }

    var cfg config
    readYmlConfig(&cfg)

    cfg.Dirs = directories{
        Input: os.Getenv("INPUT_DIR"),
        Output: os.Getenv("OUTPUT_DIR"),
        Temp: os.Getenv("TEMP_DIR"),
        Ignore: make(map[string]bool),
    }

    cfg.Dirs.Ignore = stringToMap(os.Getenv("IGNORE_DIRS"))
    cfg.TargetExts = stringToMap(os.Getenv("TARGET_EXT"))

    return cfg
}
var Config = initConfig()

func stringToMap(s string) map[string]bool {
    m := make(map[string]bool)
    if s == "" {
        return m
    }

    for _, v := range strings.Split(s, ":") {
        m[v] = true
    }
    return m
}

