package nhl

import (
    "bytes"
    "fmt"
    "github.com/golang/freetype"
    "github.com/golang/freetype/truetype"
    log "github.com/sirupsen/logrus"
    "golang.org/x/image/draw"
    "golang.org/x/image/font"
    "golang.org/x/image/font/gofont/goregular"
    "golang.org/x/image/math/fixed"
    "image"
    "image/color"
    "image/png"
)



func GenerateScoreCard(message *GameInfo) []byte {
	var width = 300
	var height = 100
	img, _ := createCard(width, height, "#ffffff", fmt.Sprintf("%d - %d", message.HomeTeam.Score, message.AwayTeam.Score))
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, img)
	return buf.Bytes()
}

func createCard(width int, height int, bg string, text string) (*image.RGBA, error) {
	bgColor, err := hexToRGBA(bg)
	if err != nil {
		log.Fatal(err)
	}
	background := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(background, background.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)
	_ = drawText(background, text)
	return background, err
}

func drawText(canvas *image.RGBA, text string) error {
	var (
		fgColor  image.Image
		fontFace *truetype.Font
		err      error
		fontSize = 32.0
	)
	fgColor = image.Black
	fontFace, err = freetype.ParseFont(goregular.TTF)
	fontDrawer := &font.Drawer{
		Dst: canvas,
		Src: fgColor,
		Face: truetype.NewFace(fontFace, &truetype.Options{
			Size:    fontSize,
			Hinting: font.HintingFull,
		}),
	}
	textBounds, _ := fontDrawer.BoundString(text)
	xPosition := (fixed.I(canvas.Rect.Max.X) - fontDrawer.MeasureString(text)) / 2
	textHeight := textBounds.Max.Y - textBounds.Min.Y
	yPosition := fixed.I((canvas.Rect.Max.Y)-textHeight.Ceil())/2 + fixed.I(textHeight.Ceil())
	fontDrawer.Dot = fixed.Point26_6{
		X: xPosition,
		Y: yPosition,
	}
	fontDrawer.DrawString(text)
	return err
}

func hexToRGBA(hex string) (color.RGBA, error) {
	var (
		rgba             color.RGBA
		err              error
		errInvalidFormat = fmt.Errorf("invalid")
	)
	rgba.A = 0xff
	if hex[0] != '#' {
		return rgba, errInvalidFormat
	}
	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}
	switch len(hex) {
	case 7:
		rgba.R = hexToByte(hex[1])<<4 + hexToByte(hex[2])
		rgba.G = hexToByte(hex[3])<<4 + hexToByte(hex[4])
		rgba.B = hexToByte(hex[5])<<4 + hexToByte(hex[6])
	case 4:
		rgba.R = hexToByte(hex[1]) * 17
		rgba.G = hexToByte(hex[2]) * 17
		rgba.B = hexToByte(hex[3]) * 17
	default:
		err = errInvalidFormat
	}
	return rgba, err
}
