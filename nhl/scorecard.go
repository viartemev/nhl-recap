package nhl

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	logos         *logos.Logos
	abbreviations map[string]string
	settings      GeneratorSettings
}

type GeneratorSettings struct {
	Width      int
	Height     int
	Background color.Color
	TextColor  color.Color
	FontSize   float64
}

func NewScoreCardGenerator(l *logos.Logos, settings GeneratorSettings) ScoreCardGenerator {
	return ScoreCardGenerator{logos: l, settings: settings, abbreviations: map[string]string{
		"New Jersey Devils":     "NJD",
		"New York Islanders":    "NYI",
		"New York Rangers":      "NYR",
		"Los Angeles Kings":     "LAK",
		"San Jose Sharks":       "SJS",
		"St. Louis Blues":       "STL",
		"Tampa Bay Lightning":   "TBL",
		"Washington Capitals":   "WSH",
		"Carolina Hurricanes":   "CAR",
		"Florida Panthers":      "FLA",
		"Chicago Blackhawks":    "CHI",
		"Colorado Avalanche":    "COL",
		"Minnesota Wild":        "MIN",
		"Nashville Predators":   "NSH",
		"Winnipeg Jets":         "WPG",
		"Anaheim Ducks":         "ANA",
		"Columbus Blue Jackets": "CBJ",
		"Dallas Starts":         "DAL",
		"Edmonton Oilers":       "EDM",
		"Vancouver Canucks":     "VAN",
		"Arizona Coyotes":       "ARI",
		"Calgary Flames":        "CGY",
		"Montréal Canadiens":    "MTL",
		"Ottawa Senators":       "OTT",
		"Philadelphia Flyers":   "PHI",
		"Pittsburgh Penguins":   "PIT",
		"Toronto Maple Leafs":   "TOR",
		"Boston Bruins":         "BOS",
		"Buffalo Sabres":        "BUF",
		"Detroit Red Wings":     "DET",
		"Seattle Kraken":        "SEA",
		"Vegas Golden Knights":  "VGK",
	}}
}

func (g *ScoreCardGenerator) GenerateScoreCard(game domain.ScheduleGame) []byte {
	background := image.NewRGBA(image.Rect(0, 0, g.settings.Width, g.settings.Height))
	draw.Draw(background, background.Bounds(), &image.Uniform{C: g.settings.Background}, image.Point{}, draw.Src)
	err := drawText(background, game, g.logos, g.abbreviations)
	if err != nil {
		log.WithError(err).Error("Can't draw text")
	}
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, background)
	return buf.Bytes()
}

func drawText(canvas *image.RGBA, game domain.ScheduleGame, l *logos.Logos, abbreviations map[string]string) error {
	var (
		fgColor  image.Image
		fontFace *truetype.Font
		err      error
		fontSize = 25.0
	)
	fgColor = image.Black
	fontFace, err = freetype.ParseFont(goregular.TTF)
	if err != nil {
		return err
	}
	fontDrawer := &font.Drawer{
		Dst: canvas,
		Src: fgColor,
		Face: truetype.NewFace(fontFace, &truetype.Options{
			Size:    fontSize,
			Hinting: font.HintingFull,
		}),
	}
	awayAbbreviation := abbreviations[game.Teams.Away.Team.Name]
	awayLogo := l.GetLogoByTeam(awayAbbreviation)
	if awayLogo == nil || awayAbbreviation == "" {
		return fmt.Errorf("can't find away team: %v", game.Teams.Away.Team.Name)
	}
	homeAbbreviation := abbreviations[game.Teams.Home.Team.Name]
	homeLogo := l.GetLogoByTeam(homeAbbreviation)
	if homeLogo == nil || homeAbbreviation == "" {
		return fmt.Errorf("can't find home team: %v", game.Teams.Home.Team.Name)
	}

	drawHomeTeam(canvas, homeLogo, homeAbbreviation, game.Teams.Home.Score, fontDrawer)
	drawAwayTeam(canvas, awayLogo, awayAbbreviation, game.Teams.Away.Score, fontDrawer)

	return nil
}

func drawHomeTeam(background *image.RGBA, logo image.Image, abbreviation string, score int, fontDrawer *font.Drawer) {
	// Draw team logo
	draw.Draw(background, logo.Bounds().Add(image.Point{X: 20, Y: 55}), logo, image.Point{}, draw.Over)

	// Draw team name
	awayTeamTextBound, _ := fontDrawer.BoundString(abbreviation)
	awayTeamXPosition := fixed.I(90)
	awayTeamTextHeight := awayTeamTextBound.Max.Y - awayTeamTextBound.Min.Y
	awayTeamYPosition := fixed.I(background.Rect.Max.Y) - fixed.I(awayTeamTextHeight.Ceil())
	fontDrawer.Dot = fixed.Point26_6{
		X: awayTeamXPosition,
		Y: awayTeamYPosition,
	}
	fontDrawer.DrawString(abbreviation)

	// Draw team score
	awayTeamXPosition = fixed.I(background.Rect.Max.X - 50)
	fontDrawer.Dot = fixed.Point26_6{
		X: awayTeamXPosition,
		Y: awayTeamYPosition,
	}
	fontDrawer.DrawString(strconv.Itoa(score))
}

func drawAwayTeam(background *image.RGBA, logo image.Image, abbreviation string, score int, fontDrawer *font.Drawer) {
	// Draw team logo
	draw.Draw(background, logo.Bounds().Add(image.Point{X: 20, Y: 15}), logo, image.Point{}, draw.Over)

	// Draw team name
	homeTeamTextBound, _ := fontDrawer.BoundString(abbreviation)
	homeTeamXPosition := fixed.I(90)
	homeTeamTextHeight := homeTeamTextBound.Max.Y - homeTeamTextBound.Min.Y
	homeTeamYPosition := fixed.I((background.Rect.Max.Y)-homeTeamTextHeight.Ceil())/4 + fixed.I(homeTeamTextHeight.Ceil())
	fontDrawer.Dot = fixed.Point26_6{
		X: homeTeamXPosition,
		Y: homeTeamYPosition,
	}
	fontDrawer.DrawString(abbreviation)

	// Draw team score
	homeTeamXPosition = fixed.I(background.Rect.Max.X - 50)
	fontDrawer.Dot = fixed.Point26_6{
		X: homeTeamXPosition,
		Y: homeTeamYPosition,
	}
	fontDrawer.DrawString(strconv.Itoa(score))
}
