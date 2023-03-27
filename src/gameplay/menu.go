package gameplay

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/Jack-Craig/gogame/src/graphics"
	"github.com/Jack-Craig/gogame/src/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// For menu state player data, such as which frame is selected, bg color, etc
type PlayerData struct {
	id            uint32
	ms            *MenuState
	im            *ebiten.Image
	name          string
	color         color.Color
	pi            *input.PlayerInput
	curIdx        int
	readyForStart bool
	// Stuff for changing idx
	lastChange      int64
	lastChangeStart int64
	changeDelayMs   int64
}

func NewPlayerData(id uint32, ms *MenuState) *PlayerData {
	r, g, b := uint8(rand.Uint32()%128), uint8(rand.Uint32()%128), uint8(rand.Uint32()%128)
	return &PlayerData{
		ms:            ms,
		id:            id,
		im:            ms.gdl.GetSpriteImage(graphics.SpriteID(ms.playerTileIds[0].id)),
		name:          ms.playerTileIds[0].name,
		color:         color.RGBA{r + 100, g + 100, b + 100, 255},
		pi:            (*ms.im.GetPlayerInputs())[id],
		changeDelayMs: 300,
	}
}

func (pd *PlayerData) Draw(screen *ebiten.Image) {
	bg := ebiten.NewImage(pd.ms.windowWidth/pd.ms.numPlayers, pd.ms.windowHeight)
	bg.Fill(pd.color)

	font := *pd.ms.gdl.GetFontNormal()
	dio := ebiten.DrawImageOptions{}
	w, h := bg.Size()
	if pd.readyForStart {
		dio.ColorM.Scale(.35, .8, .35, 1)
	}

	dio.GeoM.Translate(float64(int(pd.id)*w), 0)
	screen.DrawImage(bg, &dio)

	guyWidth := float64(w) * .8
	wGuy, hGuy := pd.im.Size()
	dio.GeoM.Reset()
	dio.GeoM.Scale(guyWidth/float64(wGuy), guyWidth/float64(hGuy))
	dio.GeoM.Translate(float64(int(pd.id)*w), 0)
	dio.GeoM.Translate(float64(w)/2-.5*guyWidth, float64(h/2)-.5*(guyWidth))

	screen.DrawImage(pd.im, &dio)

	boundRect := text.BoundString(font, pd.name)
	boundRectW, boundRectH := boundRect.Size().X, boundRect.Size().Y
	text.Draw(screen, pd.name, font, int(pd.id)*w+int(float64(w)/2)-boundRectW/2, int(float64(h/2)+.5*(guyWidth))+boundRectH, color.White)

	if pd.readyForStart {
		text.Draw(screen, "Ready", font, int(pd.id)*w, 20, color.White)
	}
}

func (pd *PlayerData) Update() {
	timeNow := time.Now().UnixMilli()
	cycle, _ := pd.pi.GetAxes()
	if cycle != 0 {
		if pd.changeDelayMs < timeNow-pd.lastChange {
			if cycle > 0 {
				pd.curIdx++
				if pd.curIdx >= len(pd.ms.playerTileIds) {
					pd.curIdx = 0
				}
			} else {
				pd.curIdx--
				if pd.curIdx < 0 {
					pd.curIdx = len(pd.ms.playerTileIds) - 1
				}
			}
			pd.im = pd.ms.gdl.GetSpriteImage(graphics.SpriteID(pd.ms.playerTileIds[pd.curIdx].id))
			pd.name = pd.ms.playerTileIds[pd.curIdx].name
			pd.lastChange = timeNow
		}
	}
	if pd.pi.IsButtonPressed(input.JoyConX) {
		if pd.changeDelayMs < timeNow-pd.lastChangeStart {
			pd.readyForStart = !pd.readyForStart
			pd.lastChangeStart = timeNow
		}
	}
}
