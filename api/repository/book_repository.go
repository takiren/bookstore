package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"bookstore/api/db"
	"bookstore/api/domain"
)

// PageCursor はカーソルページネーションの位置を示す。
type PageCursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
}

// ListBooksParams はリポジトリへの検索条件。
type ListBooksParams struct {
	Limit  int32
	Cursor *PageCursor
}

// BookRepository は蔵書データへのアクセスインターフェース。
type BookRepository interface {
	ListBooks(ctx context.Context, params ListBooksParams) ([]*domain.Book, error)
}

type bookRepository struct {
	queries *db.Queries
}

// NewBookRepository は pgxpool を受け取り BookRepository を返す。
func NewBookRepository(pool *pgxpool.Pool) BookRepository {
	return &bookRepository{queries: db.New(pool)}
}

func (r *bookRepository) ListBooks(ctx context.Context, params ListBooksParams) ([]*domain.Book, error) {
	var (
		rows []db.Book
		err  error
	)

	if params.Cursor == nil {
		rows, err = r.queries.ListBooksFirst(ctx, params.Limit)
	} else {
		rows, err = r.queries.ListBooksAfterCursor(ctx, db.ListBooksAfterCursorParams{
			CursorCreatedAt: pgtype.Timestamptz{Time: params.Cursor.CreatedAt, Valid: true},
			CursorID:        params.Cursor.ID,
			LimitCount:      params.Limit,
		})
	}
	if err != nil {
		return nil, fmt.Errorf("query books: %w", err)
	}

	books := make([]*domain.Book, 0, len(rows))
	for _, row := range rows {
		books = append(books, toDomainBook(row))
	}
	return books, nil
}

func toDomainBook(row db.Book) *domain.Book {
	b := &domain.Book{
		ID:           row.ID,
		ISBN:         row.Isbn,
		Title:        row.Title,
		Author:       row.Author,
		Publisher:    row.Publisher,
		ThumbnailURL: row.ThumbnailUrl,
		Status:       domain.BookStatus(row.Status),
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
	if row.PublishedDate.Valid {
		t := row.PublishedDate.Time
		b.PublishedDate = &t
	}
	return b
}
