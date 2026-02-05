package repository

import (
	"context"
	"fmt"

	"godocapi/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DocumentRepository struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(db *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(ctx context.Context, doc *model.Document) error {
	query := `INSERT INTO documents (id, filename, storage_path, size, content_type, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, doc.ID, doc.Filename, doc.StoragePath, doc.Size, doc.ContentType, doc.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}
	return nil
}

func (r *DocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Document, error) {
	query := `SELECT id, filename, storage_path, size, content_type, created_at FROM documents WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var doc model.Document
	err := row.Scan(&doc.ID, &doc.Filename, &doc.StoragePath, &doc.Size, &doc.ContentType, &doc.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return &doc, nil
}

func (r *DocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM documents WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}

func (r *DocumentRepository) List(ctx context.Context) ([]model.Document, error) {
	query := `SELECT id, filename, storage_path, size, content_type, created_at FROM documents ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	defer rows.Close()

	var docs []model.Document
	for rows.Next() {
		var doc model.Document
		if err := rows.Scan(&doc.ID, &doc.Filename, &doc.StoragePath, &doc.Size, &doc.ContentType, &doc.CreatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}
