package gameplay

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Jack-Craig/gogame/src/graphics"
	"github.com/Jack-Craig/gogame/src/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type GameState interface {
	GetNextState() GameState
	Update()
	Draw(screen *ebiten.Image)
}

// PLAYSTATE
type PlayState struct {
	GameState
	gdl   *graphics.GraphicsDataLoader
	im    *input.InputManager
	world *World
}

func NewPlayState(gdl *graphics.GraphicsDataLoader, im *input.InputManager, spriteIds []graphics.SpriteID) *PlayState {
	return &PlayState{gdl: gdl, im: im, world: NewWorld(gdl, im, spriteIds)}
}

func (ps *PlayState) GetNextState() GameState {
	return nil
}

func (ps *PlayState) Update() {
	ps.world.Update()

}

func (ps *PlayState) Draw(screen *ebiten.Image) {
	ps.world.Draw(screen)
}

// MENUSTATE
type MenuState struct {
	GameState
	gdl                                   *graphics.GraphicsDataLoader
	im                                    *input.InputManager
	numPlayers, windowWidth, windowHeight int
	playerData                            []*PlayerData
	playerTileIds                         []struct {
		id   uint32
		name string
	}
	readyForNextState bool
}

type playerDataJson struct {
	Players map[string]playerDataJson_ `json: "players"`
}
type playerDataJson_ struct {
	ImageId int `json: "imageId"`
}

func NewMenuState() *MenuState {
	ms := &MenuState{}
	playerJsonFile, err := os.Open("res/models.json")
	if err != nil {
		log.Fatal(err)
	}
	defer playerJsonFile.Close()
	playerJsonBytes, _ := ioutil.ReadAll(playerJsonFile)
	if err != nil {
		log.Fatal(err)
	}
	var pd playerDataJson
	json.Unmarshal(playerJsonBytes, &pd)
	log.Println(pd)
	for name, d := range pd.Players {
		ms.playerTileIds = append(ms.playerTileIds, struct {
			id   uint32
			name string
		}{
			id:   uint32(d.ImageId),
			name: name,
		})
	}
	ms.gdl = graphics.NewGraphicsDataLoader()
	ms.im = input.NewInputManager()
	ms.im.InitiateJoyConConnections()
	ms.numPlayers = len(*ms.im.GetPlayerInputs())
	for id, _ := range *ms.im.GetPlayerInputs() {
		ms.playerData = append(ms.playerData, NewPlayerData(id, ms))
	}
	return ms
}

func (ms *MenuState) GetNextState() GameState {
	if ms.readyForNextState {
		var spriteIds []graphics.SpriteID
		for _, data := range ms.playerTileIds {
			spriteIds = append(spriteIds, graphics.SpriteID(data.id))
		}
		return NewPlayState(ms.gdl, ms.im, spriteIds)
	}
	return nil
}

func (ms *MenuState) Update() {
	isEveryoneReady := true
	for _, pd := range ms.playerData {
		pd.Update()
		if !pd.readyForStart {
			isEveryoneReady = false
		}
	}
	if isEveryoneReady {
		ms.readyForNextState = true
	}
}

func (ms *MenuState) Draw(screen *ebiten.Image) {
	ms.windowWidth, ms.windowHeight = screen.Size()
	for _, pd := range ms.playerData {
		pd.Draw(screen)
	}
}

type PlayerData struct {
	id            uint32
	ms            *MenuState
	im            *ebiten.Image
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

	name := pd.ms.playerTileIds[pd.curIdx].name
	boundRect := text.BoundString(font, name)
	boundRectW, boundRectH := boundRect.Size().X, boundRect.Size().Y
	text.Draw(screen, name, font, int(pd.id)*w+int(float64(w)/2)-boundRectW/2, int(float64(h/2)+.5*(guyWidth))+boundRectH, color.White)

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

// SHOPSTATE
type ShopState struct {
	GameState
}

func (ss *ShopState) GetNextState() GameState {
	return nil
}

func (ss *ShopState) Update() {

}

func (ss *ShopState) Draw(screen *ebiten.Image) {

}
