package repository

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/weCredit/internal/domain"
)

type pgxUserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &pgxUserRepository{
		db: db,
	}
}

// CreateUser implements domain.UserRepository.
func (r *pgxUserRepository) CreateUser(ctx context.Context, entity *domain.User) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)
	q := `INSERT INTO users(full_name,user_name, role) values($1, $2, $3) RETURNING  id, created_at, updated_at`
	args := []interface{}{entity.FullName, entity.UserName, entity.Role}
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		err = tx.QueryRow(ctx, q, args...).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
	} else {
		err = r.db.QueryRow(ctx, q, args...).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
	}
	return err
}

// DeleteUser implements domain.UserRepository.
func (r *pgxUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	q := `UPDATE users SET deleted_at = NOW() WHERE id = $1`
	args := []interface{}{id}
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		_, err = tx.Exec(ctx, q, args...)
	} else {
		_, err = r.db.Exec(ctx, q, args...)
	}
	return err
}

// FindByID implements domain.UserRepository.
func (r *pgxUserRepository) FindByID(ctx context.Context, id uuid.UUID) (result domain.User, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Retrieve the data
	q := `SELECT * FROM users WHERE id = $1  LIMIT 1`
	args := []interface{}{id}
	var rows pgx.Rows
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		rows, err = tx.Query(ctx, q, args...)
	} else {
		rows, err = r.db.Query(ctx, q, args...)
	}
	defer rows.Close()

	if err != nil {
		return result, err
	}
	if rows == nil {
		return result, err
	}

	result, err = pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[domain.User])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return result, err
	}

	return result, err
}

// FindByUserName implements domain.UserRepository.
func (r *pgxUserRepository) FindByUserName(ctx context.Context, username string) (result domain.User, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Retrieve the data
	q := `SELECT * FROM users WHERE user_name = $1  LIMIT 1`
	args := []interface{}{username}
	var rows pgx.Rows
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		rows, err = tx.Query(ctx, q, args...)
	} else {
		rows, err = r.db.Query(ctx, q, args...)
	}
	defer rows.Close()

	if err != nil {
		return result, err
	}
	if rows == nil {
		return result, err
	}

	result, err = pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[domain.User])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return result, err
	}

	return result, err
}

// UpdateUser implements domain.UserRepository.
func (r *pgxUserRepository) UpdateUser(ctx context.Context, entity *domain.User) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Update the data
	q := `UPDATE users SET full_name = $1, username = $2,  role = $3,  updated_at = NOW() WHERE id = $4 RETURNING updated_at`
	args := []interface{}{entity.FullName, entity.UserName, entity.Role, entity.ID}
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		err = tx.QueryRow(ctx, q, args...).Scan(&entity.UpdatedAt)
	} else {
		err = r.db.QueryRow(ctx, q, args...).Scan(&entity.UpdatedAt)
	}

	return err
}
