package gamestate

import (
	"log"

	"github.com/Jack-Craig/gogame/src/gameobject"
	"github.com/Jack-Craig/gogame/src/graphics"
	"github.com/Jack-Craig/gogame/src/input"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameState interface {
	GetNextState() GameState
	Update()
	Draw(screen *ebiten.Image)
}

// PLAYSTATE
type PlayState struct {
	GameState
	inited        bool
	gdl           *graphics.GraphicsDataLoader
	im            *input.InputManager
	gameObjects   []*gameobject.GameObject
	playerObjects []*gameobject.Player
	// World (Blocks)
	// Entities
	//	- Players
	//	- Enemies
	// Inputs
}

func (ps *PlayState) GetNextState() GameState {
	return nil
}

func (ps *PlayState) Update() {
	if !ps.inited {
		ps.init()
	}
	for _, player := range ps.playerObjects {
		player.Update()
	}
}

func (ps *PlayState) init() {
	log.Println("Playstate Init")
	ps.gdl = graphics.NewGraphicsDataLoader("res/play")
	// Load player objects
	for id, playerInput := range *ps.im.GetPlayerInputs() {
		log.Printf("Created player: ID: %d\n", id)
		p := gameobject.NewPlayer(1, 50, 1, ps.GetGDL().GetSpriteImage(1), playerInput)
		ps.gameObjects = append(ps.gameObjects, &p.GameObject)
		ps.playerObjects = append(ps.playerObjects, p)
	}
	ps.inited = true
}

func (ps *PlayState) GetGDL() *graphics.GraphicsDataLoader {
	return ps.gdl
}

func (ps *PlayState) Draw(screen *ebiten.Image) {
	for _, gobj := range ps.gameObjects {
		gobj.Draw(screen)
	}
}

// MENUSTATE
type MenuState struct {
	GameState
	inited bool
	im     *input.InputManager
}

func (ms *MenuState) GetNextState() GameState {
	if ms.inited {
		return &PlayState{
			im: ms.im,
		}
	}
	return nil
}

func (ms *MenuState) Update() {
	if !ms.inited {
		ms.init()
	}
}

func (ms *MenuState) init() {
	ms.im = input.NewInputManager()
	ms.im.InitiateJoyConConnections()
	ms.inited = true
}

func (ms *MenuState) Draw(screen *ebiten.Image) {

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
