package core

import (
	"github.com/playdate-go/pdgo"
)

// Explosion represents an animated explosion effect
type Explosion struct {
	game *Game

	x, y   float32
	frame  int
	sprite *pdgo.LCDSprite
	alive  bool

	// Animation timing
	frameDelay   int
	delayCounter int
}

// NewExplosion creates a new explosion at the given position
func NewExplosion(game *Game, x, y int) *Explosion {
	e := &Explosion{
		game:       game,
		x:          float32(x),
		y:          float32(y),
		frame:      0,
		frameDelay: 2,
		alive:      true,
	}

	pd := game.PD()

	// Create sprite for rendering
	e.sprite = pd.Sprite.NewSprite()
	pd.Sprite.SetImage(e.sprite, game.explosionImages[0], pdgo.BitmapUnflipped)

	// Center on position
	if game.explosionImages[0] != nil {
		data := pd.Graphics.GetBitmapData(game.explosionImages[0])
		e.x -= float32(data.Width) / 2
		e.y -= float32(data.Height) / 2
	}

	pd.Sprite.MoveTo(e.sprite, e.x, e.y)
	pd.Sprite.SetZIndex(e.sprite, 2000) // On top of everything
	pd.Sprite.AddSprite(e.sprite)

	return e
}

// Update advances the explosion animation. Returns false when done.
func (e *Explosion) Update() bool {
	if !e.alive {
		return false
	}

	e.delayCounter++
	if e.delayCounter >= e.frameDelay {
		e.delayCounter = 0
		e.frame++

		if e.frame >= 8 {
			e.Kill()
			return false
		}

		// Update sprite image to next frame
		e.game.pd.Sprite.SetImage(e.sprite, e.game.explosionImages[e.frame], pdgo.BitmapUnflipped)
	}

	return true
}

// Kill removes the explosion
func (e *Explosion) Kill() {
	if e.alive {
		e.alive = false
		pd := e.game.pd
		pd.Sprite.RemoveSprite(e.sprite)
		pd.Sprite.FreeSprite(e.sprite)
	}
}
