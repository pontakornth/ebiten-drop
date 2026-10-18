package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	TICK_RATE   = 1 / 60.0
	GAME_WIDTH  = 320
	GAME_HEIGHT = 240
	TILE_SIZE   = 32
)

type Vector2D struct {
	x float32
	y float32
}

type Game struct {
	bucketX float32
	// TODO: Substitute with objects pools (or very short array)
	drops []Vector2D
	score int
	timer int
}

func (g *Game) InitGameState() {
	g.bucketX = GAME_WIDTH/2 - 32/2
	g.drops = make([]Vector2D, 0)
	// Go variables are zero-init.
}

func (g *Game) Update() error {
	g.timer += 1
	// Around 5 seconds (300 ticks)
	if g.timer > 300 {
		// Spawn drop
		dx := rand.Float32() * (GAME_WIDTH - 32)
		d := Vector2D{x: dx, y: 0}
		g.drops = append(g.drops, d)
		g.timer = 0
	}

	// Drop movement
	for i := len(g.drops) - 1; i >= 0; i-- {
		// Check collision
		x := g.drops[i].x + 16
		y := g.drops[i].y + 16
		distanceSqr := math.Pow(float64(x-g.bucketX), 2) + math.Pow(float64(y-(GAME_HEIGHT-16.0)), 2)
		if distanceSqr <= math.Pow(32, 2) {
			g.drops = g.drops[:len(g.drops)-1]
			g.score += 100
			continue
		}
		g.drops[i].y += 150 * TICK_RATE
	}

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
	// TODO: Use image instead of rect
	c := color.RGBA{R: 125, G: 125, B: 127, A: 255}
	vector.FillRect(screen, g.bucketX, GAME_HEIGHT-32, 32, 32, c, false)
	// TODO: Remove debug circle
	vector.StrokeCircle(screen, g.bucketX+16, GAME_HEIGHT-16, 16, 1, color.White, false)

	// Draw drops
	for _, d := range g.drops {
		// TODO: Use image instead of rect
		c := color.RGBA{R: 100, G: 149, B: 237, A: 255}
		vector.FillRect(screen, d.x, d.y, 32, 32, c, false)
		vector.StrokeCircle(screen, d.x+16, d.y+16, 16, 1, color.RGBA{G: 255, A: 255}, false)
	}

	// TODO: Use actual font
	scoreText := fmt.Sprintf("Score: %03d", g.score)
	ebitenutil.DebugPrint(screen, scoreText)
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
