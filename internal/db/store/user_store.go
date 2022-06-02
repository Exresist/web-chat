package store

import (
	"context"
	"database/sql"

	"webChat/internal/ctxkey"
	"webChat/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type userStore struct {
	db        *sql.DB
	tableName string
}

func NewUserStore(db *sql.DB, tableName string) *userStore {
	return &userStore{
		db:        db,
		tableName: tableName,
	}
}

func (u *userStore) GetByID(ctx context.Context, id int) (*model.User, error) {
	query, args, err := sq.Select(
		"users.username",
		"users.hashed_password",
		"users.first_name",
		"users.last_name",
		"users.photo",
		"users.email").
		From(u.tableName).
		Where(sq.Eq{"users.id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	ctxkey.GetLogger(ctx).With("query", query, "args", args).
		Debug("selecting user by id with query")

	user := &model.User{}
	row := u.db.QueryRowContext(ctx, query, args)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if err := row.Scan(user.Username, user.HashedPassword,
		user.FirstName, user.LastName,
		user.Photo, user.Email); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query, args, err := sq.Select(
		"users.username",
		"users.hashed_password",
		"users.first_name",
		"users.last_name",
		"users.photo",
		"users.email").
		From(u.tableName).
		Where(sq.Eq{"users.username": username}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	ctxkey.GetLogger(ctx).With("query", query, "args", args).
		Debug("selecting user by id with query")

	user := &model.User{}
	row := u.db.QueryRowContext(ctx, query, args)
	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if err := row.Scan(user.Username, user.HashedPassword,
		user.FirstName, user.LastName,
		user.Photo, user.Email); err != nil {
		return nil, err
	}
	return user, nil
}
func (u *userStore) Insert(ctx context.Context, user *model.User) error {
	query, args, err := sq.Insert(u.tableName).
		SetMap(map[string]interface{}{
			"username":        user.Username,
			"hashed_password": user.HashedPassword,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"photo":           user.Photo,
			"email":           user.Email,
		}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}
	ctxkey.GetLogger(ctx).With("query", query, "args", args).
		Debug("inserting user")
	if err := u.db.QueryRowContext(ctx, query, args...).Err(); err != nil {
		return err
	}
	return nil
}

func (u *userStore) Update(ctx context.Context, user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userStore) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
