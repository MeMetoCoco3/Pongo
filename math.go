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

	// Calculate normal vector from collision point to ball center
	normalX := ball.Center.X - collisionPoint.X
	normalY := ball.Center.Y - collisionPoint.Y

	// Normalize the vector
	length := float32(math.Sqrt(float64(normalX*normalX + normalY*normalY)))
	if length > 0 {
		normalX /= length
		normalY /= length
	}

	// Move ball away from brick
	ball.Center.X = collisionPoint.X + normalX*(RADIUS+1)
	ball.Center.Y = collisionPoint.Y + normalY*(RADIUS+1)

	dotProduct := ball.Direction.X*normalX + ball.Direction.Y*normalY
	ball.Direction.X = ball.Direction.X - 2*dotProduct*normalX
	ball.Direction.Y = ball.Direction.Y - 2*dotProduct*normalY
	ball.KillFrames = 1
}

func (r *Rectangle) Toggle() {
	if r.Team == Yin {
		r.Team = Yan
		return
	}
	r.Team = Yin
}

func DrawLineVec2(x rl.Vector2, y rl.Vector2, c rl.Color) {
	rl.DrawLine(int32(x.X), int32(x.Y), int32(y.X), int32(y.Y), c)
}
