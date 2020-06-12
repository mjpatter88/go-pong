package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// Render at roughly 60 fps
const msPerFrame = 16

const width = 1100
const height = 800
const playerWidth int32 = 20
const playerHeight int32 = 150
const wallWidth int32 = 20

type gameObjects struct {
	Player1 *Player
	Player2 *Player
	Ball    *Ball
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
	drawWalls(renderer)
	objs.Ball.draw(renderer)
	objs.Player1.draw(renderer)
	objs.Player2.draw(renderer)
	renderer.Present()
}

func createGameObjects(renderer *sdl.Renderer) *gameObjects {
	player1 := CreatePlayer1(renderer)
	player2 := CreatePlayer2(renderer)
	ball := createBall(renderer)
	gameObjects := gameObjects{Player1: &player1, Player2: &player2, Ball: &ball}
	return &gameObjects
}

func updateObjectsPosition(objs *gameObjects) {
	objs.Player1.updatePosition()
	objs.Player2.updatePosition()
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

func handleInput(state []uint8, objs *gameObjects) {
	if state[sdl.SCANCODE_W] == 1 {
		objs.Player1.accelerateUp()
	} else if state[sdl.SCANCODE_S] == 1 {
		objs.Player1.accelerateDown()
	} else {
		objs.Player1.stop()
	}

	if state[sdl.SCANCODE_UP] == 1 {
		objs.Player2.accelerateUp()
	} else if state[sdl.SCANCODE_DOWN] == 1 {
		objs.Player2.accelerateDown()
	} else {
		objs.Player2.stop()
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
		drawFrame(renderer, gameObjects)
		frameCount++

		elapsed := time.Since(frameStart).Milliseconds()
		if elapsed < msPerFrame {
			delay := msPerFrame - elapsed
			if delay < 0 {
				fmt.Println(delay)
				panic(delay)
			}
			sdl.Delay(uint32(delay))
		}
	}
}
