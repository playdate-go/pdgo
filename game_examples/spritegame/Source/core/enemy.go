package core

import (
	"github.com/playdate-go/pdgo"
)

// Enemy represents an enemy plane
type Enemy struct {
	game *Game

	x, y          float32
	width, height float32
	speed         float32

	sprite *pdgo.LCDSprite
	alive  bool
}

// NewEnemy creates a new enemy plane at a random position at the top
func NewEnemy(game *Game) *Enemy {
	e := &Enemy{
		game:   game,
		width:  float32(game.enemyWidth),
		height: float32(game.enemyHeight),
		speed:  4,
		alive:  true,
	}

	// Start at random X position at top of screen
	e.x = float32(game.RandomN(int(ScreenWidth - e.width)))
	e.y = -e.height - float32(game.RandomN(30))

	pd := game.PD()

	// Create sprite for rendering
	e.sprite = pd.Sprite.NewSprite()
	pd.Sprite.SetImage(e.sprite, game.enemyImage, pdgo.BitmapUnflipped)
	pd.Sprite.MoveTo(e.sprite, e.x, e.y)
	pd.Sprite.SetZIndex(e.sprite, 500)
	pd.Sprite.AddSprite(e.sprite)

	return e
}

// Update updates the enemy position. Returns false if enemy should be removed.
func (e *Enemy) Update() bool {
	if !e.alive {
		return false
	}

	// Move down
	e.y += e.speed

	// Remove if off screen
	if e.y > ScreenHeight+e.height {
		e.Kill()
		return false
	}

	// Update sprite position
	e.game.pd.Sprite.MoveTo(e.sprite, e.x, e.y)

	return true
}

// Kill removes the enemy
func (e *Enemy) Kill() {
	if e.alive {
		e.alive = false
		pd := e.game.pd
		pd.Sprite.RemoveSprite(e.sprite)
		pd.Sprite.FreeSprite(e.sprite)
	}
}
