// Package core implements a shoot-em-up game for Playdate
// Uses sprites for rendering but manual update loop for logic
// (TinyGo doesn't support sprite callbacks)
package core

import (
	"github.com/playdate-go/pdgo"
)

const (
	ScreenWidth  = 400
	ScreenHeight = 240
)

// Game holds the game state and coordinates all entities
type Game struct {
	pd *pdgo.PlaydateAPI

	score int

	// Player
	player *Player

	// Entity lists - manually updated each frame
	bullets    []*Bullet
	enemies    []*Enemy
	explosions []*Explosion
	bgPlanes   []*BackgroundPlane

	// Background
	background *Background

	// Enemy spawning
	maxEnemies  int
	enemyImage  *pdgo.LCDBitmap
	enemyWidth  int
	enemyHeight int

	// Cached images
	explosionImages [8]*pdgo.LCDBitmap

	// Random number generator state (XorShift)
	rngState uint32
}

// New creates a new sprite game instance
func New(pd *pdgo.PlaydateAPI) *Game {
	return &Game{
		pd:         pd,
		rngState:   12345,
		maxEnemies: 10,
		bullets:    make([]*Bullet, 0, 32),
		enemies:    make([]*Enemy, 0, 32),
		explosions: make([]*Explosion, 0, 16),
		bgPlanes:   make([]*BackgroundPlane, 0, 16),
	}
}

// Setup initializes the game
func (g *Game) Setup() {
	// Seed RNG with current time
	seconds, ms := g.pd.System.GetSecondsSinceEpoch()
	g.rngState = uint32(seconds) ^ uint32(ms<<16)
	if g.rngState == 0 {
		g.rngState = 12345
	}

	// Preload images
	g.preloadImages()

	// Create background (handles scrolling texture)
	g.background = NewBackground(g)

	// Create player
	g.player = NewPlayer(g, ScreenWidth/2, ScreenHeight-40)
}

// Update runs one frame of the game
func (g *Game) Update() int {
	g.handleInput()

	// Spawn enemies
	g.spawnEnemyIfNeeded()

	// Spawn background planes
	g.spawnBgPlaneIfNeeded()

	// Update all entities
	g.background.Update()
	g.player.Update()
	g.updateBullets()
	g.updateEnemies()
	g.updateExplosions()
	g.updateBgPlanes()

	// Check collisions
	g.checkCollisions()

	// Sprite system handles all drawing (including background via draw callback)
	g.pd.Sprite.UpdateAndDrawSprites()

	// Draw score on top of everything
	g.pd.Graphics.DrawText("Score: "+itoa(g.score), 10, 10)

	return 1
}

func (g *Game) handleInput() {
	current, pushed, _ := g.pd.System.GetButtonState()

	// Fire on A or B press
	if pushed&pdgo.ButtonA != 0 || pushed&pdgo.ButtonB != 0 {
		g.player.Fire()
	}

	// Move player
	g.player.HandleInput(current)

	// Crank controls enemy count
	change := g.pd.System.GetCrankChange()
	if change > 1 {
		g.maxEnemies++
		if g.maxEnemies > 50 {
			g.maxEnemies = 50
		}
	} else if change < -1 {
		g.maxEnemies--
		if g.maxEnemies < 0 {
			g.maxEnemies = 0
		}
	}
}

func (g *Game) spawnEnemyIfNeeded() {
	if len(g.enemies) < g.maxEnemies && g.maxEnemies > 0 {
		if g.RandomN(60) == 0 {
			g.enemies = append(g.enemies, NewEnemy(g))
		}
	}
}

func (g *Game) spawnBgPlaneIfNeeded() {
	if len(g.bgPlanes) < 10 {
		if g.RandomN(120) == 0 {
			g.bgPlanes = append(g.bgPlanes, NewBackgroundPlane(g))
		}
	}
}

func (g *Game) updateBullets() {
	alive := g.bullets[:0]
	for _, b := range g.bullets {
		if b.Update() {
			alive = append(alive, b)
		}
	}
	g.bullets = alive
}

func (g *Game) updateEnemies() {
	alive := g.enemies[:0]
	for _, e := range g.enemies {
		if e.Update() {
			alive = append(alive, e)
		}
	}
	g.enemies = alive
}

func (g *Game) updateExplosions() {
	alive := g.explosions[:0]
	for _, e := range g.explosions {
		if e.Update() {
			alive = append(alive, e)
		}
	}
	g.explosions = alive
}

func (g *Game) updateBgPlanes() {
	alive := g.bgPlanes[:0]
	for _, bp := range g.bgPlanes {
		if bp.Update() {
			alive = append(alive, bp)
		}
	}
	g.bgPlanes = alive
}

func (g *Game) checkCollisions() {
	// Bullets vs Enemies
	for _, b := range g.bullets {
		if !b.alive {
			continue
		}
		for _, e := range g.enemies {
			if !e.alive {
				continue
			}
			if rectsOverlap(b.x, b.y, b.width, b.height, e.x, e.y, e.width, e.height) {
				b.Kill()
				e.Kill()
				g.CreateExplosion(int(e.x+e.width/2), int(e.y+e.height/2))
				g.score++
			}
		}
	}

	// Player vs Enemies
	px, py, pw, ph := g.player.Bounds()
	for _, e := range g.enemies {
		if !e.alive {
			continue
		}
		if rectsOverlap(px, py, pw, ph, e.x, e.y, e.width, e.height) {
			e.Kill()
			g.CreateExplosion(int(e.x+e.width/2), int(e.y+e.height/2))
			g.score--
		}
	}
}

// AddBullet adds a bullet to the game
func (g *Game) AddBullet(b *Bullet) {
	g.bullets = append(g.bullets, b)
}

// CreateExplosion creates an explosion at the given position
func (g *Game) CreateExplosion(x, y int) {
	g.explosions = append(g.explosions, NewExplosion(g, x, y))
}

// PD returns the Playdate API
func (g *Game) PD() *pdgo.PlaydateAPI {
	return g.pd
}

// Random returns a pseudo-random uint32 using XorShift algorithm
func (g *Game) Random() uint32 {
	x := g.rngState
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	g.rngState = x
	return x
}

// RandomN returns a random number in range [0, n)
func (g *Game) RandomN(n int) int {
	if n <= 0 {
		return 0
	}
	return int(g.Random() % uint32(n))
}

func (g *Game) preloadImages() {
	// Enemy image
	g.enemyImage, _ = g.pd.Graphics.LoadBitmap("images/plane1")
	if g.enemyImage != nil {
		data := g.pd.Graphics.GetBitmapData(g.enemyImage)
		g.enemyWidth = data.Width
		g.enemyHeight = data.Height
	}

	// Explosion frames
	g.explosionImages[0], _ = g.pd.Graphics.LoadBitmap("images/explosion/1")
	g.explosionImages[1], _ = g.pd.Graphics.LoadBitmap("images/explosion/2")
	g.explosionImages[2], _ = g.pd.Graphics.LoadBitmap("images/explosion/3")
	g.explosionImages[3], _ = g.pd.Graphics.LoadBitmap("images/explosion/4")
	g.explosionImages[4], _ = g.pd.Graphics.LoadBitmap("images/explosion/5")
	g.explosionImages[5], _ = g.pd.Graphics.LoadBitmap("images/explosion/6")
	g.explosionImages[6], _ = g.pd.Graphics.LoadBitmap("images/explosion/7")
	g.explosionImages[7], _ = g.pd.Graphics.LoadBitmap("images/explosion/8")
}

// Helper: check if two rectangles overlap
func rectsOverlap(x1, y1, w1, h1, x2, y2, w2, h2 float32) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

// Helper: int to string (avoid fmt for TinyGo)
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
