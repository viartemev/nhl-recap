package logos

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
)

type Logos struct {
	m map[string]image.Image
}

func (l *Logos) GetLogoByTeam(team string) image.Image {
	return l.m[team]
}

func (l *Logos) GetDefault() image.Image {
	return image.NewRGBA(image.Rectangle{image.Point{X: 0, Y: 0}, image.Point{X: 30, Y: 30}})
}

func LoadLogos() *Logos {
	logos := &Logos{m: map[string]image.Image{}}
	baseDir := "nhl/logos"
	files, err := filepath.Glob(filepath.Join(baseDir, "*.png"))
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Println(err)
		}
		img, _, _ := image.Decode(bytes.NewReader(data))
		fileName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		logos.m[fileName] = img
	}
	return logos
}
