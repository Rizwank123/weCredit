package repository

import (
	"context"
	"errors"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/weCredit/internal/domain"
)

type pgxLoginCodeRepository struct {
	db *pgxpool.Pool
}

func NewLoginCodeRepository(db *pgxpool.Pool) domain.LoginCodeRepository {
	return &pgxLoginCodeRepository{
		db: db,
	}
}

func (r pgxLoginCodeRepository) FindByID(ctx context.Context, id uuid.UUID) (result domain.LoginCode, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Retrieve the data
	q := `SELECT * FROM login_codes WHERE id = $1 AND deleted_at IS NULL LIMIT 1`
	args := []interface{}{id}

	// Retrieve the data
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
		return result, domain.DataNotFoundError{}
	}

	result, err = pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[domain.LoginCode])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return result, domain.DataNotFoundError{}
	}

	return result, err
}

func (r pgxLoginCodeRepository) FindByUsername(ctx context.Context, username string) (result domain.LoginCode, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Retrieve the data
	q := `SELECT * FROM login_codes WHERE username = $1 AND deleted_at IS NULL LIMIT 1`
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
		return result, domain.DataNotFoundError{}
	}

	result, err = pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[domain.LoginCode])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return result, domain.DataNotFoundError{}
	}

	return result, err
}

func (r pgxLoginCodeRepository) Create(ctx context.Context, entity *domain.LoginCode) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Create the data
	q := `INSERT INTO login_codes (username, code, expiry_time, status) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	args := []interface{}{entity.Username, entity.Code, entity.ExpiryTime, entity.Status}
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		err = tx.QueryRow(ctx, q, args...).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
	} else {
		err = r.db.QueryRow(ctx, q, args...).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
	}

	return err
}

func (r pgxLoginCodeRepository) Update(ctx context.Context, id uuid.UUID, entity *domain.LoginCode) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Update the data
	q := `UPDATE login_codes SET username=$1, code=$2, expiry_time=$3, status=$4, updated_at=NOW() WHERE id=$5 RETURNING updated_at`
	args := []interface{}{entity.Username, entity.Code, entity.ExpiryTime, entity.Status, id}
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		err = tx.QueryRow(ctx, q, args...).Scan(&entity.UpdatedAt)
	} else {
		err = r.db.QueryRow(ctx, q, args...).Scan(&entity.UpdatedAt)
	}

	return err
}

func (r pgxLoginCodeRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Delete the data
	q := `DELETE FROM login_codes WHERE id=$1`
	args := []interface{}{id}
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		_, err = tx.Exec(ctx, q, args...)
	} else {
		_, err = r.db.Exec(ctx, q, args...)
	}

	return err
}

func (r pgxLoginCodeRepository) DeleteByUsername(ctx context.Context, username string) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	// Delete the data
	q := `DELETE FROM login_codes WHERE username = $1`
	args := []interface{}{username}
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		_, err = tx.Exec(ctx, q, args...)
	} else {
		_, err = r.db.Exec(ctx, q, args...)
	}

	return err
}
