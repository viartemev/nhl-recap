package nhl

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

func GenerateScoreCard(message *GameInfo) []byte {
	var width = 300
	var height = 100
	img, _ := createCard(width, height, color.White, message)
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, img)
	return buf.Bytes()
}

func createCard(width int, height int, color color.Color, message *GameInfo) (*image.RGBA, error) {
	background := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(background, background.Bounds(), &image.Uniform{C: color}, image.Point{}, draw.Src)
	_ = drawText(background, message)
	return background, nil
}

func drawText(canvas *image.RGBA, message *GameInfo) error {
	var (
		fgColor  image.Image
		fontFace *truetype.Font
		err      error
		fontSize = 25.0
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

	img, err := os.Open("logo.png")
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}

	logo, err := png.Decode(img)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer img.Close()

	drawHomeTeam(canvas, logo, message, fontDrawer)
	drawAwayTeam(canvas, logo, message, fontDrawer)

	return err
}

func drawAwayTeam(background *image.RGBA, logo image.Image, message *GameInfo, fontDrawer *font.Drawer) {
	// Draw team logo
	draw.Draw(background, logo.Bounds().Add(image.Point{X: 20, Y: 15}), logo, image.Point{}, draw.Over)
	
	// Draw team name
	awayTeamTextBound, _ := fontDrawer.BoundString(message.AwayTeam.Name)
	awayTeamXPosition := fixed.I(90)
	awayTeamTextHeight := awayTeamTextBound.Max.Y - awayTeamTextBound.Min.Y
	awayTeamYPosition := fixed.I(background.Rect.Max.Y) - fixed.I(awayTeamTextHeight.Ceil())
	fontDrawer.Dot = fixed.Point26_6{
		X: awayTeamXPosition,
		Y: awayTeamYPosition,
	}
	fontDrawer.DrawString(message.AwayTeam.Name)

	// Draw team score
	awayTeamXPosition = fixed.I(background.Rect.Max.X - 50)
	fontDrawer.Dot = fixed.Point26_6{
		X: awayTeamXPosition,
		Y: awayTeamYPosition,
	}
	fontDrawer.DrawString(strconv.Itoa(message.AwayTeam.Score))
}

func drawHomeTeam(background *image.RGBA, logo image.Image, message *GameInfo, fontDrawer *font.Drawer) {
	// Draw team logo
	draw.Draw(background, logo.Bounds().Add(image.Point{X: 20, Y: 55}), logo, image.Point{}, draw.Over)

	// Draw team name
	homeTeamTextBound, _ := fontDrawer.BoundString(message.HomeTeam.Name)
	homeTeamXPosition := fixed.I(90)
	homeTeamTextHeight := homeTeamTextBound.Max.Y - homeTeamTextBound.Min.Y
	homeTeamYPosition := fixed.I((background.Rect.Max.Y)-homeTeamTextHeight.Ceil())/4 + fixed.I(homeTeamTextHeight.Ceil())
	fontDrawer.Dot = fixed.Point26_6{
		X: homeTeamXPosition,
		Y: homeTeamYPosition,
	}
	fontDrawer.DrawString(message.HomeTeam.Name)

	// Draw team score
	homeTeamXPosition = fixed.I(background.Rect.Max.X - 50)
	fontDrawer.Dot = fixed.Point26_6{
		X: homeTeamXPosition,
		Y: homeTeamYPosition,
	}
	fontDrawer.DrawString(strconv.Itoa(message.HomeTeam.Score))
}
