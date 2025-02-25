package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Adition(a, b rl.Vector2) rl.Vector2 {
	return rl.Vector2{X: a.X + b.X, Y: a.Y + b.Y}
}

func ProductScalar(a int32, b rl.Vector2) rl.Vector2 {
	return rl.Vector2{b.X * float32(a), b.Y * float32(a)}
}

func Substraction(a, b rl.Vector2) rl.Vector2 {
	return rl.Vector2{X: a.X - b.X, Y: a.Y - b.Y}
}

func Magnitude(a rl.Vector2) float64 {
	return math.Sqrt(float64(a.X*a.X + a.Y*a.Y))
}

func Distance(a, b rl.Vector2) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y

	return math.Sqrt(float64(dx*dx + dy*dy))
}

func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func CheckBallRectangleCollision(ball *Ball, brick *Rectangle) (bool, rl.Vector2) {
	// Check point on Rectangle closes to circle
	closestX := Clamp(ball.Center.X, brick.X, brick.X+brick.Width)
	closestY := Clamp(ball.Center.Y, brick.Y, brick.Y+brick.Height)

	// Calculate Distance
	distanceX := ball.Center.X - closestX
	distanceY := ball.Center.Y - closestY

	distanceSquared := float64(distanceX*distanceX + distanceY*distanceY)

	return distanceSquared <= float64(RADIUS*RADIUS), rl.Vector2{X: closestX, Y: closestY}
}

func HandleBallBrickCollision(ball *Ball, brick *Rectangle, collisionPoint rl.Vector2) {
	brick.Toggle()

	// If the ball is up or down the brick, deltaY will be bigger than deltaX
	// Because, think about lines, if the center on X match to the colision point, means, in this case
	// That the center on Y does not (if not they would be one on top of the other)
	deltaX := math.Abs(float64(collisionPoint.X - ball.Center.X))
	deltaY := math.Abs(float64(collisionPoint.Y - ball.Center.Y))

	if deltaX > deltaY {
		ball.Direction.X = -ball.Direction.X
	} else {
		ball.Direction.Y = -ball.Direction.Y
	}
}

func (r *Rectangle) Toggle() {
	if r.Team == Yin {
		r.Team = Yan
		return
	}
	r.Team = Yin
}
