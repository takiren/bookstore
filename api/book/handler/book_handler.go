package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"bookstore/api/book/domain"
	"bookstore/api/book/usecase"
	api "bookstore/api/gen/api"
)

// listBooksExecutor は ListBooksUseCase を抽象化するインターフェース（テスト用）。
type listBooksExecutor interface {
	Execute(ctx context.Context, input usecase.ListBooksInput) (*usecase.ListBooksOutput, error)
}

// BookHandler は蔵書に関するハンドラ実装。
type BookHandler struct {
	listBooks listBooksExecutor
}

func NewBookHandler(listBooks listBooksExecutor) *BookHandler {
	return &BookHandler{listBooks: listBooks}
}

// ListBooks は GET /books を処理する。
func (h *BookHandler) ListBooks(ctx context.Context, request api.ListBooksRequestObject) (api.ListBooksResponseObject, error) {
	limit := 0
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}

	out, err := h.listBooks.Execute(ctx, usecase.ListBooksInput{
		Limit:     limit,
		CursorStr: request.Params.Cursor,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCursor) {
			return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return nil, err
	}

	books := make([]api.Book, 0, len(out.Books))
	for _, b := range out.Books {
		books = append(books, toAPIBook(b))
	}

	return api.ListBooks200JSONResponse{
		Books:      books,
		NextCursor: out.NextCursor,
	}, nil
}

func (h *BookHandler) CreateBooks(_ context.Context, _ api.CreateBooksRequestObject) (api.CreateBooksResponseObject, error) {
	return nil, echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}

func (h *BookHandler) DeleteBooks(_ context.Context, _ api.DeleteBooksRequestObject) (api.DeleteBooksResponseObject, error) {
	return nil, echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}

func (h *BookHandler) GetBook(_ context.Context, _ api.GetBookRequestObject) (api.GetBookResponseObject, error) {
	return nil, echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}

func (h *BookHandler) UpdateBook(_ context.Context, _ api.UpdateBookRequestObject) (api.UpdateBookResponseObject, error) {
	return nil, echo.NewHTTPError(http.StatusNotImplemented, "not implemented")
}

func toAPIBook(b *domain.Book) api.Book {
	book := api.Book{
		Id:           b.ID,
		Isbn:         b.ISBN,
		Title:        b.Title,
		Author:       b.Author,
		Publisher:    b.Publisher,
		Status:       api.BookStatus(b.Status),
		ThumbnailUrl: b.ThumbnailURL,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}
	if b.PublishedDate != nil {
		d := openapi_types.Date{Time: *b.PublishedDate}
		book.PublishedDate = &d
	}
	return book
}
