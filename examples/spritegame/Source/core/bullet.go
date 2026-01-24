package core

import (
	"github.com/playdate-go/pdgo"
)

// Bullet represents a player's bullet
type Bullet struct {
	game *Game

	x, y          float32
	width, height float32
	speed         float32

	sprite *pdgo.LCDSprite
	alive  bool
}

// NewBullet creates a new bullet at the given position
func NewBullet(game *Game, x, y float32, image *pdgo.LCDBitmap, w, h int) *Bullet {
	b := &Bullet{
		game:   game,
		x:      x,
		y:      y,
		width:  float32(w),
		height: float32(h),
		speed:  20,
		alive:  true,
	}

	pd := game.PD()

	// Create sprite for rendering
	b.sprite = pd.Sprite.NewSprite()
	pd.Sprite.SetImage(b.sprite, image, pdgo.BitmapUnflipped)
	pd.Sprite.MoveTo(b.sprite, x, y)
	pd.Sprite.SetZIndex(b.sprite, 999)
	pd.Sprite.AddSprite(b.sprite)

	return b
}

// Update updates the bullet position. Returns false if bullet should be removed.
func (b *Bullet) Update() bool {
	if !b.alive {
		return false
	}

	// Move up
	b.y -= b.speed

	// Remove if off screen
	if b.y < -b.height {
		b.Kill()
		return false
	}

	// Update sprite position
	b.game.pd.Sprite.MoveTo(b.sprite, b.x, b.y)

	return true
}

// Kill removes the bullet
func (b *Bullet) Kill() {
	if b.alive {
		b.alive = false
		pd := b.game.pd
		pd.Sprite.RemoveSprite(b.sprite)
		pd.Sprite.FreeSprite(b.sprite)
	}
}
