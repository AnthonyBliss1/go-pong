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

type Object struct {
	X, Y, W, H int
}

type Paddle struct {
	Object
}

type Ball struct {
	Object
	dxdt int // x velocity
	dydt int // y velocity
}

type Game struct {
	paddle    Paddle
	ball      Ball
	score     int
	highScore int
}

func main() {
	ebiten.SetWindowTitle("Go-Pong")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	paddle := Paddle{
		Object: Object{
			X: 600,
			Y: 200,
			W: 15,
			H: 100,
		},
	}

	ball := Ball{
		Object: Object{
			X: 0,
			Y: 0,
			W: 15,
			H: 15,
		},
		dxdt: ballSpeed,
		dydt: ballSpeed,
	}

	g := &Game{
		paddle: paddle,
		ball:   ball,
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
	return nil
}

func (p *Paddle) MoveOnKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && p.Y <= screenHeight-p.H {
		p.Y += paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && p.Y >= 0 {
		p.Y -= paddleSpeed
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
	if g.ball.X >= screenWidth {
		g.Reset()
	} else if g.ball.X <= 0 {
		g.ball.dxdt = ballSpeed
	} else if g.ball.Y <= 0 {
		g.ball.dydt = ballSpeed
	} else if g.ball.Y+g.ball.W >= screenHeight {
		g.ball.dydt = -ballSpeed
	}
}

func (g *Game) CollideWithPaddle() {
	if g.ball.X >= g.paddle.X-g.paddle.W && g.ball.Y >= g.paddle.Y && g.ball.Y <= g.paddle.Y+g.paddle.H {
		g.ball.dxdt = -g.ball.dxdt
		g.score++
		if g.score > g.highScore {
			g.highScore = g.score
		}
	}
}
