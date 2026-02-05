package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"godocapi/internal/model"
	"godocapi/internal/repository"
	"godocapi/internal/storage"

	"github.com/google/uuid"
)

type DocumentService struct {
	repo    *repository.DocumentRepository
	storage storage.Storage
}

func NewDocumentService(repo *repository.DocumentRepository, storage storage.Storage) *DocumentService {
	return &DocumentService{
		repo:    repo,
		storage: storage,
	}
}

func (s *DocumentService) UploadDocument(ctx context.Context, file io.Reader, filename string, size int64, contentType string) (*model.Document, error) {
	// Upload file to storage
	storagePath, err := s.storage.Upload(ctx, file, filename, contentType, size)
	if err != nil {
		return nil, err
	}

	// Create document record
	doc := &model.Document{
		ID:          uuid.New(),
		Filename:    filename,
		StoragePath: storagePath,
		Size:        size,
		ContentType: contentType,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, doc); err != nil {
		// Try to delete file from storage if DB insert fails
		_ = s.storage.Delete(ctx, storagePath)
		return nil, err
	}

	return doc, nil
}

func (s *DocumentService) GetDocument(ctx context.Context, id uuid.UUID) (*model.Document, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DocumentService) DownloadDocument(ctx context.Context, id uuid.UUID) (io.ReadCloser, *model.Document, error) {
	doc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if doc == nil {
		return nil, nil, nil
	}

	stream, err := s.storage.Download(ctx, doc.StoragePath)
	if err != nil {
		return nil, nil, err
	}

	return stream, doc, nil
}

func (s *DocumentService) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	doc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if doc == nil {
		return fmt.Errorf("document not found")
	}

	// Delete from storage
	if err := s.storage.Delete(ctx, doc.StoragePath); err != nil {
		// Log error but continue to delete from DB? or fail?
		// For now we return error, but ideally we should ensure consistency
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete from DB
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *DocumentService) ListDocuments(ctx context.Context) ([]model.Document, error) {
	return s.repo.List(ctx)
}
