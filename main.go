package main

import (
	"./generator"
	"flag"
	"log"
)

const (
	HEIGHT      = 2000
	WIDTH       = HEIGHT * 5
	RECT_WIDTH  = 70
	//RECT_HEIGHT = 30
	IMG_FOLDER = "boobs"
	TEXT       = "geeks"
	FONT_NAME  = "NotoSans-Bold.ttf"
)


func main() {
	var (
		width        int
		imageWidth   int
		text         string
		exampleImage string
		fontName     string
		imagesFolder string
	)
	flag.StringVar(&text, "text", TEXT, "a string")
	flag.IntVar(&width, "width", WIDTH, "an int")
	flag.IntVar(&imageWidth, "image_width", RECT_WIDTH, "an int")
	flag.StringVar(&exampleImage, "example", "", "image path/filename or empty for text")
	flag.StringVar(&fontName, "font", FONT_NAME, "filename in folder fonts")
	flag.StringVar(&imagesFolder, "images_folder", IMG_FOLDER, "path to folder with images")
	flag.Parse()

	generator.Reload(imageWidth)

	var err error
	if exampleImage != "" {
		_, err = generator.GenerateImageForImage(exampleImage, IMG_FOLDER)
	} else {
		_, err = generator.GenerateImageForText(text, fontName, imagesFolder, HEIGHT, width)
	}
	if err != nil {
		log.Panic(err)
	}
}