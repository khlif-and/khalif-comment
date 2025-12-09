package domain

import (
	"context"
	"time"

)

// --- ENTITIES ---

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StoryUUID string    `gorm:"index;not null" json:"story_id"` // UUID Story dari service khalif-stories
	UserID    string    `gorm:"index;not null" json:"user_id"`  // User ID dari JWT
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// --- INTERFACES ---

// RedisRepository (Tetap dipertahankan untuk Caching)
type RedisRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	DeletePrefix(ctx context.Context, prefix string) error
}

// CommentRepository (Kontrak untuk akses Database)
type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	GetByStoryUUID(ctx context.Context, storyUUID string) ([]Comment, error)
	GetByID(ctx context.Context, id uint) (*Comment, error)
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, id uint) error
}

// CommentUseCase (Kontrak untuk Business Logic)
type CommentUseCase interface {
	// Create: User memposting komentar baru
	Create(ctx context.Context, storyUUID, userID, content string) (*Comment, error)
	
	// GetByStoryUUID: Mengambil semua komentar untuk story tertentu
	GetByStoryUUID(ctx context.Context, storyUUID string) ([]Comment, error)
	
	// Update: Mengedit komentar (hanya pemilik komentar)
	Update(ctx context.Context, id uint, userID, content string) (*Comment, error)
	
	// Delete: Menghapus komentar (hanya pemilik komentar)
	Delete(ctx context.Context, id uint, userID string) error
}

// --- ERRORS ---
// (Opsional: Jika error defined di file errors.go terpisah, bagian ini tidak perlu. 
// Tapi jika ingin disatukan di domain package, bisa taruh sini atau di errors.go)