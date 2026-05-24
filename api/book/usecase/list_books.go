package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"bookstore/api/book/domain"
	"bookstore/api/book/repository"
)

// ErrInvalidCursor はカーソルのデコードに失敗した場合のエラー。
var ErrInvalidCursor = errors.New("invalid cursor")

const (
	defaultLimit = 20
	maxLimit     = 100
)

// pageCursor はカーソルのエンコード形式。
type pageCursor struct {
	CreatedAt time.Time `json:"ca"`
	ID        uuid.UUID `json:"id"`
}

// ListBooksInput は蔵書一覧取得のユースケース入力。
type ListBooksInput struct {
	Limit     int
	CursorStr *string
}

// ListBooksOutput は蔵書一覧取得のユースケース出力。
type ListBooksOutput struct {
	Books      []*domain.Book
	NextCursor *string
}

// ListBooksUseCase は蔵書一覧取得のユースケース。
type ListBooksUseCase struct {
	repo repository.BookRepository
}

func NewListBooksUseCase(repo repository.BookRepository) *ListBooksUseCase {
	return &ListBooksUseCase{repo: repo}
}

func (uc *ListBooksUseCase) Execute(ctx context.Context, input ListBooksInput) (*ListBooksOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	var cursor *repository.PageCursor
	if input.CursorStr != nil && *input.CursorStr != "" {
		pc, err := decodeCursor(*input.CursorStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidCursor, err)
		}
		cursor = &repository.PageCursor{
			CreatedAt: pc.CreatedAt,
			ID:        pc.ID,
		}
	}

	books, err := uc.repo.ListBooks(ctx, repository.ListBooksParams{
		Limit:  int32(limit),
		Cursor: cursor,
	})
	if err != nil {
		return nil, fmt.Errorf("list books: %w", err)
	}

	// 取得件数が limit と一致する場合、次ページが存在する可能性がある。
	var nextCursor *string
	if len(books) == limit {
		last := books[len(books)-1]
		encoded := encodeCursor(pageCursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID,
		})
		nextCursor = &encoded
	}

	return &ListBooksOutput{
		Books:      books,
		NextCursor: nextCursor,
	}, nil
}

func encodeCursor(c pageCursor) string {
	b, _ := json.Marshal(c)
	return base64.URLEncoding.EncodeToString(b)
}

func decodeCursor(s string) (pageCursor, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return pageCursor{}, err
	}
	var c pageCursor
	if err := json.Unmarshal(b, &c); err != nil {
		return pageCursor{}, err
	}
	if c.ID == uuid.Nil {
		return pageCursor{}, errors.New("cursor id is zero")
	}
	return c, nil
}
