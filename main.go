package main

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREENWIDTH    = 600
	SCREENHEIGHT   = 600
	RECTANGLE_SIZE = 30
	RADIUS         = 16
)

type State int

const (
	PAUSE State = iota
	PLAY
	NEWBALL
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

type Game struct {
	State           State
	PressPosition   rl.Vector2
	ReleasePosition rl.Vector2
	balls           []Ball
	field           [20][20]Rectangle
}

type Ball struct {
	Team       Team
	Center     rl.Vector2
	Direction  rl.Vector2
	KillFrames int
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
			Center:     rl.Vector2{50, SCREENHEIGHT / 2},
			Direction:  rl.Vector2{X: 5, Y: 2},
			Team:       Yin,
			KillFrames: 2,
		}, {
			Center:     rl.Vector2{SCREENWIDTH - 50, SCREENHEIGHT / 2},
			Direction:  rl.Vector2{X: 5, Y: 4},
			Team:       Yan,
			KillFrames: 2,
		},
	}
	game := Game{
		PAUSE,
		rl.Vector2{0, 0},
		rl.Vector2{0, 0},
		balls, Field,
	}

	rl.SetTargetFPS(60)

	// ==========================================
	//            GAME LOOP
	// ==========================================

	for !rl.WindowShouldClose() {
		game.GetInput()
		if game.State != PAUSE {
			for i := range game.balls {
				if game.balls[i].Center.X-RADIUS <= 0 || game.balls[i].Center.X+RADIUS >= SCREENWIDTH {
					game.balls[i].Direction.X = -game.balls[i].Direction.X
				} else if game.balls[i].Center.Y-RADIUS <= 0 || game.balls[i].Center.Y+RADIUS >= SCREENHEIGHT {
					game.balls[i].Direction.Y = -game.balls[i].Direction.Y
				}
				// I use double Direction in hopes of getting farther from the border.
				game.balls[i].Center = Adition(game.balls[i].Center, ProductScalar(2, game.balls[i].Direction))

			}

			for i := 0; i < 20; i++ {
				for j := 0; j < 20; j++ {
					for b := 0; b < len(game.balls); b++ {
						if game.field[i][j].Team == game.balls[b].Team || game.balls[b].KillFrames > 0 {
							game.balls[b].KillFrames--
							continue
						}
						collide, collisionPoint := CheckBallRectangleCollision(
							&game.balls[b],
							&game.field[i][j])
						if !collide {
							continue
						}
						HandleBallBrickCollision(&game.balls[b], &game.field[i][j], collisionPoint)
					}
				}
			}
		}

		rl.BeginDrawing()

		for i := 0; i < 20; i++ {
			for j := 0; j < 20; j++ {
				var color rl.Color

				if game.field[i][j].Team == Yin {
					color = yinC
				} else {
					color = yanC
				}

				rl.DrawRectangleRec(game.field[i][j].Rectangle, color)
			}
		}
		for i := range game.balls {
			var color rl.Color

			if game.balls[i].Team == Yin {
				color = yanC
			} else {
				color = yinC
			}

			rl.DrawCircle(int32(game.balls[i].Center.X), int32(game.balls[i].Center.Y), RADIUS, color)
		}

		if game.State == NEWBALL {
			DrawLineVec2(game.PressPosition, game.ReleasePosition, rl.Red)
		}

		rl.EndDrawing()
	}
}

func (g *Game) GetInput() {
	if rl.IsKeyPressed(rl.KeySpace) {
		log.Println("SPACE")
		if g.State == PAUSE {
			g.State = PLAY
		} else {
			g.State = PAUSE
		}
	}
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		g.State = NEWBALL
		g.PressPosition = rl.GetMousePosition()
	}

	if g.State == NEWBALL && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		g.ReleasePosition = rl.GetMousePosition()
		return
	}

	if g.State == NEWBALL {
		newDir := rl.Vector2{
			(g.ReleasePosition.X - g.PressPosition.X) / 100,
			(g.ReleasePosition.Y - g.PressPosition.Y) / 100,
		}
		ix := int(g.PressPosition.X / RECTANGLE_SIZE)
		iy := int(g.PressPosition.Y / RECTANGLE_SIZE)

		team := g.field[iy][ix].Team

		newBall := Ball{
			Center:     g.PressPosition,
			Direction:  newDir,
			Team:       team,
			KillFrames: 1,
		}
		g.balls = append(g.balls, newBall)
		g.State = PLAY
	}
}
