package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const r = 0xFF
const g = 0xFF
const b = 0xFF

const maxPlayerVelocity = 20
const minPlayerY = 10 + wallWidth
const maxPlayerY = (height - 30) - playerHeight
const playerVelocityIncrement = 4
const playerOffset = 10 + wallWidth + 20
const playerY = (height / 2) - (playerHeight / 2)

type Player struct {
	YVelocity int32
	Entity
}

func (p *Player) draw(renderer *sdl.Renderer) {
	p.Entity.draw(renderer)
}

func (p *Player) updatePosition() {
	// Make sure the new position is in the game field
	newY := p.Rect.Y + p.YVelocity
	if newY > maxPlayerY {
		newY = maxPlayerY
	} else if newY < minPlayerY {
		newY = minPlayerY
	}
	p.Rect.Y = newY
}

func (p *Player) accelerateUp() {
	// First zero out velocity in opposite direction if necessary.
	if p.YVelocity > 0 {
		p.YVelocity = 0
	}

	// Check against max velocity
	p.YVelocity -= maxPlayerVelocity
	if p.YVelocity < -maxPlayerVelocity {
		p.YVelocity = -maxPlayerVelocity
	}
}

func (p *Player) accelerateDown() {
	// First zero out velocity in opposite direction if necessary.
	if p.YVelocity < 0 {
		p.YVelocity = 0
	}

	// Check against max velocity
	p.YVelocity += maxPlayerVelocity
	if p.YVelocity > maxPlayerVelocity {
		p.YVelocity = maxPlayerVelocity
	}
}

func (p *Player) stop() {
	p.YVelocity = 0
}

func CreatePlayer1(renderer *sdl.Renderer) Player {
	x := playerOffset
	y := playerY
	return CreatePlayer(renderer, x, y, playerWidth, playerHeight)
}

func CreatePlayer2(renderer *sdl.Renderer) Player {
	x := width - (playerOffset + playerWidth)
	y := playerY
	return CreatePlayer(renderer, x, y, playerWidth, playerHeight)
}

func CreatePlayer(renderer *sdl.Renderer, x int32, y int32, w int32, h int32) Player {
	tex, err := renderer.CreateTexture(
		uint32(sdl.PIXELFORMAT_RGBA32),
		sdl.TEXTUREACCESS_STREAMING,
		w,
		h,
	)
	if err != nil {
		panic(err)
	}

	rect := sdl.Rect{X: x, Y: y, W: w, H: h}

	// Ignore the pitch for now
	bytes, _, err := tex.Lock(nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < int(w*h); i++ {
		bytes[i*4] = r
		bytes[i*4+1] = g
		bytes[i*4+2] = b
		bytes[i*4+3] = 0xFF
	}
	tex.Unlock()

	ent := Entity{Rect: &rect, Texture: tex}
	return Player{YVelocity: 0, Entity: ent}
}
