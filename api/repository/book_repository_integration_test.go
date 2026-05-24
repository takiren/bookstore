//go:build integration

package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"bookstore/api/config"
	"bookstore/api/domain"
	"bookstore/api/repository"
)

func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	cfg := config.Load()
	pool, err := pgxpool.New(context.Background(), cfg.DB.DSN())
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("ping db: %v — docker compose up -d を先に実行してください", err)
	}
	return pool
}

func cleanBooks(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), "DELETE FROM books")
	if err != nil {
		t.Fatalf("cleanup: %v", err)
	}
}

func insertBook(t *testing.T, pool *pgxpool.Pool, isbn, title string, status domain.BookStatus) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	err := pool.QueryRow(context.Background(),
		"INSERT INTO books (isbn, title, status) VALUES ($1, $2, $3) RETURNING id",
		isbn, title, string(status),
	).Scan(&id)
	if err != nil {
		t.Fatalf("insert book: %v", err)
	}
	return id
}

func TestBookRepository_ListBooks_Empty(t *testing.T) {
	pool := newTestPool(t)
	defer pool.Close()
	cleanBooks(t, pool)

	repo := repository.NewBookRepository(pool)
	books, err := repo.ListBooks(context.Background(), repository.ListBooksParams{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 0 {
		t.Errorf("got %d books, want 0", len(books))
	}
}

func TestBookRepository_ListBooks_ReturnsBooks(t *testing.T) {
	pool := newTestPool(t)
	defer pool.Close()
	cleanBooks(t, pool)

	insertBook(t, pool, "9784873119038", "Book A", domain.BookStatusUnread)
	insertBook(t, pool, "9784873119039", "Book B", domain.BookStatusRead)

	repo := repository.NewBookRepository(pool)
	books, err := repo.ListBooks(context.Background(), repository.ListBooksParams{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 2 {
		t.Errorf("got %d books, want 2", len(books))
	}
}

func TestBookRepository_ListBooks_LimitRespected(t *testing.T) {
	pool := newTestPool(t)
	defer pool.Close()
	cleanBooks(t, pool)

	for i := range 5 {
		insertBook(t, pool, fmt.Sprintf("9780000000%03d", i), fmt.Sprintf("Book %d", i), domain.BookStatusUnread)
	}

	repo := repository.NewBookRepository(pool)
	books, err := repo.ListBooks(context.Background(), repository.ListBooksParams{Limit: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 3 {
		t.Errorf("got %d books, want 3", len(books))
	}
}

func TestBookRepository_ListBooks_OrderedByCreatedAtDesc(t *testing.T) {
	pool := newTestPool(t)
	defer pool.Close()
	cleanBooks(t, pool)

	for i := range 3 {
		insertBook(t, pool, fmt.Sprintf("9780000000%03d", i), fmt.Sprintf("Book %d", i), domain.BookStatusUnread)
		time.Sleep(20 * time.Millisecond)
	}

	repo := repository.NewBookRepository(pool)
	books, err := repo.ListBooks(context.Background(), repository.ListBooksParams{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 3 {
		t.Fatalf("got %d books, want 3", len(books))
	}
	for i := 1; i < len(books); i++ {
		if books[i-1].CreatedAt.Before(books[i].CreatedAt) {
			t.Errorf("books not ordered by created_at DESC: index %d (%v) before %d (%v)",
				i-1, books[i-1].CreatedAt, i, books[i].CreatedAt)
		}
	}
}

func TestBookRepository_ListBooks_CursorPagination(t *testing.T) {
	pool := newTestPool(t)
	defer pool.Close()
	cleanBooks(t, pool)

	for i := range 5 {
		insertBook(t, pool, fmt.Sprintf("9780000000%03d", i), fmt.Sprintf("Book %d", i), domain.BookStatusUnread)
		time.Sleep(20 * time.Millisecond)
	}

	repo := repository.NewBookRepository(pool)

	page1, err := repo.ListBooks(context.Background(), repository.ListBooksParams{Limit: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(page1) != 3 {
		t.Fatalf("page1: got %d books, want 3", len(page1))
	}

	last := page1[len(page1)-1]
	page2, err := repo.ListBooks(context.Background(), repository.ListBooksParams{
		Limit: 3,
		Cursor: &repository.PageCursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(page2) != 2 {
		t.Errorf("page2: got %d books, want 2", len(page2))
	}

	seen := make(map[uuid.UUID]bool)
	for _, b := range page1 {
		seen[b.ID] = true
	}
	for _, b := range page2 {
		if seen[b.ID] {
			t.Errorf("duplicate book %v found in page2", b.ID)
		}
	}
}

func TestBookRepository_ListBooks_FieldMapping(t *testing.T) {
	pool := newTestPool(t)
	defer pool.Close()
	cleanBooks(t, pool)

	isbn := "9784873119038"
	title := "Go言語プログラミング"
	insertBook(t, pool, isbn, title, domain.BookStatusReading)

	repo := repository.NewBookRepository(pool)
	books, err := repo.ListBooks(context.Background(), repository.ListBooksParams{Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 1 {
		t.Fatalf("got %d books, want 1", len(books))
	}
	b := books[0]
	if b.ISBN != isbn {
		t.Errorf("isbn = %q, want %q", b.ISBN, isbn)
	}
	if b.Title != title {
		t.Errorf("title = %q, want %q", b.Title, title)
	}
	if b.Status != domain.BookStatusReading {
		t.Errorf("status = %q, want %q", b.Status, domain.BookStatusReading)
	}
	if b.ID == uuid.Nil {
		t.Error("id should not be nil UUID")
	}
	if b.CreatedAt.IsZero() {
		t.Error("created_at should not be zero")
	}
}
