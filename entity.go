package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Entity struct {
	Rect    *sdl.Rect
	Texture *sdl.Texture
}

func (e *Entity) draw(renderer *sdl.Renderer) {
	err := renderer.Copy(e.Texture, nil, e.Rect)
	if err != nil {
		panic(err)
	}
}
