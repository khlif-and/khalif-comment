package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"khalif-comment/internal/domain"

)

type CommentRepo struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

// Create: Menyimpan komentar baru ke database
func (r *CommentRepo) Create(ctx context.Context, c *domain.Comment) error {
	return r.db.WithContext(ctx).Create(c).Error
}

// GetByStoryUUID: Mengambil semua komentar berdasarkan Story UUID
func (r *CommentRepo) GetByStoryUUID(ctx context.Context, storyUUID string) ([]domain.Comment, error) {
	var comments []domain.Comment
	// Kita urutkan descending (terbaru di atas) agar user langsung melihat komentar baru
	err := r.db.WithContext(ctx).
		Where("story_uuid = ?", storyUUID).
		Order("created_at desc").
		Find(&comments).Error
	return comments, err
}

// GetByID: Mencari satu komentar (penting untuk validasi update/delete)
func (r *CommentRepo) GetByID(ctx context.Context, id uint) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.WithContext(ctx).First(&comment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCommentNotFound // Return error domain yang spesifik
		}
		return nil, err
	}
	return &comment, nil
}

// Update: Menyimpan perubahan konten komentar
func (r *CommentRepo) Update(ctx context.Context, c *domain.Comment) error {
	// Menggunakan Save agar UpdatedAt otomatis diperbarui oleh GORM
	return r.db.WithContext(ctx).Save(c).Error
}

// Delete: Menghapus komentar berdasarkan ID
func (r *CommentRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Comment{}, id).Error
}