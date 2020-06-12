package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const r = 0xFF
const g = 0xFF
const b = 0xFF

type Player struct {
	XVelocity int32
	YVelocity int32
	Entity
}

func (p *Player) draw(renderer *sdl.Renderer) {
	p.Entity.draw(renderer)
}

func (p *Player) updatePosition() {
	// // Limit the velocities in each Y direction.
	if p.YVelocity > maxPlayerVelocity {
		p.YVelocity = maxPlayerVelocity
	} else if p.YVelocity < -maxPlayerVelocity {
		p.YVelocity = -maxPlayerVelocity
	}

	newY := p.Rect.Y + p.YVelocity
	if newY > maxPlayerY {
		newY = maxPlayerY
	} else if newY < minPlayerY {
		newY = minPlayerY
	}
	p.Rect.Y = newY
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
	return Player{XVelocity: 0, YVelocity: 0, Entity: ent}
}
