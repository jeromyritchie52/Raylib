package main

import (
	"unsafe"
	"github.com/ctessum/raylib-go/raylib"
)

const (
	// Increase batch sizes to reduce flush frequency
	RL_MAX_BATCH_BUFFERS = 8192  // Increased from default
	RL_MAX_BATCH_DRAWS   = 32768  // Increased from default
)

// Sprite represents a 2D sprite with texture and transform data
type Sprite struct {
	Texture     rl.Texture2D
	Position     rl.Vector2
	SourceRec    rl.Rectangle
	DestRec      rl.Rectangle
	Origin       rl.Vector2
	Rotation     float32
	Color        rl.Color
}

// SpriteBatch groups sprites by texture to minimize state changes
type SpriteBatch struct {
	Texture rl.Texture2D
	Sprites []Sprite
}

// OptimizedSpriteRenderer handles efficient batching of sprites
type OptimizedSpriteRenderer struct {
	batches map[uint32]*SpriteBatch // key: texture ID
}

// NewOptimizedSpriteRenderer creates a new sprite renderer
func NewOptimizedSpriteRenderer() *OptimizedSpriteRenderer {
	return &OptimizedSpriteRenderer{
		batches: make(map[uint22.Texture2D]),
	}
}

// AddSprite adds a sprite to the appropriate batch based on its texture
func (osr *OptimizedSpriteRenderer) AddSprite(sprite Sprite) {
	textureID := *(*uint32)(unsafe.Pointer(&sprite.Texture.ID))
	if _, exists := osr.batches[textureID]; !exists {
		osr.batches[textureID] = &SpriteBatch{
			Texture: sprite.Texture,
			Sprites: make([]Sprite, 0),
		}
	}
	osr.batches[textureID].Sprites = append(osr.batches[textureID].Sprites, sprite)
}

// Draw renders all sprite batches
func (osr *OptimizedSpriteRenderer) Draw() {
	// Draw each batch by texture to minimize state changes
	for _, batch := range osr.batches {
		rl.BeginTextureMode(batch.Texture)
		for _, sprite := range batch.Sprites {
			rl.DrawTexturePro(sprite.Texture, sprite.SourceRec, sprite.DestRec, sprite.Origin, sprite.Rotation, sprite.Color)
		}
		rl.EndTextureMode()
	}
	
	// Clear sprites for next frame
	for _, batch := range osr.batches {
		batch.Sprites = batch.Sprites[:0]
	}
}

// Pre-sorted sprite rendering to minimize texture switches
func DrawSpritesOptimized(sprites []Sprite) {
	if len(sprites) == 0 {
		return
	}
	
	// Group sprites by texture to minimize state changes
	textureGroups := make(map[uint32][]Sprite)
	
	// Group sprites by their texture ID
	for _, sprite := range sprites {
		textureID := *(*uint32)(unsafe.Pointer(&sprite.Texture.ID))
		textureGroups[textureID] = append(textureGroups[textureID], sprite)
	}
	
	// Render each group
	for _, group := range textureGroups {
		if len(group) == 0 {
			continue
		}
		
		// All sprites in this group share the same texture
		firstSprite := group[0]
		rl.BeginTextureMode(firstSprite.Texture)
		
		for _, sprite := range group {
			rl.DrawTexturePro(sprite.Texture, sprite.SourceRec, sprite.DestRec, sprite.Origin, sprite.Rotation, sprite.Color)
		}
		
		rl.EndTextureMode()
	}
}

// Alternative approach: Pre-sort sprites by texture to minimize state changes
func DrawSpritesSortedByTexture(sprites []Sprite) {
	// This would be implemented with a more sophisticated sorting mechanism
	// that groups sprites by texture ID and draws them in batches
}

func main() {
	rl.InitWindow(800, 600, "Optimized 2D Sprite Batching")
	rl.SetTargetFPS(60)
	
	// Create renderer
	renderer := NewOptimizedSpriteRenderer()
	
	// Example usage would go here
	// This is a simplified example structure
	
	rl.CloseWindow()
}