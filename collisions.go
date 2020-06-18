package main

import "github.com/veandco/go-sdl2/sdl"

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
	//TODO detect collision with either side wall (goal)
	//TODO update the score, etc.
}
