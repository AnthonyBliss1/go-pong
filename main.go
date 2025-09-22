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
	screenWidth  = 640
	screenHeight = 480
	ballSpeed    = 3
	paddleSpeed  = 6
)

var (
	Red   = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	Green = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	Blue  = color.RGBA{R: 0, G: 0, B: 255, A: 255}
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
		[]Block{
			Block{
				Object: Object{
					X: (screenWidth / 2) - 25,
					Y: 10,
					W: 50,
					H: 15,
				},
				hits: 0,
			},
		},
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

	switch true {
	case g.blocks.Blocks[0].hits == 0:
		vector.DrawFilledRect(screen,
			float32(g.blocks.Blocks[0].X), float32(g.blocks.Blocks[0].Y),
			float32(g.blocks.Blocks[0].W), float32(g.blocks.Blocks[0].H),
			color.White, false,
		)
	case g.blocks.Blocks[0].hits == 1:
		vector.DrawFilledRect(screen,
			float32(g.blocks.Blocks[0].X), float32(g.blocks.Blocks[0].Y),
			float32(g.blocks.Blocks[0].W), float32(g.blocks.Blocks[0].H),
			Green, false,
		)
	case g.blocks.Blocks[0].hits == 2:
		vector.DrawFilledRect(screen,
			float32(g.blocks.Blocks[0].X), float32(g.blocks.Blocks[0].Y),
			float32(g.blocks.Blocks[0].W), float32(g.blocks.Blocks[0].H),
			Blue, false,
		)
	case g.blocks.Blocks[0].hits == 3:
		vector.DrawFilledRect(screen,
			float32(g.blocks.Blocks[0].X), float32(g.blocks.Blocks[0].Y),
			float32(g.blocks.Blocks[0].W), float32(g.blocks.Blocks[0].H),
			Red, false,
		)
	}

	scoreS := "Score: " + fmt.Sprint(g.score)
	text.Draw(screen, scoreS, basicfont.Face7x13, 10, 10, color.White)

	highScoreS := "High Score: " + fmt.Sprint(g.highScore)
	text.Draw(screen, highScoreS, basicfont.Face7x13, 10, 30, color.White)
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
		g.ball.dydt = -g.ball.dydt
	}
}

func (g *Game) CollideWithBlock() {
	if g.ball.X+g.ball.W >= g.blocks.Blocks[0].X && g.ball.X <= g.blocks.Blocks[0].X+g.blocks.Blocks[0].W && g.ball.Y+g.ball.H >= g.blocks.Blocks[0].Y && g.ball.Y <= g.blocks.Blocks[0].Y+g.blocks.Blocks[0].H {
		g.ball.dydt = -g.ball.dydt
		g.blocks.Blocks[0].hits++
		g.score++
		if g.score > g.highScore {
			g.highScore = g.score
		}
		if g.blocks.Blocks[0].hits == 4 {
			g.blocks.Blocks[0].X = 0
			g.blocks.Blocks[0].Y = screenWidth + g.blocks.Blocks[0].W
		}
		fmt.Println("BLOCK HIT!")
	}
}
