package nhl

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

var img []byte

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


func generateImage(w, h int, pixelColor color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			img.Set(x, y, pixelColor)
		}
	}
	return img
}

func randomColor() color.RGBA {
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	return color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}
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

func GeneratePreview()  {
	var width = 200
	var height = 200
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
	img = buf.Bytes()
}

func ShowImage() []byte {
	return img
}