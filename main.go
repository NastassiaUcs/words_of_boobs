package main

import (
	"image"
	"github.com/golang/freetype"
	"io/ioutil"
	"log"
	"image/draw"
	"github.com/fogleman/gg"
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


	height := 500
	width := height * 3

	size := float64(height)


	// Initialize the context.
	fg, bg := image.Black, image.White

	img := image.NewRGBA(image.Rect(0, 0, width, height))
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

	ctx := gg.NewContextForImage(img)
	//if err = drawImage(ctx, "./boobs.png", 5, 5); err != nil {
	//	log.Panicln(err)
	//	return
	//}

	rectW := 30
	rectH := 30

	var blackCount int
	for i := 0; i < width; i += rectW + 1 {
		for j := 0; j < height; j += rectH + 1 {
			blackCount = 0
			for ri := i; ri < i + rectW; ri++ {
				for rj := j; rj < j + rectH; rj++ {
					rgba := img.RGBAAt(ri, rj)
					if rgba.R == 0 && rgba.G == 0 && rgba.B == 0 {
						blackCount++
					}
				}
			}
			if float64(blackCount) > float64(rectW * rectH) * 0.5 {
				//ctx.SetRGB(0, 128, 0)
				//ctx.DrawRectangle(float64(i), float64(j), float64(rectW), float64(rectH))
				//ctx.Fill()
				drawImage(ctx, "boobs.png", i, j, rectW, rectH)
			}
		}
	}


	ctx.SavePNG("test.png")
}