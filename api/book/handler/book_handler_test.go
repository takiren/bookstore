package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"bookstore/api/book/domain"
	"bookstore/api/book/handler"
	"bookstore/api/book/usecase"
	api "bookstore/api/gen/api"
)

type mockListBooksUseCase struct {
	output *usecase.ListBooksOutput
	err    error
	called usecase.ListBooksInput
}

func (m *mockListBooksUseCase) Execute(_ context.Context, input usecase.ListBooksInput) (*usecase.ListBooksOutput, error) {
	m.called = input
	return m.output, m.err
}

func setupEcho() *echo.Echo {
	return echo.New()
}

func TestBookHandler_ListBooks_EmptyResult(t *testing.T) {
	mock := &mockListBooksUseCase{
		output: &usecase.ListBooksOutput{Books: nil, NextCursor: nil},
	}
	h := handler.NewBookHandler(mock)
	strictHandler := api.NewStrictHandler(h, nil)

	e := setupEcho()
	api.RegisterHandlers(e, strictHandler)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}

	var resp api.BookListResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Books) != 0 {
		t.Errorf("books = %d, want 0", len(resp.Books))
	}
	if resp.NextCursor != nil {
		t.Error("next_cursor should be nil")
	}
}

func TestBookHandler_ListBooks_WithBooks(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	id := uuid.New()
	author := "著者A"
	books := []*domain.Book{
		{
			ID:        id,
			ISBN:      "9784873119038",
			Title:     "Go Programming",
			Author:    &author,
			Status:    domain.BookStatusUnread,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	next := "nextcursor"
	mock := &mockListBooksUseCase{
		output: &usecase.ListBooksOutput{Books: books, NextCursor: &next},
	}
	h := handler.NewBookHandler(mock)
	strictHandler := api.NewStrictHandler(h, nil)

	e := setupEcho()
	api.RegisterHandlers(e, strictHandler)

	req := httptest.NewRequest(http.MethodGet, "/books?limit=1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}

	var resp api.BookListResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Books) != 1 {
		t.Fatalf("books = %d, want 1", len(resp.Books))
	}
	if resp.Books[0].Title != "Go Programming" {
		t.Errorf("title = %q, want \"Go Programming\"", resp.Books[0].Title)
	}
	if resp.NextCursor == nil || *resp.NextCursor != "nextcursor" {
		t.Errorf("next_cursor = %v, want \"nextcursor\"", resp.NextCursor)
	}
	if mock.called.Limit != 1 {
		t.Errorf("limit passed to usecase = %d, want 1", mock.called.Limit)
	}
}

func TestBookHandler_ListBooks_InvalidCursor_Returns400(t *testing.T) {
	mock := &mockListBooksUseCase{
		err: fmt.Errorf("%w: bad data", usecase.ErrInvalidCursor),
	}
	h := handler.NewBookHandler(mock)
	strictHandler := api.NewStrictHandler(h, nil)

	e := setupEcho()
	api.RegisterHandlers(e, strictHandler)

	req := httptest.NewRequest(http.MethodGet, "/books?cursor=invalid!!!", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestBookHandler_ListBooks_InternalError_Returns500(t *testing.T) {
	mock := &mockListBooksUseCase{
		err: errors.New("db error"),
	}
	h := handler.NewBookHandler(mock)
	strictHandler := api.NewStrictHandler(h, nil)

	e := setupEcho()
	api.RegisterHandlers(e, strictHandler)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

func TestBookHandler_ListBooks_CursorPassedToUseCase(t *testing.T) {
	cursor := "somecursor"
	mock := &mockListBooksUseCase{
		output: &usecase.ListBooksOutput{},
	}
	h := handler.NewBookHandler(mock)
	strictHandler := api.NewStrictHandler(h, nil)

	e := setupEcho()
	api.RegisterHandlers(e, strictHandler)

	req := httptest.NewRequest(http.MethodGet, "/books?cursor="+cursor, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if mock.called.CursorStr == nil || *mock.called.CursorStr != cursor {
		t.Errorf("cursor passed to usecase = %v, want %q", mock.called.CursorStr, cursor)
	}
}
