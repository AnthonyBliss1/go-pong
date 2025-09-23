package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

type GameState int

const (
	screenWidth  = 665
	screenHeight = 480
	ballSpeed    = 4
	paddleSpeed  = 6

	GameStateStart GameState = iota
	GameStatePlaying
	GameStatePaused
	GameStateGameOver
)

type Object struct {
	X, Y, W, H int
}

type Paddle struct {
	Object
}

type Block struct {
	Object
	hits int
}

type Blocks struct {
	Blocks []Block
}

type Ball struct {
	Object
	dxdt int // x velocity
	dydt int // y velocity
}

type Game struct {
	currentGameState GameState
	paddle           Paddle
	ball             Ball
	blocks           Blocks
	score            int
	highScore        int
}

func main() {
	ebiten.SetWindowTitle("Go-Pong")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	paddle, blocks, ball := SetObjects()

	g := &Game{
		currentGameState: GameStateStart,
		paddle:           paddle,
		ball:             ball,
		blocks:           blocks,
	}

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.currentGameState {
	case GameStateStart:
		t1 := "Welcome to BrickBreaker"
		t2 := "Press Space to Start Playing..."
		t3 := "Made by Anthony Bliss"

		text.Draw(screen, t1, basicfont.Face7x13, screenWidth/2-len(t1)*7/2, (screenHeight/2)+10, color.White)
		text.Draw(screen, t2, basicfont.Face7x13, screenWidth/2-len(t2)*7/2, (screenHeight/2)+30, color.White)
		text.Draw(screen, t3, basicfont.Face7x13, screenWidth/2-len(t3)*7/2, screenHeight-10, color.White)

	case GameStatePlaying:
		vector.DrawFilledRect(screen,
			float32(g.paddle.X), float32(g.paddle.Y),
			float32(g.paddle.W), float32(g.paddle.H),
			color.White, false,
		)
		vector.DrawFilledRect(screen,
			float32(g.ball.X), float32(g.ball.Y),
			float32(g.ball.W), float32(g.ball.H),
			color.White, false,
		)

		for i := 0; i < 10; i++ {
			if g.blocks.Blocks[i].hits > 0 {
				vector.DrawFilledRect(screen,
					float32(g.blocks.Blocks[i].X), float32(g.blocks.Blocks[i].Y),
					float32(g.blocks.Blocks[i].W), float32(g.blocks.Blocks[i].H),
					color.White, false,
				)
			}
		}

		scoreS := "Score: " + fmt.Sprint(g.score)
		text.Draw(screen, scoreS, basicfont.Face7x13, 10, (screenHeight/2)+10, color.White)

		highScoreS := "High Score: " + fmt.Sprint(g.highScore)
		text.Draw(screen, highScoreS, basicfont.Face7x13, 10, (screenHeight/2)+30, color.White)

	case GameStateGameOver:
		t1 := "Game Over"
		t2 := "Your Score is " + fmt.Sprint(g.score)
		t3 := "Press Space to Play Again..."

		text.Draw(screen, t1, basicfont.Face7x13, screenWidth/2-len(t1)*7/2, (screenHeight/2)+10, color.White)
		text.Draw(screen, t2, basicfont.Face7x13, screenWidth/2-len(t2)*7/2, (screenHeight/2)+30, color.White)
		text.Draw(screen, t3, basicfont.Face7x13, screenWidth/2-len(t3)*7/2, (screenHeight/2)+50, color.White)
	}
}

func (g *Game) Update() error {
	switch g.currentGameState {
	case GameStateStart:
		g.StartGame()
	case GameStatePlaying:
		g.paddle.MoveOnKeyPress()
		g.ball.Move()
		g.CollideWithWall()
		g.CollideWithPaddle()
		g.CollideWithBlock()
	case GameStateGameOver:
		g.ResetGame()
	}
	return nil
}

func (p *Paddle) MoveOnKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && p.X <= screenWidth-p.W {
		p.X += paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && p.X >= 0 {
		p.X -= paddleSpeed
	}
}

func (b *Ball) Move() {
	b.X += b.dxdt
	b.Y += b.dydt
}

func (g *Game) Reset() {
	g.ball.X = 0
	g.ball.Y = 0

	g.score = 0
}

func (g *Game) CollideWithWall() {
	if g.ball.Y+g.ball.H >= screenHeight {
		g.currentGameState = GameStateGameOver // ball goign through bottom wall
	} else if g.ball.X <= 0 {
		g.ball.dxdt = ballSpeed
	} else if g.ball.Y <= 0 {
		g.ball.dydt = ballSpeed
	} else if g.ball.X+g.ball.W >= screenWidth {
		g.ball.dxdt = -ballSpeed
	}
}

func (g *Game) CollideWithPaddle() {
	if g.ball.X+g.ball.W >= g.paddle.X && g.ball.X <= g.paddle.X+g.paddle.W &&
		g.ball.Y+g.ball.H >= g.paddle.Y && g.ball.Y <= g.paddle.Y+g.paddle.H {

		if g.ball.dydt > 0 {
			g.ball.Y = g.paddle.Y - g.ball.H

			ballCenterX := float64(g.ball.X + g.ball.W/2)
			paddleCenterX := float64(g.paddle.X + g.paddle.W/2)
			paddleHalf := float64(g.paddle.W / 2)

			u := (ballCenterX - paddleCenterX) / paddleHalf

			// clampS
			if u < -1.0 {
				u = -1.0
			} else if u > 1.0 {
				u = 1.0
			}

			maxAngle := 60.0 * math.Pi / 180.0
			theta := u * maxAngle

			speed := math.Sqrt(float64(g.ball.dxdt*g.ball.dxdt + g.ball.dydt*g.ball.dydt))
			if speed < float64(ballSpeed) {
				speed = float64(ballSpeed)
			}

			newDX := speed * math.Sin(theta)
			newDY := -speed * math.Cos(theta)

			g.ball.dxdt = int(math.Round(newDX))
			g.ball.dydt = int(math.Round(newDY))

			if g.ball.dydt == 0 {
				g.ball.dydt = -ballSpeed
			}
			if g.ball.dxdt == 0 && u != 0 {
				if u > 0 {
					g.ball.dxdt = 1
				} else {
					g.ball.dxdt = -1
				}
			}
		}
	}
}

func (g *Game) CollideWithBlock() {
	for i := 0; i < 10; i++ {
		if g.ball.X+g.ball.W >= g.blocks.Blocks[i].X && g.ball.X <= g.blocks.Blocks[i].X+g.blocks.Blocks[i].W &&
			g.ball.Y+g.ball.H >= g.blocks.Blocks[i].Y && g.ball.Y <= g.blocks.Blocks[i].Y+g.blocks.Blocks[i].H {
			overlapX := min(g.ball.X+g.ball.W-g.blocks.Blocks[i].X, g.blocks.Blocks[i].X+g.blocks.Blocks[i].W-g.ball.X)
			overlapY := min(g.ball.Y+g.ball.H-g.blocks.Blocks[i].Y, g.blocks.Blocks[i].Y+g.blocks.Blocks[i].H-g.ball.Y)

			if overlapX < overlapY {
				if g.ball.X < g.blocks.Blocks[i].X {
					g.ball.X = g.blocks.Blocks[i].X - g.ball.W
				} else {
					g.ball.X = g.blocks.Blocks[i].X + g.blocks.Blocks[i].W
				}
				g.ball.dxdt = -g.ball.dxdt
			} else {
				if g.ball.Y < g.blocks.Blocks[i].Y {
					g.ball.Y = g.blocks.Blocks[i].Y - g.ball.H
				} else {
					g.ball.Y = g.blocks.Blocks[i].Y + g.blocks.Blocks[i].H
				}
				g.ball.dydt = -g.ball.dydt
			}
			// instead of subtracting hits, i want to also make the block shrink. think thats better visually
			g.blocks.Blocks[i].H -= 5
			g.blocks.Blocks[i].hits--
			g.score++
			if g.score > g.highScore {
				g.highScore = g.score
			}
			if g.blocks.Blocks[i].hits == 0 {
				g.blocks.Blocks[i].X = 0
				g.blocks.Blocks[i].Y = screenWidth + g.blocks.Blocks[i].W
			}
			fmt.Printf("BLOCK[%d] HIT!\n", i)
		}
	}
}

func (g *Game) StartGame() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.currentGameState = GameStatePlaying
	}
}

func (g *Game) ResetGame() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.paddle, g.blocks, g.ball = SetObjects()

		g.score = 0
		g.highScore = 0

		g.currentGameState = GameStatePlaying
	}
}

func SetObjects() (Paddle, Blocks, Ball) {
	paddle := Paddle{
		Object: Object{
			X: screenWidth / 2,
			Y: screenHeight - 30,
			W: 100,
			H: 15,
		},
	}

	blocks := Blocks{
		[]Block{},
	}

	var c int //tracking positioning of blocksS
	for i := 0; i < 10; i++ {
		c += 15 // add spacing in between blocks
		blocks.Blocks = append(blocks.Blocks, Block{
			Object: Object{
				X: c,
				Y: 10,
				W: 50,
				H: 15,
			},
			hits: 3,
		})
		c += 50 // add width of block
	}

	ball := Ball{
		Object: Object{
			X: screenWidth / 2,
			Y: screenHeight / 2,
			W: 15,
			H: 15,
		},
		dxdt: ballSpeed,
		dydt: ballSpeed,
	}

	return paddle, blocks, ball
}

// Helper function to return smaller int
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
