package main

import "github.com/veandco/go-sdl2/sdl"

const wallWidth int32 = 20

type Walls struct {
	LeftWall   sdl.Rect
	RightWall  sdl.Rect
	TopWall    sdl.Rect
	BottomWall sdl.Rect
}

func createWalls() Walls {
	leftWall := sdl.Rect{X: 10, Y: 10, W: wallWidth, H: height - 20}
	rightWall := sdl.Rect{X: width - (wallWidth + 10), Y: 10, W: wallWidth, H: height - 20}
	topWall := sdl.Rect{X: wallWidth + 10, Y: 10, W: width - (2*wallWidth + 2*10), H: 20}
	bottomWall := sdl.Rect{X: wallWidth + 10, Y: height - 30, W: width - (2*wallWidth + 2*10), H: 20}
	return Walls{LeftWall: leftWall, RightWall: rightWall, TopWall: topWall, BottomWall: bottomWall}
}

func (w *Walls) getWalls() []sdl.Rect {
	return []sdl.Rect{w.LeftWall, w.RightWall, w.TopWall, w.BottomWall}
}

func (w *Walls) draw(renderer *sdl.Renderer) {
	err := renderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
	if err != nil {
		panic(err)
	}
	renderer.FillRects(w.getWalls())
}
