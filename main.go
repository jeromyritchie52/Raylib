package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

// Sprite represents a 2D sprite with position, speed, rotation, and color.
type Sprite struct {
	Position rl.Vector2
	Speed    rl.Vector2
	Rotation float32
	Color    rl.Color
}

var (
	customBatch       rl.RlRenderBatch
	customBatchActive bool
)

// SetRenderBatchCapacity configures or dynamically scales the internal rlgl render batch capacity.
// Passing 0 for both parameters resets to the default internal batch.
func SetRenderBatchCapacity(numBuffers int32, maxDraws int32) {
	if customBatchActive {
		rl.RlUnloadRenderBatch(customBatch)
		customBatchActive = false
	}

	if numBuffers > 0 && maxDraws > 0 {
		customBatch = rl.RlLoadRenderBatch(numBuffers, maxDraws)
		rl.RlSetRenderBatchActive(&customBatch)
		customBatchActive = true
	} else {
		rl.RlSetRenderBatchActive(nil)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	rl.InitWindow(screenWidth, screenHeight, "Raylib 2D Sprite Stress Test & Optimization")
	defer rl.CloseWindow()

	// Generate a simple texture atlas / sprite sheet dynamically to avoid external file dependency
	img := rl.GenImageChecked(64, 64, 8, 8, rl.Red, rl.Blue)
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	defer rl.UnloadTexture(tex)

	maxSprites := 25000
	sprites := make([]Sprite, maxSprites)

	for i := 0; i < maxSprites; i++ {
		sprites[i] = Sprite{
			Position: rl.NewVector2(float32(rand.Intn(screenWidth)), float32(rand.Intn(screenHeight))),
			Speed:    rl.NewVector2((rand.Float32()*500-250)/100.0, (rand.Float32()*500-250)/100.0),
			Rotation: rand.Float32() * 360.0,
			Color: rl.NewColor(
				uint8(rand.Intn(206)+50),
				uint8(rand.Intn(206)+50),
				uint8(rand.Intn(206)+50),
				255,
			),
		}
	}

	// Initialize with optimized custom batch capacity (32 buffers, 16384 draws)
	SetRenderBatchCapacity(32, 16384)
	defer func() {
		if customBatchActive {
			rl.RlUnloadRenderBatch(customBatch)
		}
	}()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		// Handle input to toggle batching optimization
		if rl.IsKeyPressed(rl.KeySpace) {
			if customBatchActive {
				SetRenderBatchCapacity(0, 0) // Reset to default
			} else {
				SetRenderBatchCapacity(32, 16384) // Enable optimized batch
			}
		}

		// Adjust sprite count dynamically
		if rl.IsKeyDown(rl.KeyUp) {
			if maxSprites < 100000 {
				maxSprites += 500
				for i := 0; i < 500; i++ {
					sprites = append(sprites, Sprite{
						Position: rl.NewVector2(float32(rand.Intn(screenWidth)), float32(rand.Intn(screenHeight))),
						Speed:    rl.NewVector2((rand.Float32()*500-250)/100.0, (rand.Float32()*500-250)/100.0),
						Rotation: rand.Float32() * 360.0,
						Color: rl.NewColor(
							uint8(rand.Intn(206)+50),
							uint8(rand.Intn(206)+50),
							uint8(rand.Intn(206)+50),
							255,
						),
					})
				}
			}
		}
		if rl.IsKeyDown(rl.KeyDown) {
			if maxSprites > 500 {
				maxSprites -= 500
				sprites = sprites[:maxSprites]
			}
		}

		// Update sprite positions and bounce off screen edges
		for i := 0; i < maxSprites; i++ {
			sprites[i].Position.X += sprites[i].Speed.X
			sprites[i].Position.Y += sprites[i].Speed.Y

			if sprites[i].Position.X < 0 || sprites[i].Position.X > float32(screenWidth) {
				sprites[i].Speed.X *= -1
			}
			if sprites[i].Position.Y < 0 || sprites[i].Position.Y > float32(screenHeight) {
				sprites[i].Speed.Y *= -1
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Draw sprites using DrawTexturePro to support rotation, scaling, and tinting
		sourceRec := rl.NewRectangle(0, 0, float32(tex.Width), float32(tex.Height))
		origin := rl.NewVector2(float32(tex.Width)/2.0, float32(tex.Height)/2.0)

		for i := 0; i < maxSprites; i++ {
			destRec := rl.NewRectangle(sprites[i].Position.X, sprites[i].Position.Y, float32(tex.Width), float32(tex.Height))
			rl.DrawTexturePro(tex, sourceRec, destRec, origin, sprites[i].Rotation, sprites[i].Color)
		}

		// Draw UI / Info overlay
		rl.DrawRectangle(10, 10, 380, 140, rl.Fade(rl.SkyBlue, 0.5))
		rl.DrawFPS(20, 20)
		rl.DrawText(fmt.Sprintf("Sprites: %d", maxSprites), 20, 50, 20, rl.DarkGray)

		if customBatchActive {
			rl.DrawText("Custom Batch: ACTIVE (Optimized)", 20, 80, 20, rl.DarkGreen)
		} else {
			rl.DrawText("Custom Batch: INACTIVE (Default)", 20, 80, 20, rl.Red)
		}

		rl.DrawText("Press SPACE to toggle batching optimization", 20, 110, 16, rl.Gray)
		rl.DrawText("Press UP/DOWN to change sprite count", 20, 130, 16, rl.Gray)

		rl.EndDrawing()
	}
}