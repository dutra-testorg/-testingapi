package micro

import (
	"github.com/pkg/errors"
)

var (
	// ErrInvalidID bad request due to an invalid UUIDv4
	ErrInvalidID = errors.New("invalid id (must be an UUIDv4)")
	// ErrInternal generic internal error
	ErrInternal = errors.New("internal error")
)
