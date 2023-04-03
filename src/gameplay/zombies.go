package gameplay

import (
	"math"
	"math/rand"
	"time"

	"github.com/Jack-Craig/gogame/src/graphics"
)

type Zombie struct {
	Entity
	zai ZombieAI
}

func NewZombie(x, y float32, world *World, ai ZombieAI) *Zombie {
	z := &Zombie{}
	ai.Init(z)
	z.zai = ai
	z.health = 100
	z.GameObject = *NewGameObject(10, x, y, TILEWIDTH-1, TILEWIDTH-1, 0, world, world.gdl.GetSpriteImage(graphics.Bullet), true)
	z.facingDir.X = 1
	z.gravityMultiplier = 1
	z.walkAnimation = *world.gdl.GenerateAnimation(graphics.UserWalkFrame1, graphics.UserWalkFrame6)
	z.idleAnimation = *world.gdl.GenerateAnimation(graphics.UserIdleFrame1, graphics.UserIdleFrame3)
	return z
}

func (z *Zombie) Update() {
	z.zai.Update()
}

type ZombieAI interface {
	Init(z *Zombie)
	Update()
}

type BaseZombieAI struct {
	ZombieAI
	z               *Zombie
	p               *Player
	speed           float32
	hearingDistance float32
	attackDistance  float64

	attackCooldown int64
	lastAttack     int64
}

func NewBaseZombie(x, y float32, w *World) *Zombie {
	return NewZombie(x, y, w, &BaseZombieAI{})
}

func (zai *BaseZombieAI) Init(z *Zombie) {
	zai.z = z
	zai.attackDistance = float64(TILEWIDTH / 2)
	zai.attackCooldown = 1000
	zai.speed = 1.5 + float32((rand.Int()%100))/75
	zai.hearingDistance = 10*TILEWIDTH + float32((rand.Int() % (8 * int(TILEWIDTH))))
}

func (zai *BaseZombieAI) Update() {
	if zai.p == nil {
		// Get nearest player
		var nearestPlayer *Player
		nearestDist := float64(-1)
		for _, player := range zai.z.w.players {
			dist := math.Abs(float64(player.x-zai.z.x)) + math.Abs(float64(player.y-zai.z.y))
			if dist > float64(zai.hearingDistance) {
				continue
			}
			if nearestPlayer == nil || dist < nearestDist {
				nearestDist = dist
				nearestPlayer = player
			}
		}
		zai.p = nearestPlayer
	}
	// No target. Maybe roam?
	if zai.p == nil {
		return
	}
	dx := float64(zai.z.x - zai.p.x)
	dy := math.Abs(float64(zai.z.y - zai.p.y))
	if dx < -zai.attackDistance {
		zai.z.vx = zai.speed
	} else if dx > zai.attackDistance {
		zai.z.vx = -zai.speed
	} else {
		zai.z.vx = 0
		// Try attack
		if dy < float64(TILEWIDTH/2) {
			timeNow := time.Now().UnixMilli()
			if zai.lastAttack == 0 {
				zai.lastAttack = timeNow
			}
			if timeNow > zai.lastAttack+zai.attackCooldown {
				zai.p.health -= 25
				zai.lastAttack = timeNow
			}
		}
	}

}
