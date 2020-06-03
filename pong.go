package main

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const width = 256
const height = 240

func render(surface *sdl.Surface) {
	pixelFormat := surface.Format
	bitspp := int32(pixelFormat.BitsPerPixel)
	bytespp := int32(pixelFormat.BytesPerPixel)
	color := sdl.MapRGBA(pixelFormat, 0x00, 0xFF, 0x00, 0xFF)
	rectWidth := 20
	rectHeight := 20
	var data [200 * 120]uint32
	for i := 0; i < rectWidth*rectHeight; i++ {
		data[i] = color
	}
	newSurface, err := sdl.CreateRGBSurfaceWithFormatFrom(
		unsafe.Pointer(&data),
		int32(rectWidth),
		int32(rectHeight),
		bitspp,
		bytespp*int32(rectWidth),
		pixelFormat.Format,
	)
	if err != nil {
		panic(err)
	}
	newSurface.Blit(nil, surface, &sdl.Rect{X: 20, Y: 20, W: int32(rectWidth), H: int32(rectHeight)})
}
func clearFrame(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0x00, 0x00, 0x00, 0xFF)
	renderer.Clear()
}

func drawFrame(renderer *sdl.Renderer) {
	rectWidth := 20
	rectHeight := 20

	tex, err := renderer.CreateTexture(
		uint32(sdl.PIXELFORMAT_RGBA32),
		sdl.TEXTUREACCESS_STREAMING,
		int32(rectWidth),
		int32(rectHeight),
	)
	if err != nil {
		panic(err)
	}
	defer tex.Destroy()

	rect := sdl.Rect{X: 0, Y: 0, W: int32(rectWidth), H: int32(rectHeight)}

	bytes, _, err := tex.Lock(nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < rectHeight*rectWidth; i++ {
		bytes[i*4] = 77
		bytes[i*4+1] = 200
		bytes[i*4+2] = 233
		bytes[i*4+3] = 0xFF
	}
	tex.Unlock()
	err = renderer.Copy(tex, nil, &rect)
	if err != nil {
		panic(err)
	}
	renderer.Present()
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, renderer, err := sdl.CreateWindowAndRenderer(
		width,
		height,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer renderer.Destroy()
	window.SetTitle("Go Pong")
	if err != nil {
		panic(err)
	}

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

	running := true
	for running {
		frameStart := time.Now()
		clearFrame(renderer)
		drawFrame(renderer)
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
