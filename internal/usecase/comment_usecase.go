package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"khalif-comment/internal/domain"

)

type CommentUC struct {
	commentRepo domain.CommentRepository
	redisRepo   domain.RedisRepository
}

func NewCommentUseCase(repo domain.CommentRepository, redis domain.RedisRepository) *CommentUC {
	return &CommentUC{
		commentRepo: repo,
		redisRepo:   redis,
	}
}

// Create: Menambahkan komentar baru
func (uc *CommentUC) Create(ctx context.Context, storyUUID, userID, content string) (*domain.Comment, error) {
	if content == "" {
		return nil, domain.ErrEmptyContent
	}

	comment := &domain.Comment{
		StoryUUID: storyUUID,
		UserID:    userID,
		Content:   content,
	}

	if err := uc.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// Invalidate Cache untuk Story ini agar komentar baru muncul
	if uc.redisRepo != nil {
		cacheKey := fmt.Sprintf("comments:%s", storyUUID)
		_ = uc.redisRepo.Del(ctx, cacheKey)
	}

	return comment, nil
}

// GetByStoryUUID: Mengambil daftar komentar (dengan Caching)
func (uc *CommentUC) GetByStoryUUID(ctx context.Context, storyUUID string) ([]domain.Comment, error) {
	cacheKey := fmt.Sprintf("comments:%s", storyUUID)

	// 1. Cek Redis
	if uc.redisRepo != nil {
		cachedData, err := uc.redisRepo.Get(ctx, cacheKey)
		if err == nil && cachedData != "" {
			var comments []domain.Comment
			if err := json.Unmarshal([]byte(cachedData), &comments); err == nil {
				return comments, nil
			}
		}
	}

	// 2. Ambil dari DB jika cache miss
	comments, err := uc.commentRepo.GetByStoryUUID(ctx, storyUUID)
	if err != nil {
		return nil, err
	}

	// 3. Simpan ke Redis (TTL 10 menit)
	if uc.redisRepo != nil {
		if data, err := json.Marshal(comments); err == nil {
			_ = uc.redisRepo.Set(ctx, cacheKey, data, 10*time.Minute)
		}
	}

	return comments, nil
}

// Update: Mengubah isi komentar (Hanya Pemilik)
func (uc *CommentUC) Update(ctx context.Context, id uint, userID, content string) (*domain.Comment, error) {
	if content == "" {
		return nil, domain.ErrEmptyContent
	}

	// Ambil data existing untuk validasi kepemilikan
	comment, err := uc.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validasi: Apakah yang update adalah pemilik komentar?
	if comment.UserID != userID {
		return nil, domain.ErrUnauthorizedAction
	}

	comment.Content = content
	
	if err := uc.commentRepo.Update(ctx, comment); err != nil {
		return nil, err
	}

	// Invalidate Cache
	if uc.redisRepo != nil {
		cacheKey := fmt.Sprintf("comments:%s", comment.StoryUUID)
		_ = uc.redisRepo.Del(ctx, cacheKey)
	}

	return comment, nil
}

// Delete: Menghapus komentar (Hanya Pemilik)
func (uc *CommentUC) Delete(ctx context.Context, id uint, userID string) error {
	// Ambil data existing untuk tahu StoryUUID (buat hapus cache) & validasi user
	comment, err := uc.commentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Validasi: Apakah yang hapus adalah pemilik komentar?
	if comment.UserID != userID {
		return domain.ErrUnauthorizedAction
	}

	if err := uc.commentRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate Cache
	if uc.redisRepo != nil {
		cacheKey := fmt.Sprintf("comments:%s", comment.StoryUUID)
		_ = uc.redisRepo.Del(ctx, cacheKey)
	}

	return nil
}