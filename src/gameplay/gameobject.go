package gameplay

import (
	"github.com/Jack-Craig/gogame/src/input"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game objects are anything that has a texture and location
type GameObject struct {
	id                  uint32
	x, y, width, height float32
	im                  *ebiten.Image
	w                   *World
}

func NewGameObject(id uint32, x, y, width, height float32, w *World, im *ebiten.Image) *GameObject {
	return &GameObject{
		id, x, y, width, height, im, w,
	}
}

func (gobj *GameObject) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	geo := gobj.w.camera.GetRenderOffset()
	geo.Translate(float64(gobj.x), float64(gobj.y))
	op.GeoM = geo
	screen.DrawImage(gobj.im, &op)
}

// Tiles are game objects with collision, the world is made of tiles
type Tile struct {
	GameObject
	isPassable bool
}

func NewTile(id uint32, x, y float32, w *World, im *ebiten.Image) *Tile {
	return &Tile{
		GameObject{id, x, y, TILEWIDTH, TILEWIDTH, im, w},
		false,
	}
}

// Entities are similar to game objects but also have movement
type Entity struct {
	GameObject
	vx, vy           float32
	stayWithinCamera bool
}

func (e *Entity) Update() {
	// Limit velocity
	if e.vx > 0 && e.vx > MAXVEL {
		e.vx = MAXVEL
	} else if e.vx < 0 && e.vx < -MAXVEL {
		e.vx = -MAXVEL
	}
	if e.vy > 0 && e.vy > MAXVEL {
		e.vy = MAXVEL
	} else if e.vy < 0 && e.vy < -MAXVEL {
		e.vy = -MAXVEL
	}
	// World collision
	expectedX := e.x + e.vx
	if e.vx > 0 {
		expectedX += e.width
	}
	collisionX := e.w.IsWorldCollision(expectedX, e.y) || e.w.IsWorldCollision(expectedX, e.y+e.height)
	if e.stayWithinCamera {
		collisionX = collisionX || !e.w.camera.IsInsideCamera(expectedX, e.y) || !e.w.camera.IsInsideCamera(expectedX, e.y+e.height)
	}
	if collisionX {
		e.vx = 0
	} else {
		e.x += e.vx
	}
	expectedY := e.y + e.vy
	if e.vy > 0 {
		expectedY += e.height
	}
	collisionY := e.w.IsWorldCollision(e.x, expectedY) || e.w.IsWorldCollision(e.x+e.width, expectedY)
	if e.stayWithinCamera {
		collisionY = collisionY || !e.w.camera.IsInsideCamera(e.x, expectedY) || !e.w.camera.IsInsideCamera(e.x+e.width, expectedY)
	}
	if collisionY {
		e.vy = 0
	} else {
		e.y += e.vy
	}

}

func (e *Entity) AddVel(dx, dy float32) {
	e.vx += dx
	e.vy += dy
}

// Players are entities with controls
type Player struct {
	Entity
	pi *input.PlayerInput
}

func NewPlayer(id uint32, x, y, width, height float32, w *World, im *ebiten.Image, pip *input.PlayerInput) *Player {
	return &Player{
		Entity: Entity{
			GameObject: GameObject{
				id, x, y, width, height, im, w,
			},
			vx:               0,
			vy:               0,
			stayWithinCamera: true,
		},
		pi: pip,
	}
}

func (p *Player) Update() {
	_, xAxis := p.pi.GetAxes()
	var magn float32 = 5
	p.vx = magn * xAxis
	if p.pi.IsButtonPressed(input.JoyConB) && p.w.IsWorldCollision(p.x, p.y+p.height+5) {
		p.vy -= 8
	}
}