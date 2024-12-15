package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid/v5"
)

type (
	// JSONB represents a JSONB type
	JSONB map[string]interface{} // @name JSONB
)

type (
	// Address Defines the model for address
	Address struct {
		Location string `json:"location" example:"Ahmedabad"`
		Street   string `json:"street" example:"Near Railway Station"`
		City     string `json:"city" example:"Ahmedabad"`
		State    string `json:"state" example:"Gujarat"`
		Country  string `json:"country" example:"India"`
		Pincode  string `json:"pincode" example:"380009"`
	} // @name Address

	// Base define the base model
	Base struct {
		ID uuid.UUID `json:"id" db:"id" example:""`
	} // @name Base
	// BaseAudit define the base audit model
	BaseAudit struct {
		CreatedAt time.Time `json:"created_at" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	} // @name BaseAudit
)
type (
	// BaseResponse defines base response fields.
	BaseResponse struct {
		Data interface{} `json:"data"`
	} // @name BaseResponse

	// PaginationResponse defines pagination response fields.
	PaginationResponse struct {
		Data  interface{} `json:"data"`
		Total int64       `json:"total" example:"1000"`
		Size  int64       `json:"size" example:"10"`
		Page  int64       `json:"page" example:"1"`
	} // @name PaginationResponse
)

type (
	// Transactioner defines the methods that any transactioner should implement.
	Transactioner interface {
		Begin(ctx context.Context) (result context.Context, err error)
		Commit(ctx context.Context) (err error)
		Rollback(ctx context.Context, err error)
	}
)

// Value implements the driver.Valuer interface,
func (j *JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

// Scan implements the sql.Scanner interface,
func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal([]byte(value.(string)), &j); err != nil {
		return err
	}
	return nil
}
