package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookStatus string

const (
	BookStatusUnread  BookStatus = "unread"
	BookStatusReading BookStatus = "reading"
	BookStatusRead    BookStatus = "read"
	BookStatusLent    BookStatus = "lent"
)

func (s BookStatus) Valid() bool {
	switch s {
	case BookStatusUnread, BookStatusReading, BookStatusRead, BookStatusLent:
		return true
	}
	return false
}

type Book struct {
	ID            uuid.UUID
	ISBN          string
	Title         string
	Author        *string
	Publisher     *string
	PublishedDate *time.Time
	ThumbnailURL  *string
	Status        BookStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
