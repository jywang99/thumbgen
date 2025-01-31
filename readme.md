# Thumbnail Generator

## Purpose
Generate thumbnails for video files and image folders in a directory recursively.
Use case: websites serving video and image contents

## Dependencies
1. [FFmpeg](https://ffmpeg.org/) - version n7.0.1
2. [ImageMagick](https://imagemagick.org/index.php) - Version: ImageMagick 7.1.1-34 Q16-HDRI x86_64 22301

## Features
1. Generate thumbnails (1 image, 1 gif) for each video file in the directory
2. Generate the same thumbnails for each image folder in the directory

## Usage
1. Write the configuration file. Example: [config.json](conf/config.yml)
2. Build
```bash
go build -o bin/thumbgen ./src
```
3. Run \
Replace `conf/config.yml` with the path to your configuration file.
```bash
./bin/thumbgen -f conf/config.yml
```

## Example
Let's say we have the following directory structure for our `directories.input`:
```
/data/
    folder1/
        video1.mp4
        video2.mp4
        images1/
            image1.jpg
            image2.jpg
    folder2/
        video3.mp4
        videos1/
            video4.mp4
    folder3/
        mixture/
            video5.mp4
            more/
                image3.jpg
                image4.jpg
                evenmore/
                    video6.mp4
                    video7.mp4
                    image5.jpg
    $RECYCLE.BIN/
        ...
    System Volume Information/
        ...
```
And we've defined the following configuration:
```yaml
ffmpeg:
# ...
directories:
  input: /mnt/f
  output: /soft/video-prep/work/out
  temp: /soft/video-prep/work/tmp
  ignore: $RECYCLE.BIN:System Volume Information
  maxDepth: 1
files:
# ...
```
Then the script will generate the following thumbnails:
```
/soft/video-prep/work/out/ - level 0
    folder1/ - level 1
        video1.png
        video1.gif
        video2.png
        video2.gif
        images1.png - randomly selected from image1.jpg and image2.jpg, originally in images1/
        images1.gif - gif cycling through both image1.jpg and image2.jpg
    folder2/
        video3.jpg
        video3.gif
        videos1.png
        videos1.gif
    folder3/
        mixture.png - randomly selected from image3, image4, image5, and cutscenes from video5, video6, video7
        mixture.gif - can contain randomly selected cutscenes from video5, video6, video7, and/or image3, image4, image5
```
Each directory/file inside the level 1 directory (I call them leaf nodes) produce one gif + one png.\
**NOTE**: The `ffmpeg.maxCuts` specifies how many cutscenes a gif can contain. An image occupies half of the `ffmpeg.cutDuration`. Thus, specifying `maxCuts: 5` will generate a gif with 10 images if all the selected cutscenes are images.

