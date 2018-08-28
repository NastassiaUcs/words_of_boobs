package generator

import (
	"regexp"
	"image"
	"github.com/fogleman/gg"
	"math"
	"io/ioutil"
	"log"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"path/filepath"
	"image/draw"
	"math/rand"
	"time"
	"image/color"
)

const (
	FONTS_FOLDER   = "./fonts/"
	IMAGES_FOLDER  = "./img/"
	RESULTS_FOLDER = "./results/"
	EXAMPLES_FOLDER = "./examples/"
	FONT_POINTS = 750
)

var (
	g generator
	textContent = gg.NewContext(50, 50)
)

func init() {
	g = generator{}
	g.imageSets = make(map[string][]image.Image)
	g.fonts = make(map[string]*truetype.Font)

	textContent.LoadFontFace(FONTS_FOLDER + "Symbola.ttf", FONT_POINTS)
}


func Reload(imageWidth int) {
	g.loadFonts()
	g.loadImagesSets(imageWidth)
}


func prepareFont(fontName string) (f *truetype.Font, err error) {
	fontBytes, err := ioutil.ReadFile(fontName)
	if err != nil {
		log.Panicln(err)
		return
	}
	f, err = freetype.ParseFont(fontBytes)
	return
}


func isPng(filename string) bool {
	re := regexp.MustCompile("\\.png$")
	return re.MatchString(filename)
}

func prepareImage(filename string, width int) (source image.Image, err error) {
	if isPng(filename) {
		source, err = gg.LoadPNG(filename)
	} else {
		source, err = gg.LoadImage(filename)
	}
	if err != nil {
		return nil, err
	}

	size := source.Bounds().Size()

	var scale float64

	scale = float64(width) / float64(size.X)

	height := int(math.Round(float64(size.Y) * scale))

	ctx := gg.NewContext(width, height)
	ctx.Scale(scale, scale)
	ctx.DrawImage(source, 0, 0)

	return ctx.Image(), err
}


func drawImage(ctx *gg.Context, img image.Image, x int, y int, filledDots *dotsManager) {
	size := img.Bounds().Size()
	imgWidth, imgHeight := size.X, size.Y
	x -= imgWidth / 2
	y -= imgHeight / 2
	for i := x; i < x + imgWidth; i++ {
		for j := y; j < y + imgHeight; j++ {
			filledDots.addDot(i, j)
		}
	}
	//log.Printf("x = %d, y = %d, width = %d, height = %d\n", x, y, imgWidth, imgHeight)
	ctx.DrawImage(img, x, y)
}

func getFilename() string {
	t := time.Now()
	return t.Format("20060102150405") + ".png"
}


type generator struct {
	imageSets map[string][]image.Image
	fonts map[string]*truetype.Font
}



func (this *generator) loadFonts() {
	var fonts = make(map[string]*truetype.Font)
	filenames, err := filepath.Glob(FONTS_FOLDER + "*.ttf")
	if err != nil {
		log.Panic(err)
	}
	for _, filename := range filenames {
		if fonts[filepath.Base(filename)], err = prepareFont(filename); err != nil {
			log.Panic(err)
		}
	}
	this.fonts = fonts
}

func (this *generator) loadImagesSets(imageWidth int) {
	this.imageSets = make(map[string][]image.Image)
	dirs, err := filepath.Glob(IMAGES_FOLDER + "*")
	if err != nil {
		log.Panic(err)
	}
	for _, dirName := range dirs {
		filenames, err := filepath.Glob(dirName + "/*")
		if err != nil {
			log.Panic(err)
		}
		re := regexp.MustCompile("\\.(png|jpg|jpeg)$")
		images := make([]image.Image, len(filenames))
		for i := range filenames {
			if !re.MatchString(filenames[i]) {
				continue
			}
			if images[i], err = prepareImage(filenames[i], imageWidth); err != nil {
				log.Panic(err)
			}
		}
		this.imageSets[filepath.Base(dirName)] = images
	}
}




func (this *generator) process(source image.Image, imgSet string) (filename string) {
	var (
		bg = image.White
		img draw.Image
		sourceCtx = gg.NewContextForImage(source)
		ctx    *gg.Context
		allDots = createDots()
		filledDots = createDots()
		r, g, b, a uint32
	)

	img = image.NewRGBA(source.Bounds())
	draw.Draw(img, img.Bounds(), bg, image.ZP, draw.Src)
	ctx = gg.NewContextForImage(img)

	for i := 0; i < sourceCtx.Width(); i ++ {
		for j := 0; j < sourceCtx.Height(); j ++ {
			r, g, b, a = source.At(i, j).RGBA()
			if r == 0 && g == 0 && b == 0 && a != 0 {
				allDots.addDot(i, j)
			}
		}
	}


	var (
		images []image.Image
		points = allDots.getList(true)
	)
	images = this.imageSets[imgSet]

	imgCount := len(images)
	drawnCount := 0
	for i := 0; i < len(points); i++ {
		p := points[i]
		if filledDots.checkDot(p.x, p.y) {
			continue
		}
		drawImage(ctx, images[rand.Intn(imgCount)], p.x, p.y, filledDots)
		drawnCount++
	}

	filename = getFilename()
	ctx.SavePNG(RESULTS_FOLDER + filename)
	log.Printf("%d images was drawn\n", drawnCount)
	return
}

func GenerateImageForText(text, fontName, imgSet string, height, width int) (filename string, err error) {
	tw, th := textContent.MeasureString(text)

	var (
		padding = 50
	)

	textHeight := int(th)
	ctx := gg.NewContext(int(tw) + padding*2, textHeight + textHeight / 40 * 16 + padding)
	ctx.SetColor(color.White)
	ctx.Clear()
	ctx.SetColor(color.Black)
	if err = ctx.LoadFontFace(FONTS_FOLDER + "Symbola.ttf", FONT_POINTS); err != nil {
		return
	}

	ctx.DrawString(text, float64(padding), th + float64(padding))

	filename = g.process(ctx.Image(), imgSet)

	return
}


func GenerateImageForImage(imageName, imgSet string) (filename string, err error) {
	var img image.Image

	if img, err = gg.LoadImage(EXAMPLES_FOLDER + imageName); err != nil {
		return
	}

	g.process(img, imgSet)

	return
}