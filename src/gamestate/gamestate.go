package gamestate

import (
	"github.com/Jack-Craig/gogame/src/gameplay"
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
	inited bool
	gdl    *graphics.GraphicsDataLoader
	im     *input.InputManager
	world  *gameplay.World
}

func (ps *PlayState) GetNextState() GameState {
	return nil
}

func (ps *PlayState) Update() {
	if !ps.inited {
		ps.init()
	}
	ps.world.Update()

}

func (ps *PlayState) init() {
	ps.gdl = graphics.NewGraphicsDataLoader("res/play")
	ps.world = gameplay.NewWorld(ps.gdl, ps.im)
	ps.inited = true
}

func (ps *PlayState) Draw(screen *ebiten.Image) {
	if !ps.inited {
		return
	}
	ps.world.Draw(screen)
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
