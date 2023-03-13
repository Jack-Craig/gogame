package gameplay

import (
	"time"

	"github.com/Jack-Craig/gogame/src/graphics"
	"github.com/Jack-Craig/gogame/src/input"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game objects are anything that has a texture and location
type GameObject struct {
	id                  uint32
	x, y, width, height float32
	im                  *ebiten.Image
	w                   *World
	shouldRemove        bool
}

func NewGameObject(id uint32, x, y, width, height float32, w *World, im *ebiten.Image) *GameObject {
	return &GameObject{
		id, x, y, width, height, im, w, false,
	}
}

func (gobj *GameObject) Draw(screen *ebiten.Image) {
	if gobj.im == nil {
		return
	}
	op := ebiten.DrawImageOptions{}
	camOffX, camOffY := gobj.w.camera.GetRenderOffset()
	w, h := gobj.im.Size()
	op.GeoM.Scale(float64(gobj.width/float32(w)), float64(float32(gobj.height/float32(h))))
	op.GeoM.Translate(float64(camOffX), float64(camOffY))
	op.GeoM.Translate(float64(gobj.x), float64(gobj.y))
	screen.DrawImage(gobj.im, &op)
}

// Tiles are game objects with collision, the world is made of tiles
type Tile struct {
	GameObject
	isPassable bool
}

func NewTile(id uint32, x, y float32, w *World, im *ebiten.Image) *Tile {
	return &Tile{
		GameObject{id, x, y, TILEWIDTH, TILEWIDTH, im, w, false},
		false,
	}
}

// Entities are similar to game objects but also have movement
type Entity struct {
	GameObject
	vx, vy                    float32
	stayWithinCamera          bool
	health, gravityMultiplier float32
	// Maintained by world every Update()
	collidingEntities []*Entity
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

	// TOOD: Move to gun struct
	fireRate     int64 // Milliseconds
	lastShotTime int64 // millseconds
}

func NewPlayer(id uint32, x, y, width, height float32, w *World, im *ebiten.Image, pip *input.PlayerInput) *Player {
	return &Player{
		Entity: Entity{
			GameObject: GameObject{
				id, x, y, width, height, im, w, false,
			},
			vx:                0,
			vy:                0,
			stayWithinCamera:  true,
			gravityMultiplier: 1,
		},
		pi:       pip,
		fireRate: 150,
	}
}

func (p *Player) Update() {
	if !p.pi.IsButtonPressed(input.JoyConTriggerLeft) {
		_, xAxis := p.pi.GetAxes()
		var magn float32 = 5
		p.vx = magn * xAxis
	} else {
		p.vx = 0
	}
	if p.pi.IsButtonPressed(input.JoyConB) && (p.w.IsWorldCollision(p.x, p.y+p.height+2) || p.w.IsWorldCollision(p.x+p.width, p.y+p.height+2)) {
		p.vy -= 8.5
	}

	if p.pi.IsButtonPressed(input.JoyConA) {
		p.Shoot()
	}
}

// TOOD: Move to gun object
func (p *Player) Shoot() {
	curTime := time.Now().UnixMilli()
	if p.lastShotTime < curTime-p.fireRate {
		p.lastShotTime = curTime
		p := NewBullet(p.x, p.y+p.height/3, 12, 0, 10, p.w)
		p.w.AddProjectile(p)
	}
}

type Projectile struct {
	Entity
	damage float32
}

func NewProjectile(id uint32, x, y, width, height, vx, vy, damage float32, w *World, im *ebiten.Image) *Projectile {
	return &Projectile{
		Entity: Entity{
			GameObject: GameObject{
				id: id, x: x, y: y, width: width, height: height, im: im, w: w,
			},
			vx:                vx,
			vy:                vy,
			stayWithinCamera:  false,
			health:            0,
			collidingEntities: nil,
		},
		damage: damage,
	}
}

func NewBullet(x, y, vx, vy, damage float32, w *World) *Projectile {
	return NewProjectile(0, x, y, 15, 10, vx, vy, damage, w, w.gdl.GetSpriteImage(graphics.Bullet))
}
