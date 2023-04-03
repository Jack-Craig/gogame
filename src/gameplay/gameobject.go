package gameplay

import (
	"math"
	"time"

	"github.com/Jack-Craig/gogame/src/common"
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
	theta               float64
	normals             []common.Vec2
	hasAnimation        bool
}

func NewGameObject(id uint32, x, y, width, height float32, theta float64, w *World, im *ebiten.Image, hasAnimation bool) *GameObject {
	gameObj := &GameObject{
		id, x, y, width, height, im, w, false, theta, nil, hasAnimation,
	}
	gameObj.normals = make([]common.Vec2, 4)
	gameObj.CalcNormals()
	return gameObj
}

func (gobj *GameObject) CalcNormals() {
	projAxis := common.Normalize(common.NewVec2(1, math.Tan(gobj.theta)))

	tl := common.NewVec2(float64(gobj.x), float64(gobj.y))
	tr := common.NewVec2(float64(gobj.x+gobj.width), float64(gobj.y))
	br := common.NewVec2(float64(gobj.x+gobj.width), float64(gobj.y+gobj.height))
	bl := common.NewVec2(float64(gobj.x), float64(gobj.y+gobj.height))

	tEdge := common.Normal(common.Normalize(common.NewVec2(1, common.Dot(projAxis, common.Sub(tl, tr)))))
	rEdge := common.Normal(common.Normalize(common.NewVec2(1, common.Dot(projAxis, common.Sub(tr, br)))))
	bEdge := common.Normal(common.Normalize(common.NewVec2(1, common.Dot(projAxis, common.Sub(br, bl)))))
	lEdge := common.Normal(common.Normalize(common.NewVec2(1, common.Dot(projAxis, common.Sub(bl, tl)))))

	gobj.normals[0] = tEdge
	gobj.normals[1] = rEdge
	gobj.normals[2] = bEdge
	gobj.normals[3] = lEdge
}

func (gobj *GameObject) Draw(screen *ebiten.Image) {
	if gobj.im == nil || gobj.hasAnimation {
		return
	}
	op := ebiten.DrawImageOptions{}
	camOffX, camOffY := gobj.w.camera.GetRenderOffset()
	w, h := gobj.im.Size()

	op.GeoM.Scale(float64(gobj.width/float32(w)), float64(float32(gobj.height/float32(h))))
	op.GeoM.Rotate(gobj.theta)

	op.GeoM.Translate(float64(gobj.x), float64(gobj.y))
	op.GeoM.Translate(float64(camOffX), float64(camOffY))
	screen.DrawImage(gobj.im, &op)
}

// Tiles are game objects with collision, the world is made of tiles
type Tile struct {
	GameObject
	isPassable, isClimbable bool
}

func NewTile(id uint32, x, y float32, w *World, im *ebiten.Image) *Tile {
	return &Tile{
		GameObject{id, x, y, TILEWIDTH, TILEWIDTH, im, w, false, 0, nil, false},
		false,
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
	immuneToGuns      bool
	walkAnimation     graphics.Animation
	idleAnimation     graphics.Animation
	facingDir         common.Vec2
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
	if e.vx > 0 {
		e.facingDir.X = 1
	} else if e.vx < 0 {
		e.facingDir.X = -1
	}
	// World collision
	expectedX := e.x + e.vx
	if e.vx > 0 {
		expectedX += e.width
	}
	collisionX, collisionY := e.WillCollideWithWorld()
	if e.stayWithinCamera {
		if e.vx < 0 || !e.w.canLeave {
			collisionX = collisionX || !e.w.camera.IsInsideCamera(expectedX, e.y) || !e.w.camera.IsInsideCamera(expectedX, e.y+e.height)
		}
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
	if e.stayWithinCamera {
		collisionY = collisionY || !e.w.camera.IsInsideCamera(e.x, expectedY) || !e.w.camera.IsInsideCamera(e.x+e.width, expectedY)
	}
	if collisionY {
		e.vy = 0
	} else {
		e.y += e.vy
	}

	// TODO: Only update when theta changes
	// e.CalcNormals()

}

func (e *Entity) WillCollideWithWorld() (bool, bool) {
	// Check collisions on left top and bottom if vx < 0, else right top and bottom
	// Check collisions on left and right bottom if by > 0, else left and right top
	// X and Y will stay the same, bottom Y and right X will change
	cosTheta, sinTheta := math.Cos(e.theta), math.Sin(e.theta)
	baseX := float64(e.x)
	baseY := float64(e.y)
	tl := common.NewVec2(baseX, baseY)
	tr := common.NewVec2(baseX+float64(e.width)*cosTheta, baseY)
	bl := common.NewVec2(baseX, baseY+float64(e.height)*cosTheta)
	br := common.NewVec2(baseX+float64(e.width)*cosTheta-float64(e.height)*sinTheta, baseY+float64(e.height)*cosTheta+float64(e.width)*sinTheta)
	collisionX := e.w.IsWorldCollision(float32(tl.X)+e.vx, float32(tl.Y)) || e.w.IsWorldCollision(float32(tr.X)+e.vx, float32(tr.Y)) || e.w.IsWorldCollision(float32(bl.X)+e.vx, float32(bl.Y)) || e.w.IsWorldCollision(float32(br.X)+e.vx, float32(br.Y))
	collisionY := e.w.IsWorldCollision(float32(tl.X), float32(tl.Y)+e.vy) || e.w.IsWorldCollision(float32(tr.X), float32(tr.Y)+e.vy) || e.w.IsWorldCollision(float32(bl.X), float32(bl.Y)+e.vy) || e.w.IsWorldCollision(float32(br.X), float32(br.Y)+e.vy)
	return collisionX, collisionY
}

func (e *Entity) Draw(screen *ebiten.Image) {
	if !e.hasAnimation {
		return
	}
	op := ebiten.DrawImageOptions{}
	camOffX, camOffY := e.w.camera.GetRenderOffset()
	//w, h := TILEWIDTH, TILEWIDTH

	op.GeoM.Scale(float64(e.width/float32(24)), float64(float32(e.height/float32(24))))
	if e.facingDir.X < 0 {
		// Facing left
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(e.width), 0)
	}

	op.GeoM.Rotate(e.theta)

	op.GeoM.Translate(float64(e.x), float64(e.y))
	op.GeoM.Translate(float64(camOffX), float64(camOffY))

	// Animation shiz
	if e.vx != 0 {
		e.walkAnimation.Draw(screen, &op)
	} else {
		e.idleAnimation.Draw(screen, &op)
	}
}

func (e *Entity) AddVel(dx, dy float32) {
	e.vx += dx
	e.vy += dy
}

// Players are entities with controls
type Player struct {
	Entity
	pi     *input.PlayerInput
	isDead bool
	name   string
	// TOOD: Move to gun struct
	fireRate     int64 // Milliseconds
	lastShotTime int64 // millseconds
}

func NewPlayer(id uint32, name string, w *World, im *ebiten.Image, pip *input.PlayerInput) *Player {
	return &Player{
		Entity: Entity{
			GameObject:        *NewGameObject(id, 0, 0, TILEWIDTH-1, TILEWIDTH-1, 0, w, im, true),
			vx:                0,
			vy:                0,
			stayWithinCamera:  true,
			gravityMultiplier: 1,
			immuneToGuns:      true,
		},
		pi:       pip,
		fireRate: 100,
		name:     name,
		isDead:   true,
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
		yDir, xDir := p.pi.GetAxes()
		if math.Abs(float64(yDir)) < .05 && math.Abs(float64(xDir)) < .05 {
			yDir = 0
			xDir = float32(p.facingDir.X)
		} else {
			m := float32(math.Sqrt(float64(yDir*yDir + xDir*xDir)))
			yDir /= m
			xDir /= m
		}

		bulletSpeed := float32(30)

		p.lastShotTime = curTime
		p := NewBullet(p.x+p.width/2, p.y+p.height/3, xDir*bulletSpeed, yDir*bulletSpeed, 25, p.w)
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
			GameObject:        *NewGameObject(id, x, y, width, height, math.Atan2(float64(vy), float64(vx)), w, im, false),
			vx:                vx,
			vy:                vy,
			stayWithinCamera:  false,
			health:            1,
			collidingEntities: nil,
			immuneToGuns:      true,
		},
		damage: damage,
	}
}

func NewBullet(x, y, vx, vy, damage float32, w *World) *Projectile {
	return NewProjectile(2, x, y, 18, 4, vx, vy, damage, w, w.gdl.GetSpriteImage(graphics.Bullet))
}

func (p *Projectile) Update() {
	xCol, yCol := p.WillCollideWithWorld()
	p.shouldRemove = xCol || yCol
	for _, e := range p.collidingEntities {
		if e.immuneToGuns {
			continue
		}
		e.health -= p.damage
		p.shouldRemove = true
	}
}
