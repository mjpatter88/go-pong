package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const width = 800
const height = 600
const rectWidth = 20
const rectHeight = 20
const numCols = width / rectWidth
const numRows = height / rectHeight

type entity struct {
	Rect    *sdl.Rect
	Texture *sdl.Texture
}

func clearFrame(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0x00, 0x00, 0x00, 0xFF)
	renderer.Clear()
}

func drawFrame(renderer *sdl.Renderer, entities []entity) {
	for _, entity := range entities {
		err := renderer.Copy(entity.Texture, nil, entity.Rect)
		if err != nil {
			panic(err)
		}
	}
	renderer.Present()
}

func createBlock(renderer *sdl.Renderer, x int, y int, w int, h int, r uint8, g uint8, b uint8) entity {

	tex, err := renderer.CreateTexture(
		uint32(sdl.PIXELFORMAT_RGBA32),
		sdl.TEXTUREACCESS_STREAMING,
		int32(w),
		int32(h),
	)
	if err != nil {
		panic(err)
	}

	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}

	// Ignore the pitch for now
	bytes, _, err := tex.Lock(nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < w*h; i++ {
		bytes[i*4] = r
		bytes[i*4+1] = g
		bytes[i*4+2] = b
		bytes[i*4+3] = 0xFF
	}
	tex.Unlock()
	return entity{Rect: &rect, Texture: tex}
}

func createEntites(renderer *sdl.Renderer) []entity {
	var entities [numRows * numCols]entity

	for row := 0; row < numRows; row++ {
		y := rectHeight * row
		for col := 0; col < numCols; col++ {
			x := rectWidth * col
			var r, g, b uint8
			if (row+col)%2 == 0 {
				r = 77
				g = 200
				b = 233
			} else {
				r = 65
				g = 132
				b = 164
			}
			block := createBlock(renderer, x, y, rectWidth, rectHeight, r, g, b)
			entities[row*numCols+col] = block
		}
	}

	return entities[:]
}

func initialize() (*sdl.Window, *sdl.Renderer) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, renderer, err := sdl.CreateWindowAndRenderer(
		width,
		height,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(err)
	}
	window.SetTitle("Go Pong")
	if err != nil {
		panic(err)
	}
	return window, renderer
}

func main() {
	window, renderer := initialize()
	defer sdl.Quit()
	defer window.Destroy()
	defer renderer.Destroy()

	var frameCount int = 0
	framesProcessed := 0
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				frames := frameCount - framesProcessed
				framesProcessed = frameCount
				fmt.Println("fps: ", frames)
			}
		}
	}()

	entities := createEntites(renderer)

	running := true
	for running {
		frameStart := time.Now()
		clearFrame(renderer)
		drawFrame(renderer, entities)
		frameCount++
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Quit")
				running = false
				done <- false
				break
			}
		}
		// Render at roughly 60 fps
		elapsed := time.Since(frameStart).Milliseconds()
		if elapsed < 16 {
			delay := 16 - elapsed
			if delay < 0 {
				fmt.Println(delay)
				panic(delay)
			}
			sdl.Delay(uint32(delay))
		}
	}
}
