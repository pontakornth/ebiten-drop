package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	TICK_RATE   = 1 / 60.0
	GAME_WIDTH  = 320
	GAME_HEIGHT = 240
)

type Vector2D struct {
	x float32
	y float32
}

type Game struct {
	bucketX float32
	drops   []Vector2D
	score   int
	timer   int
}

func (g *Game) InitGameState() {
	g.bucketX = GAME_WIDTH/2 - 32/2
	g.drops = make([]Vector2D, 0)
	// Go variables are zero-init.
}

func (g *Game) Update() error {
	// Bucket movement
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.bucketX += 300 * TICK_RATE
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.bucketX -= 300 * TICK_RATE
	}
	// Mouse movement
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		mx, _ := ebiten.CursorPositionF()
		g.bucketX = float32(mx) - 16
	}
	// clamp
	if g.bucketX < 0 {
		g.bucketX = 0
	}
	if g.bucketX+32 > GAME_WIDTH {
		g.bucketX = GAME_WIDTH - 32
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	// Draw bucket
	c := color.RGBA{R: 125, G: 125, B: 127}
	vector.FillRect(screen, g.bucketX, GAME_HEIGHT-32, 32, 32, c, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	game := &Game{}
	game.InitGameState()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
