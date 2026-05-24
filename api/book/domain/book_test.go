package domain_test

import (
	"testing"

	"bookstore/api/book/domain"
)

func TestBookStatus_Valid(t *testing.T) {
	tests := []struct {
		status domain.BookStatus
		want   bool
	}{
		{domain.BookStatusUnread, true},
		{domain.BookStatusReading, true},
		{domain.BookStatusRead, true},
		{domain.BookStatusLent, true},
		{"invalid", false},
		{"", false},
		{"UNREAD", false},
	}
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.Valid(); got != tt.want {
				t.Errorf("BookStatus(%q).Valid() = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}
