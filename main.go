package main

import (
	"image"
	"github.com/golang/freetype"
	"io/ioutil"
	"log"
	"image/draw"
	"github.com/fogleman/gg"
	"math/rand"
)

const (
	HEIGHT      = 500
	WIDTH       = HEIGHT * 3
	RECT_WIDTH  = 30
	RECT_HEIGHT = 30
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
	if _, ok := self.dots[x]; !ok {
		return self.dots[x][y]
	}
	return false
}

func (self *dotsManager) getList() []point {
	points := make([]point, self.count)
	i := 0
	for x, column := range self.dots {
		for y, v := range column {
			if v {
				points[i] = point{x: x, y: y}
				i++
			}
		}
	}
	return points
}

func (self *dotsManager) getRandomDot() point {
	list := self.getList()
	return list[rand.Intn(self.count)]
}

func createDots() *dotsManager {
	var d = dotsManager{}
	d.dots = make(map[int]map[int]bool)
	d.count = 0
	return &d
}

var (
	d = createDots()
)

func drawImage(ctx *gg.Context, filename string, x int, y int, w int, h int) error {

	source, err := gg.LoadPNG(filename)
	if err != nil {
		return err
	}

	size := source.Bounds().Size()

	var scaleX, scaleY float64

	scaleX = float64(w) / float64(size.X)
	scaleY = float64(h) / float64(size.Y)

	ctx.Scale(scaleX, scaleY)
	log.Printf("x=%d y=%d scaleX=%f scaleY=%f\n", x, y, scaleX, scaleY)
	xx := int(float64(x) / scaleX)
	yy := int(float64(y) / scaleY)
	ctx.DrawImage(source, xx, yy)
	ctx.Scale(1 / scaleX, 1 / scaleY)

	return nil
}

func main() {
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


	size := float64(HEIGHT)


	// Initialize the context.
	fg, bg := image.Black, image.White

	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	draw.Draw(img, img.Bounds(), bg, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(fg)


	pt := freetype.Pt(2, int(size) - 34)

	s := "ага"

	_, err = c.DrawString(s, pt)
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Printf("text %s was drawn\n", s)
	}

	cleanImg := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	cleanCtx := gg.NewContextForImage(cleanImg)

	var blackCount int
	for i := 0; i < WIDTH; i += RECT_WIDTH + 1 {
		for j := 0; j < HEIGHT; j += RECT_HEIGHT + 1 {
			blackCount = 0
			for ri := i; ri < i + RECT_WIDTH; ri++ {
				for rj := j; rj < j + RECT_HEIGHT; rj++ {
					rgba := img.RGBAAt(ri, rj)
					if rgba.R == 0 && rgba.G == 0 && rgba.B == 0 {
						blackCount++
						d.addDot(ri, rj)
					}
				}
			}
			if float64(blackCount) > float64(RECT_WIDTH * RECT_HEIGHT) * 0.5 {
				//ctx.SetRGB(0, 128, 0)
				//ctx.DrawRectangle(float64(i), float64(j), float64(rectW), float64(rectH))
				//ctx.Fill()
				//drawImage(ctx, "boobs.png", i, j, RECT_WIDTH, RECT_HEIGHT)
			}
		}
	}

	for i := 0; i < 1000; i++ {
		p := d.getRandomDot()
		drawImage(cleanCtx, "boobs.png", p.x, p.y, RECT_WIDTH, RECT_HEIGHT)
	}


	cleanCtx.SavePNG("test.png")
}