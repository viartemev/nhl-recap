package nhl

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
)
type Scene struct {
	Width, Height int
	Img           *image.RGBA
}

func NewScene(width int, height int) *Scene {
	return &Scene{
		Width:  width,
		Height: height,
		Img:    image.NewRGBA(image.Rect(0, 0, width, height)),
		}
}

func (s *Scene) EachPixel(colorFunction func(int, int) color.RGBA) {
	for x := 0; x < s.Width; x++ {
		for y := 0; y < s.Height; y++ {
			s.Img.Set(x, y, colorFunction(x, y))
		}
	}
}

func (s *Scene) Save(filename string) {
	f, err := os.Create(filename)

	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, s.Img)
}

func GenerateScoreCard() []byte {
	//TODO use game info to generate image

	var width = 300
	var height = 100
	scene := NewScene(width, height)
	scene.EachPixel(func(x, y int) color.RGBA {
		return color.RGBA{
			uint8(x * 255 / width),
			uint8(y * 255 / height),
			100,
			255,
			}
	})
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, scene.Img)
	return buf.Bytes()
}