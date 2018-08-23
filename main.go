package main

import (
	"image"
	"github.com/golang/freetype"
	"io/ioutil"
	"log"
	"image/draw"
	"github.com/fogleman/gg"
	"math/rand"
	"path/filepath"
	"math"
	"flag"
)

const (
	HEIGHT      = 2000
	WIDTH       = HEIGHT * 5
	RECT_WIDTH  = 25
	//RECT_HEIGHT = 30
	IMG_FOLDER = "./img/"
	EXAMPLE_IMG = "example.jpg"
	TEXT = "aga"
)

type dotsManager struct {
	dots map[int]map[int]bool
	count int
}

type point struct {
	x, y int
}

func (self *dotsManager) addDot(x int, y int) {
	if _, ok := self.dots[x]; !ok {
		self.dots[x] = make(map[int]bool)
	}
	if _, v := self.dots[x][y]; !v {
		self.dots[x][y] = true
		self.count++
	}
}

func (self *dotsManager) removeDot(x int, y int) {
	if _, ok := self.dots[x]; ok {
		if self.dots[x][y] {
			self.dots[x][y] = false
			self.count--
		}
	}
}

func (self *dotsManager) checkDot(x int, y int) bool {
	if _, ok := self.dots[x]; ok {
		return self.dots[x][y]
	}
	return false
}

func (self *dotsManager) getList(shuffle bool) []point {
	points := make([]point, self.count)
	count := 0
	for x, column := range self.dots {
		for y, v := range column {
			if v {
				points[count] = point{x: x, y: y}
				count++
			}
		}
	}
	if shuffle {
		rand.Shuffle(count, func(i, j int) {
			points[i], points[j] = points[j], points[i]
		})
	}
	return points
}

func (self *dotsManager) getRandomDot() point {
	list := self.getList(false)
	return list[rand.Intn(self.count)]
}

func createDots() *dotsManager {
	var d = dotsManager{}
	d.dots = make(map[int]map[int]bool)
	d.count = 0
	return &d
}

func prepareImage(filename string, width int) (img image.Image, err error) {
	source, err := gg.LoadPNG(filename)
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

func main() {
	var width int
	text := flag.String("text", "geeks", "a string")
	flag.IntVar(&width, "width", WIDTH, "an int")
	flag.Parse()
	log.Println(*text)

	fontBytes, err := ioutil.ReadFile("./fonts/NotoSans-Bold.ttf")
	if err != nil {
		log.Panicln(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Panicln(err)
		return
	}


	size := float64(HEIGHT * 0.75)


	// Initialize the context.
	fg, bg := image.Black, image.White //image.NewUniform(color.Gray16{0xaaaa})

	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	draw.Draw(img, img.Bounds(), bg, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(fg)


	pt := freetype.Pt(2, int(size) - 34)


	s := *text

	_, err = c.DrawString(s, pt)
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Printf("text %s is drawing, please wait\n", s)
	}

	cleanImg := image.NewRGBA(image.Rect(0, 0, width, HEIGHT))
	draw.Draw(cleanImg, cleanImg.Bounds(), bg, image.ZP, draw.Src)
	cleanCtx := gg.NewContextForImage(cleanImg)

	var (
		allDots = createDots()
		filledDots = createDots()
	)

	var example image.Image
	if EXAMPLE_IMG != "" {
		if example, err = gg.LoadImage("example.jpg"); err != nil {
			panic(err)
		}
	}

	var (
		blackCount int
		r, g, b uint32
	)
	for i := 0; i < width; i ++ {
		for j := 0; j < HEIGHT; j ++ {
			blackCount = 0
			if EXAMPLE_IMG == "" {
				rgba := img.RGBAAt(i, j)
				r, g, b = uint32(rgba.R), uint32(rgba.G), uint32(rgba.B)
			} else {
				r, g, b, _ = example.At(i, j).RGBA()
			}
			if r == 0 && g == 0 && b == 0 {
				blackCount++
				allDots.addDot(i, j)
			}
		}
	}

	var (
		images []image.Image
		filenames []string
		points = allDots.getList(true)
	)
	if filenames, err = filepath.Glob(IMG_FOLDER + "*.png"); err != nil {
		panic(err)
	} else {
		images = make([]image.Image, len(filenames))
		for i := range filenames {
			if images[i], err = prepareImage(filenames[i], RECT_WIDTH); err != nil {
				panic(err)
			}
		}
	}
	imgCount := len(images)

	for i := 0; i < len(points); i++ {
		p := points[i]
		if filledDots.checkDot(p.x, p.y) {
			continue
		}
		drawImage(cleanCtx, images[rand.Intn(imgCount)], p.x, p.y, filledDots)
	}

	cleanCtx.SavePNG("result.png")
	log.Println("done, check file result.png")
}