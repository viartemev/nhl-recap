package nhl

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"nhl-recap/nhl/domain"
	"nhl-recap/nhl/logos"
	"strconv"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

type ScoreCardGenerator struct {
	logos logos.Logos
}

func NewScoreCardGenerator(l logos.Logos) ScoreCardGenerator {
	return ScoreCardGenerator{logos: l}
}

func (g *ScoreCardGenerator) GenerateScoreCard(game domain.ScheduleGame) []byte {
	var width = 300
	var height = 100
	background := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(background, background.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	_ = drawText(background, game, &g.logos)
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, background)
	return buf.Bytes()
}

func drawText(canvas *image.RGBA, game domain.ScheduleGame, l *logos.Logos) error {
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

	homeLogo := l.GetLogoByTeam(game.Teams.Home.Team.Name)
	drawAwayTeam(canvas, homeLogo, game, fontDrawer)
	awayLogo := l.GetLogoByTeam(game.Teams.Away.Team.Name)
	drawHomeTeam(canvas, awayLogo, game, fontDrawer)

	return err
}

func drawHomeTeam(background *image.RGBA, logo image.Image, game domain.ScheduleGame, fontDrawer *font.Drawer) {
	// Draw team logo
	draw.Draw(background, logo.Bounds().Add(image.Point{X: 20, Y: 15}), logo, image.Point{}, draw.Over)

	// Draw team name
	awayTeamTextBound, _ := fontDrawer.BoundString(game.Teams.Home.Team.Name)
	awayTeamXPosition := fixed.I(90)
	awayTeamTextHeight := awayTeamTextBound.Max.Y - awayTeamTextBound.Min.Y
	awayTeamYPosition := fixed.I(background.Rect.Max.Y) - fixed.I(awayTeamTextHeight.Ceil())
	fontDrawer.Dot = fixed.Point26_6{
		X: awayTeamXPosition,
		Y: awayTeamYPosition,
	}
	fontDrawer.DrawString(game.Teams.Home.Team.Name)

	// Draw team score
	awayTeamXPosition = fixed.I(background.Rect.Max.X - 50)
	fontDrawer.Dot = fixed.Point26_6{
		X: awayTeamXPosition,
		Y: awayTeamYPosition,
	}
	fontDrawer.DrawString(strconv.Itoa(game.Teams.Home.Score))
}

func drawAwayTeam(background *image.RGBA, logo image.Image, message domain.ScheduleGame, fontDrawer *font.Drawer) {
	// Draw team logo
	draw.Draw(background, logo.Bounds().Add(image.Point{X: 20, Y: 55}), logo, image.Point{}, draw.Over)

	// Draw team name
	homeTeamTextBound, _ := fontDrawer.BoundString(message.Teams.Away.Team.Name)
	homeTeamXPosition := fixed.I(90)
	homeTeamTextHeight := homeTeamTextBound.Max.Y - homeTeamTextBound.Min.Y
	homeTeamYPosition := fixed.I((background.Rect.Max.Y)-homeTeamTextHeight.Ceil())/4 + fixed.I(homeTeamTextHeight.Ceil())
	fontDrawer.Dot = fixed.Point26_6{
		X: homeTeamXPosition,
		Y: homeTeamYPosition,
	}
	fontDrawer.DrawString(message.Teams.Away.Team.Name)

	// Draw team score
	homeTeamXPosition = fixed.I(background.Rect.Max.X - 50)
	fontDrawer.Dot = fixed.Point26_6{
		X: homeTeamXPosition,
		Y: homeTeamYPosition,
	}
	fontDrawer.DrawString(strconv.Itoa(message.Teams.Away.Score))
}
