package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// Render at roughly 60 fps
const msPerFrame = 16

const width int32 = 1100
const height int32 = 800
const playerWidth int32 = 20
const playerHeight int32 = 150

type gameObjects struct {
	Player1 *Player
	Player2 *Player
	Ball    *Ball
	Walls   *Walls
}

func clearFrame(renderer *sdl.Renderer) {
	err := renderer.SetDrawColor(0x00, 0x00, 0x00, 0xFF)
	if err != nil {
		panic(err)
	}
	renderer.Clear()
}

func drawFrame(renderer *sdl.Renderer, objs *gameObjects) {
	objs.Walls.draw(renderer)
	objs.Ball.draw(renderer)
	objs.Player1.draw(renderer)
	objs.Player2.draw(renderer)
	renderer.Present()
}

func createGameObjects(renderer *sdl.Renderer) *gameObjects {
	player1 := CreatePlayer1(renderer)
	player2 := CreatePlayer2(renderer)
	ball := createBall(renderer)
	walls := createWalls()
	gameObjects := gameObjects{Player1: &player1, Player2: &player2, Ball: &ball, Walls: &walls}
	return &gameObjects
}

func updateOjects(objs *gameObjects) {
	objs.Ball.update()
	objs.Player1.update()
	objs.Player2.update()
}

//TODO: better place for this?
func didCollide(first *sdl.Rect, second *sdl.Rect) bool {
	colliding := true
	firstTop := first.Y
	firstBottom := first.Y + first.H
	firstLeft := first.X
	firstRight := first.X + first.W

	secondTop := second.Y
	secondBottom := second.Y + second.H
	secondLeft := second.X
	secondRight := second.X + second.W

	if firstBottom < secondTop || firstTop > secondBottom {
		colliding = false
	}
	if firstRight < secondLeft || firstLeft > secondRight {
		colliding = false
	}

	return colliding
}

func handleCollisions(objs *gameObjects) {
	ballEntity := objs.Ball.Entity
	player1Entity := objs.Player1.Entity
	player2Entity := objs.Player2.Entity
	if didCollide(ballEntity.Rect, player1Entity.Rect) || didCollide(ballEntity.Rect, player2Entity.Rect) {
		objs.Ball.Xvelocity *= -1
	}
	if didCollide(ballEntity.Rect, &objs.Walls.BottomWall) || didCollide(ballEntity.Rect, &objs.Walls.TopWall) {
		objs.Ball.Yvelocity *= -1
	}
	//TODO handle collisions with right and left walls (goals)
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
		handleCollisions(gameObjects)
		updateOjects(gameObjects)
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
