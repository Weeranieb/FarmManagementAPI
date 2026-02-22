package transaction

import (
	"context"

	"gorm.io/gorm"
)

// Manager runs database work inside a transaction. WithTransaction commits on nil error and rolls back on error or panic.
type Manager interface {
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type manager struct {
	db *gorm.DB
}

// NewManager returns a transaction manager that uses the given DB.
func NewManager(db *gorm.DB) Manager {
	return &manager{db: db}
}

func (m *manager) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return m.db.WithContext(ctx).Transaction(fn)
}
