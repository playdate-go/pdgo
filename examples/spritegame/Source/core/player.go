package core

import (
	"github.com/playdate-go/pdgo"
)

// Player represents the player's plane
type Player struct {
	game *Game

	x, y          float32
	width, height float32

	sprite      *pdgo.LCDSprite
	image       *pdgo.LCDBitmap
	bulletImage *pdgo.LCDBitmap
	bulletW     int
	bulletH     int

	speed float32
}

// NewPlayer creates a new player at the given position
func NewPlayer(game *Game, centerX, centerY int) *Player {
	p := &Player{
		game:  game,
		x:     float32(centerX),
		y:     float32(centerY),
		speed: 4,
	}

	pd := game.PD()

	// Load player image
	p.image, _ = pd.Graphics.LoadBitmap("images/player")
	if p.image != nil {
		data := pd.Graphics.GetBitmapData(p.image)
		p.width = float32(data.Width)
		p.height = float32(data.Height)
	} else {
		p.width = 32
		p.height = 32
	}

	// Load bullet image
	p.bulletImage, _ = pd.Graphics.LoadBitmap("images/doubleBullet")
	if p.bulletImage != nil {
		data := pd.Graphics.GetBitmapData(p.bulletImage)
		p.bulletW = data.Width
		p.bulletH = data.Height
	} else {
		p.bulletW = 8
		p.bulletH = 16
	}

	// Center position
	p.x -= p.width / 2
	p.y -= p.height / 2

	// Create sprite for rendering
	p.sprite = pd.Sprite.NewSprite()
	pd.Sprite.SetImage(p.sprite, p.image, pdgo.BitmapUnflipped)
	pd.Sprite.MoveTo(p.sprite, p.x, p.y)
	pd.Sprite.SetZIndex(p.sprite, 1000)
	pd.Sprite.AddSprite(p.sprite)

	return p
}

// HandleInput processes player input for movement
func (p *Player) HandleInput(buttons pdgo.PDButtons) {
	dx := float32(0)
	dy := float32(0)

	if buttons&pdgo.ButtonUp != 0 {
		dy = -p.speed
	} else if buttons&pdgo.ButtonDown != 0 {
		dy = p.speed
	}
	if buttons&pdgo.ButtonLeft != 0 {
		dx = -p.speed
	} else if buttons&pdgo.ButtonRight != 0 {
		dx = p.speed
	}

	// Apply movement with bounds checking
	p.x += dx
	p.y += dy

	// Clamp to screen bounds
	if p.x < 0 {
		p.x = 0
	}
	if p.x+p.width > ScreenWidth {
		p.x = ScreenWidth - p.width
	}
	if p.y < 0 {
		p.y = 0
	}
	if p.y+p.height > ScreenHeight {
		p.y = ScreenHeight - p.height
	}
}

// Update updates the player sprite position
func (p *Player) Update() {
	p.game.pd.Sprite.MoveTo(p.sprite, p.x, p.y)
}

// Fire shoots a bullet from the player's position
func (p *Player) Fire() {
	bulletX := p.x + p.width/2 - float32(p.bulletW)/2
	bulletY := p.y - float32(p.bulletH)

	b := NewBullet(p.game, bulletX, bulletY, p.bulletImage, p.bulletW, p.bulletH)
	p.game.AddBullet(b)
}

// Bounds returns the player's collision bounds (slightly smaller than image)
func (p *Player) Bounds() (x, y, w, h float32) {
	margin := float32(5)
	return p.x + margin, p.y + margin, p.width - margin*2, p.height - margin*2
}
