package main

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREENWIDTH    = 600
	SCREENHEIGHT   = 600
	RECTANGLE_SIZE = 30
	RADIUS         = 16
)

type Team int

const (
	Yin Team = iota
	Yan
)

var (
	yinC = rl.Color{R: 72, G: 229, B: 194, A: 255}
	yanC = rl.Color{R: 51, G: 51, B: 51, A: 255}
)

type Ball struct {
	Team      Team
	Center    rl.Vector2
	Direction rl.Vector2
	Killer    bool
}

type Rectangle struct {
	rl.Rectangle
	Team Team
}

func main() {
	rl.InitWindow(SCREENWIDTH, SCREENHEIGHT, "PONGO")
	defer rl.CloseWindow()

	Field := [20][20]Rectangle{}

	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			var team Team
			if j < 10 {
				team = Yin
			} else {
				team = Yan
			}
			Field[i][j] = Rectangle{
				Rectangle: rl.Rectangle{
					X:      float32(j * RECTANGLE_SIZE),
					Y:      float32(i * RECTANGLE_SIZE),
					Width:  RECTANGLE_SIZE,
					Height: RECTANGLE_SIZE,
				},
				Team: team,
			}
		}
	}

	balls := []Ball{
		{
			Center:    rl.Vector2{50, SCREENHEIGHT / 2},
			Direction: rl.Vector2{X: 5, Y: 2},
			Team:      Yin,
			Killer:    false,
		}, {
			Center:    rl.Vector2{SCREENWIDTH - 50, SCREENHEIGHT / 2},
			Direction: rl.Vector2{X: 5, Y: 4},
			Team:      Yan,
			Killer:    false,
		},
	}

	rl.SetTargetFPS(60)

	// ==========================================
	//            GAME LOOP
	// ==========================================

	for !rl.WindowShouldClose() {

		for i := range balls {
			if balls[i].Center.X-RADIUS <= 0 || balls[i].Center.X+RADIUS >= SCREENWIDTH {
				balls[i].Direction.X = -balls[i].Direction.X
			} else if balls[i].Center.Y-RADIUS <= 0 || balls[i].Center.Y+RADIUS >= SCREENHEIGHT {
				balls[i].Direction.Y = -balls[i].Direction.Y
			}
			// I use double Direction in hopes of getting farther from the border.
			balls[i].Center = Adition(balls[i].Center, ProductScalar(2, balls[i].Direction))

		}

		for i := 0; i < 20; i++ {
			for j := 0; j < 20; j++ {
				for b := 0; b < len(balls); b++ {
					if Field[i][j].Team == balls[b].Team || balls[b].Killer {
						continue
					}
					collide, collisionPoint := CheckBallRectangleCollision(
						&balls[b],
						&Field[i][j])
					if !collide {
						continue
					}
					HandleBallBrickCollision(&balls[b], &Field[i][j], collisionPoint)
					balls[b].Killer = true
					go func() {
						time.Sleep(time.Millisecond * 10)
						balls[b].Killer = false
					}()
				}
			}
		}

		rl.BeginDrawing()

		for i := 0; i < 20; i++ {
			for j := 0; j < 20; j++ {
				var color rl.Color

				if Field[i][j].Team == Yin {
					color = yinC
				} else {
					color = yanC
				}

				rl.DrawRectangleRec(Field[i][j].Rectangle, color)
			}
		}
		for i := range balls {
			var color rl.Color

			if balls[i].Team == Yin {
				color = yanC
			} else {
				color = yinC
			}

			rl.DrawCircle(int32(balls[i].Center.X), int32(balls[i].Center.Y), RADIUS, color)
		}

		rl.EndDrawing()
	}
}
