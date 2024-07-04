package cli

import (
	"fmt"
	"os/exec"
	"path"
)

// example: magick -delay 5 -loop 0 range1.gif range2.gif combined.gif
func CombineGifs(srcDir, output string) error {
    logger.INFO.Println("Combining gifs to", output)
    cmd := exec.Command(
        "magick", 
        "-delay", "5",
        "-loop", "0",
        path.Join(srcDir, "*"),
        output,
    )
    _, err := execCmd(cmd)
    if err != nil {
        logger.ERROR.Println("CombineGifs command error")
        return err
    }

    return nil
}

// example: magick aaa.jpg -resize 360x240 -gravity center -background black -extent 360x240 output.png
func ResizeImgTo(input string, output string) error {
    logger.INFO.Printf("Resizing %v to %v\n", input, output)
    cmd := exec.Command(
        "magick", input, 
        "-resize", fmt.Sprintf("%dx%d", cfg.ScaleWidth, cfg.ScaleHeight),
        "-gravity", "center",
        "-background", "black",
        "-extent", fmt.Sprintf("%dx%d", cfg.ScaleWidth, cfg.ScaleHeight),
        output,
    )

    _, err := execCmd(cmd)
    if err != nil {
        logger.ERROR.Println("ResizeImgTo command error")
        return err
    }

    return nil
}

// example: magick -delay 100 -loop 0 DSC_5696.jpg  DSC_5762.jpg  DSC_5863.jpg  DSC_5918.jpg  DSC_6006.jpg  DSC_6052.jpg  DSC_6141.jpg  DSC_6261.jpg  DSC_6335.jpg  DSC_6387.jpg out.gif
func ImagesToGif(srcDir string, output string) error {
    logger.INFO.Println("Converting images to ", output)
    cmd := exec.Command(
        "magick", 
        "-delay", "100", // TODO config
        "-loop", "0",
        path.Join(srcDir, "*"),
        output,
    )
    _, err := execCmd(cmd)
    if err != nil {
        logger.ERROR.Println("ImagesToGif command error")
        return err
    }

    return err
}

