package micro

import (
	"context"
)

// Repository basic struct to be filled with data storage info
type Repository struct {
}

// NewRepository to receive data storage info
func NewRepository() *Repository {
	return &Repository{}
}

// Demo to retrieve data from storage
func (r *Repository) Demo(ctx context.Context, uid string) (Demo, error) {
	return Demo{ID: uid}, nil
}
