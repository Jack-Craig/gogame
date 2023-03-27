package gameplay

import (
	"github.com/Jack-Craig/gogame/src/common"
	"github.com/Jack-Craig/gogame/src/graphics"
	"github.com/Jack-Craig/gogame/src/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type Handler struct {
	gdl     *graphics.GraphicsDataLoader
	im      *input.InputManager
	players []*Player
}

type GameState interface {
	GetNextState() GameState
	Update()
	Draw(screen *ebiten.Image)
}

// PLAYSTATE
type PlayState struct {
	GameState
	Handler
	world *World
}

func NewPlayState(handler Handler) *PlayState {
	return &PlayState{
		Handler: handler,
		world:   NewWorld(handler),
	}
}

func (ps *PlayState) GetNextState() GameState {
	if ps.world.allPlayersDoneOrDead {
		return NewShopState(ps.Handler)
	}
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
	Handler
	numPlayers, windowWidth, windowHeight int
	playerData                            []*PlayerData
	playerTileIds                         []struct {
		id   uint32
		name string
	}
	readyForNextState bool
}

func NewMenuState() *MenuState {
	ms := &MenuState{}
	var pd common.PlayerDataJson
	common.LoadJSON("res/models.json", &pd)
	for name, d := range pd.Players {
		ms.playerTileIds = append(ms.playerTileIds, struct {
			id   uint32
			name string
		}{
			id:   uint32(d.ImageId),
			name: name,
		})
	}
	// Initialize handler objects
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
		for _, data := range ms.playerData {
			ms.players = append(ms.players, NewPlayer(data.id+1, data.name, nil, data.im, data.pi))
		}
		return NewPlayState(ms.Handler)
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

// SHOPSTATE
type ShopState struct {
	GameState
	Handler
}

func NewShopState(handler Handler) *ShopState {
	return &ShopState{Handler: handler}
}

func (ss *ShopState) GetNextState() GameState {
	return NewPlayState(ss.Handler)
}

func (ss *ShopState) Update() {

}

func (ss *ShopState) Draw(screen *ebiten.Image) {

}
