package core

import (
	"github.com/playdate-go/pdgo"
)

// Background handles the scrolling background using a sprite with draw callback
type Background struct {
	game *Game

	sprite *pdgo.LCDSprite
	image  *pdgo.LCDBitmap
	y      int
	height int
	speed  int
}

// NewBackground creates a new scrolling background
func NewBackground(game *Game) *Background {
	bg := &Background{
		game:  game,
		speed: 1,
	}

	pd := game.PD()

	// Load background image
	bg.image, _ = pd.Graphics.LoadBitmap("images/background")
	if bg.image != nil {
		data := pd.Graphics.GetBitmapData(bg.image)
		bg.height = data.Height
	} else {
		bg.height = 240
	}

	// Create background sprite with draw callback for tiled rendering
	bg.sprite = pd.Sprite.NewSprite()
	pd.Sprite.SetBounds(bg.sprite, pdgo.PDRect{X: 0, Y: 0, Width: ScreenWidth, Height: ScreenHeight})
	pd.Sprite.SetDrawFunction(bg.sprite, bg.draw)
	pd.Sprite.SetZIndex(bg.sprite, 0) // Behind everything
	pd.Sprite.AddSprite(bg.sprite)

	return bg
}

// Update scrolls the background
func (bg *Background) Update() {
	bg.y += bg.speed
	if bg.y >= bg.height {
		bg.y = 0
	}
	// Mark sprite dirty so draw callback is called
	bg.game.pd.Sprite.MarkDirty(bg.sprite)
}

// draw is the sprite draw callback - renders tiled background
func (bg *Background) draw(sprite *pdgo.LCDSprite, bounds, drawRect pdgo.PDRect) {
	if bg.image != nil {
		bg.game.pd.Graphics.DrawBitmap(bg.image, 0, bg.y, pdgo.BitmapUnflipped)
		bg.game.pd.Graphics.DrawBitmap(bg.image, 0, bg.y-bg.height, pdgo.BitmapUnflipped)
	}
}

// BackgroundPlane is a decorative plane in the background (no collision)
type BackgroundPlane struct {
	game *Game

	x, y   float32
	height float32
	speed  float32

	sprite *pdgo.LCDSprite
	alive  bool
}

// NewBackgroundPlane creates a new decorative background plane
func NewBackgroundPlane(game *Game) *BackgroundPlane {
	bp := &BackgroundPlane{
		game:  game,
		speed: 2,
		alive: true,
	}

	pd := game.PD()

	// Load plane image
	image, _ := pd.Graphics.LoadBitmap("images/plane2")
	if image != nil {
		data := pd.Graphics.GetBitmapData(image)
		bp.height = float32(data.Height)
		w := data.Width

		// Start at random X
		bp.x = float32(game.RandomN(ScreenWidth - w))
	}

	bp.y = -bp.height

	// Create sprite for rendering
	bp.sprite = pd.Sprite.NewSprite()
	pd.Sprite.SetImage(bp.sprite, image, pdgo.BitmapUnflipped)
	pd.Sprite.MoveTo(bp.sprite, bp.x, bp.y)
	pd.Sprite.SetZIndex(bp.sprite, 100) // Behind everything else
	pd.Sprite.AddSprite(bp.sprite)

	return bp
}

// Update moves the background plane. Returns false if it should be removed.
func (bp *BackgroundPlane) Update() bool {
	if !bp.alive {
		return false
	}

	bp.y += bp.speed

	if bp.y > ScreenHeight+bp.height {
		bp.Kill()
		return false
	}

	// Update sprite position
	bp.game.pd.Sprite.MoveTo(bp.sprite, bp.x, bp.y)

	return true
}

// Kill removes the background plane
func (bp *BackgroundPlane) Kill() {
	if bp.alive {
		bp.alive = false
		pd := bp.game.pd
		pd.Sprite.RemoveSprite(bp.sprite)
		pd.Sprite.FreeSprite(bp.sprite)
	}
}
