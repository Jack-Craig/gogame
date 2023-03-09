package gameobject

import (
	"github.com/Jack-Craig/gogame/src/input"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game objects are anything that has a texture and location
type GameObject struct {
	id   uint32
	x, y float32
	im   *ebiten.Image
}

func NewGameObject(id uint32, x, y float32, im *ebiten.Image) *GameObject {
	return &GameObject{
		id, x, y, im,
	}
}

func (gobj *GameObject) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(gobj.x), float64(gobj.y))
	screen.DrawImage(gobj.im, op)
}

// Entities are similar to game objects but also have movement
type Entity struct {
	GameObject
	vx, vy float32
}

func (e *Entity) Update() {
	e.x += e.vx
	e.y += e.vy
}

// Players are entities with controls
type Player struct {
	Entity
	pi *input.PlayerInput
}

func NewPlayer(id uint32, x, y float32, im *ebiten.Image, pip *input.PlayerInput) *Player {
	return &Player{
		Entity: Entity{
			GameObject: GameObject{
				id, x, y, im,
			},
			vx: 0,
			vy: 0,
		},
		pi: pip,
	}
}

func (p *Player) Update() {
	yAxis, xAxis := p.pi.GetAxes()
	var magn float32 = 5
	p.vx = magn * xAxis
	p.vy = magn * yAxis
	p.Entity.Update()
}
