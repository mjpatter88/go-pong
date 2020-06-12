package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const ballWidth = 20
const ballHeight = 20
const startingX = (width + ballWidth) / 2
const startingY = (height + ballHeight) / 2
const startingVelocityX = 5
const startingVelocityY = -2

type Ball struct {
	Xvelocity int32
	Yvelocity int32
	Entity
}

func (b *Ball) draw(renderer *sdl.Renderer) {
	b.Entity.draw(renderer)
}

func (b *Ball) update() {
	b.Rect.X += b.Xvelocity
	b.Rect.Y += b.Yvelocity
}

func createBall(renderer *sdl.Renderer) Ball {
	tex, err := renderer.CreateTexture(
		uint32(sdl.PIXELFORMAT_RGBA32),
		sdl.TEXTUREACCESS_STREAMING,
		ballWidth,
		ballHeight,
	)
	if err != nil {
		panic(err)
	}

	rect := sdl.Rect{X: startingX, Y: startingY, W: ballWidth, H: ballHeight}

	// Ignore the pitch for now
	bytes, _, err := tex.Lock(nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < int(ballWidth*ballHeight); i++ {
		bytes[i*4] = r
		bytes[i*4+1] = g
		bytes[i*4+2] = b
		bytes[i*4+3] = 0xFF
	}
	tex.Unlock()

	ent := Entity{Rect: &rect, Texture: tex}
	return Ball{Xvelocity: startingVelocityX, Yvelocity: startingVelocityY, Entity: ent}
}
