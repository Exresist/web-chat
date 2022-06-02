package db

import (
	"context"

	"webChat/internal/model"
)

//go:generate mockgen -source=user_store.go -destination=user_store_mock.go -package=store
type UserStore interface {
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Insert(ctx context.Context, customer *model.User) error
	Update(ctx context.Context, customer *model.User) error
	Delete(ctx context.Context, id string) error
}
