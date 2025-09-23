package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 665
	screenHeight = 480
	ballSpeed    = 3
	paddleSpeed  = 6
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
	paddle    Paddle
	ball      Ball
	blocks    Blocks
	score     int
	highScore int
}

func main() {
	ebiten.SetWindowTitle("Go-Pong")
	ebiten.SetWindowSize(screenWidth, screenHeight)
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

	g := &Game{
		paddle: paddle,
		ball:   ball,
		blocks: blocks,
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
}

func (g *Game) Update() error {
	g.paddle.MoveOnKeyPress()
	g.ball.Move()
	g.CollideWithWall()
	g.CollideWithPaddle()
	g.CollideWithBlock()
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
		g.Reset()
	} else if g.ball.X <= 0 {
		g.ball.dxdt = ballSpeed
	} else if g.ball.Y <= 0 {
		g.ball.dydt = ballSpeed
	} else if g.ball.X+g.ball.W >= screenWidth {
		g.ball.dxdt = -ballSpeed
	}
}

func (g *Game) CollideWithPaddle() {
	if g.ball.X+g.ball.W >= g.paddle.X && g.ball.X <= g.paddle.X+g.paddle.W && g.ball.Y+g.ball.H >= g.paddle.Y && g.ball.Y <= g.paddle.Y+g.paddle.H {
		overlapX := min(g.ball.X+g.ball.W-g.paddle.X, g.paddle.X+g.paddle.W-g.ball.X)
		overlapY := min(g.ball.Y+g.ball.H-g.paddle.Y, g.paddle.Y+g.paddle.H-g.ball.Y)

		if overlapX < overlapY {
			if g.ball.X < g.paddle.X {
				g.ball.X = g.paddle.X - g.ball.W
			} else {
				g.ball.X = g.paddle.X + g.paddle.W
			}
			g.ball.dxdt = -g.ball.dxdt
		} else {
			g.ball.Y = g.paddle.Y - g.ball.H
			g.ball.dydt = -g.ball.dydt
		}
	}
}

func (g *Game) CollideWithBlock() {
	for i := 0; i < 10; i++ {
		if g.ball.X+g.ball.W >= g.blocks.Blocks[i].X && g.ball.X <= g.blocks.Blocks[i].X+g.blocks.Blocks[i].W && g.ball.Y+g.ball.H >= g.blocks.Blocks[i].Y && g.ball.Y <= g.blocks.Blocks[i].Y+g.blocks.Blocks[i].H {
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

// Helper function to return smaller int
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
