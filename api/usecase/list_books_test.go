package usecase_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"bookstore/api/domain"
	"bookstore/api/repository"
	"bookstore/api/usecase"
)

// mockBookRepository はテスト用のリポジトリモック。
type mockBookRepository struct {
	books            []*domain.Book
	err              error
	calledWithCursor *repository.PageCursor
	calledWithLimit  int32
}

func (m *mockBookRepository) ListBooks(_ context.Context, params repository.ListBooksParams) ([]*domain.Book, error) {
	m.calledWithCursor = params.Cursor
	m.calledWithLimit = params.Limit
	return m.books, m.err
}

func makeBook(id uuid.UUID, createdAt time.Time) *domain.Book {
	return &domain.Book{
		ID:        id,
		ISBN:      "9784873119038",
		Title:     "Test Book",
		Status:    domain.BookStatusUnread,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

func encodeCursorForTest(createdAt time.Time, id uuid.UUID) string {
	type pageCursor struct {
		CreatedAt time.Time `json:"ca"`
		ID        uuid.UUID `json:"id"`
	}
	b, _ := json.Marshal(pageCursor{CreatedAt: createdAt, ID: id})
	return base64.URLEncoding.EncodeToString(b)
}

func TestListBooksUseCase_DefaultLimit(t *testing.T) {
	repo := &mockBookRepository{}
	uc := usecase.NewListBooksUseCase(repo)

	_, err := uc.Execute(context.Background(), usecase.ListBooksInput{})
	if err != nil {
		t.Fatal(err)
	}
	if repo.calledWithLimit != 20 {
		t.Errorf("default limit = %d, want 20", repo.calledWithLimit)
	}
}

func TestListBooksUseCase_MaxLimitClamped(t *testing.T) {
	repo := &mockBookRepository{}
	uc := usecase.NewListBooksUseCase(repo)

	_, err := uc.Execute(context.Background(), usecase.ListBooksInput{Limit: 999})
	if err != nil {
		t.Fatal(err)
	}
	if repo.calledWithLimit != 100 {
		t.Errorf("clamped limit = %d, want 100", repo.calledWithLimit)
	}
}

func TestListBooksUseCase_NoCursor_FirstPage(t *testing.T) {
	repo := &mockBookRepository{}
	uc := usecase.NewListBooksUseCase(repo)

	_, err := uc.Execute(context.Background(), usecase.ListBooksInput{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if repo.calledWithCursor != nil {
		t.Error("cursor should be nil for first page")
	}
}

func TestListBooksUseCase_NextCursor_WhenFullPage(t *testing.T) {
	now := time.Now().UTC()
	books := make([]*domain.Book, 10)
	for i := range books {
		books[i] = makeBook(uuid.New(), now.Add(-time.Duration(i)*time.Hour))
	}

	repo := &mockBookRepository{books: books}
	uc := usecase.NewListBooksUseCase(repo)

	out, err := uc.Execute(context.Background(), usecase.ListBooksInput{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if out.NextCursor == nil {
		t.Error("next_cursor should not be nil when full page returned")
	}
}

func TestListBooksUseCase_NoNextCursor_WhenPartialPage(t *testing.T) {
	books := []*domain.Book{makeBook(uuid.New(), time.Now())}

	repo := &mockBookRepository{books: books}
	uc := usecase.NewListBooksUseCase(repo)

	out, err := uc.Execute(context.Background(), usecase.ListBooksInput{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if out.NextCursor != nil {
		t.Error("next_cursor should be nil when partial page returned")
	}
}

func TestListBooksUseCase_InvalidCursor(t *testing.T) {
	repo := &mockBookRepository{}
	uc := usecase.NewListBooksUseCase(repo)

	invalid := "not-valid-base64!!!"
	_, err := uc.Execute(context.Background(), usecase.ListBooksInput{CursorStr: &invalid})
	if !errors.Is(err, usecase.ErrInvalidCursor) {
		t.Errorf("expected ErrInvalidCursor, got %v", err)
	}
}

func TestListBooksUseCase_ValidCursor_PassedToRepository(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	id := uuid.New()
	encoded := encodeCursorForTest(now, id)

	repo := &mockBookRepository{}
	uc := usecase.NewListBooksUseCase(repo)

	_, err := uc.Execute(context.Background(), usecase.ListBooksInput{CursorStr: &encoded})
	if err != nil {
		t.Fatal(err)
	}
	if repo.calledWithCursor == nil {
		t.Fatal("cursor should be passed to repository")
	}
	if !repo.calledWithCursor.CreatedAt.Equal(now) {
		t.Errorf("cursor.CreatedAt = %v, want %v", repo.calledWithCursor.CreatedAt, now)
	}
	if repo.calledWithCursor.ID != id {
		t.Errorf("cursor.ID = %v, want %v", repo.calledWithCursor.ID, id)
	}
}

func TestListBooksUseCase_RepositoryError_Propagated(t *testing.T) {
	repo := &mockBookRepository{err: errors.New("db down")}
	uc := usecase.NewListBooksUseCase(repo)

	_, err := uc.Execute(context.Background(), usecase.ListBooksInput{})
	if err == nil {
		t.Error("expected error from repository to propagate")
	}
}

func TestListBooksUseCase_EmptyResult(t *testing.T) {
	repo := &mockBookRepository{books: nil}
	uc := usecase.NewListBooksUseCase(repo)

	out, err := uc.Execute(context.Background(), usecase.ListBooksInput{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Books) != 0 {
		t.Errorf("expected 0 books, got %d", len(out.Books))
	}
	if out.NextCursor != nil {
		t.Error("next_cursor should be nil when empty result")
	}
}
