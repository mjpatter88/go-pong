package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const width = 1100
const height = 800
const playerWidth int32 = 20
const playerHeight int32 = 150
const wallWidth int32 = 20
const maxPlayerVelocity = 15
const minPlayerY = 10 + wallWidth
const maxPlayerY = (height - 30) - playerHeight

type gameObjects struct {
	Player1 *entity
	Player2 *entity
	Ball    *entity
}

type entity struct {
	Rect      *sdl.Rect
	Texture   *sdl.Texture
	XVelocity int32
	YVelocity int32
}

func clearFrame(renderer *sdl.Renderer) {
	err := renderer.SetDrawColor(0x00, 0x00, 0x00, 0xFF)
	if err != nil {
		panic(err)
	}
	renderer.Clear()
}

func drawWalls(renderer *sdl.Renderer) {
	err := renderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
	if err != nil {
		panic(err)
	}
	var rects [4]sdl.Rect
	// Left, Right, Top, Bottom
	rects[0] = sdl.Rect{X: 10, Y: 10, W: wallWidth, H: height - 20}
	rects[1] = sdl.Rect{X: width - (wallWidth + 10), Y: 10, W: wallWidth, H: height - 20}
	rects[2] = sdl.Rect{X: wallWidth + 10, Y: 10, W: width - (2*wallWidth + 2*10), H: 20}
	rects[3] = sdl.Rect{X: wallWidth + 10, Y: height - 30, W: width - (2*wallWidth + 2*10), H: 20}
	renderer.FillRects(rects[:])
}

func drawFrame(renderer *sdl.Renderer, objs *gameObjects) {
	player1 := objs.Player1
	err := renderer.Copy(player1.Texture, nil, player1.Rect)
	if err != nil {
		panic(err)
	}
	player2 := objs.Player2
	err = renderer.Copy(player2.Texture, nil, player2.Rect)
	if err != nil {
		panic(err)
	}
	renderer.Present()
}

func createPlayer(renderer *sdl.Renderer, x int32, y int32, w int32, h int32, r uint8, g uint8, b uint8) entity {
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
	return entity{Rect: &rect, Texture: tex, XVelocity: 0, YVelocity: 0}
}

func createGameObjects(renderer *sdl.Renderer) *gameObjects {
	playerOffset := 10 + wallWidth + 20
	playerY := (height / 2) - (playerHeight / 2)

	player1 := createPlayer(renderer, playerOffset, playerY, playerWidth, playerHeight, 0xFF, 0xFF, 0xFF)
	player2 := createPlayer(renderer, width-(playerOffset+playerWidth), playerY, playerWidth, playerHeight, 0xFF, 0xFF, 0xFF)

	gameObjects := gameObjects{Player1: &player1, Player2: &player2, Ball: nil}
	return &gameObjects
}

func updateObjectsPosition(objs *gameObjects) {
	player1 := objs.Player1
	player2 := objs.Player2
	// Limit the velocities in each Y direction.
	if player1.YVelocity > maxPlayerVelocity {
		player1.YVelocity = maxPlayerVelocity
	} else if player1.YVelocity < -maxPlayerVelocity {
		player1.YVelocity = -maxPlayerVelocity
	}
	if player2.YVelocity > maxPlayerVelocity {
		player2.YVelocity = maxPlayerVelocity
	} else if player2.YVelocity < -maxPlayerVelocity {
		player2.YVelocity = -maxPlayerVelocity
	}

	newY := player1.Rect.Y + player1.YVelocity
	if newY > maxPlayerY {
		newY = maxPlayerY
	} else if newY < minPlayerY {
		newY = minPlayerY
	}
	player1.Rect.Y = newY

	newY = player2.Rect.Y + player2.YVelocity
	if newY > maxPlayerY {
		newY = maxPlayerY
	} else if newY < minPlayerY {
		newY = minPlayerY
	}
	player2.Rect.Y = newY
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

// TODO: rework this approach make movement nice.
// Maybe apply acceleration to object velocities rather than updating the position directly?
// TODO: prevent moving too far up or down
func handleInput(state []uint8, objs *gameObjects) {
	if state[sdl.SCANCODE_W] == 1 {
		objs.Player1.YVelocity -= 2
	} else if state[sdl.SCANCODE_S] == 1 {
		objs.Player1.YVelocity += 1
	} else {
		objs.Player1.YVelocity = 0
	}

	if state[sdl.SCANCODE_UP] == 1 {
		objs.Player2.YVelocity -= 4
	} else if state[sdl.SCANCODE_DOWN] == 1 {
		objs.Player2.YVelocity += 4
	} else {
		objs.Player2.YVelocity = 0
	}
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

	gameObjects := createGameObjects(renderer)
	running := true
	for running {
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Quit")
				running = false
				done <- true
			}
		}
		handleInput(sdl.GetKeyboardState(), gameObjects)
		updateObjectsPosition(gameObjects)
		clearFrame(renderer)
		drawWalls(renderer)
		drawFrame(renderer, gameObjects)
		frameCount++

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
