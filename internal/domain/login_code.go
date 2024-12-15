package domain

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

// LoginCodeStatus defines model for LoginCode.Status.
type LoginCodeStatus string

type (
	// LoginCode defines model for LoginCode.
	LoginCode struct {
		Base
		Username     string          `db:"username" json:"-"`
		Code         string          `db:"code" json:"-"`
		ExpiryTime   time.Time       `db:"expiry_time" json:"-"`
		Status       LoginCodeStatus `db:"status" json:"-"`
		ResponseMeta *string         `db:"response_meta" json:"-"`
		BaseAudit
		DeletedAt *time.Time `db:"deleted_at" json:"-"`
	} // @name LoginCode
	// OtpMessage defines model for OtpMessage.
	OtpMessage struct {
		To  string `json:"-"`
		Otp string `json:"-"`
	} // @name OtpMessage
)

type (
	// LoginCodeRepository defines the methods that any login-code repository should implement.
	LoginCodeRepository interface {
		// FindByID returns a record by id
		FindByID(ctx context.Context, id uuid.UUID) (result LoginCode, err error)
		// FindByUsername returns a record by username
		FindByUsername(ctx context.Context, username string) (result LoginCode, err error)
		// Create creates a new record
		Create(ctx context.Context, entity *LoginCode) (err error)
		// Update updates an existing record
		Update(ctx context.Context, id uuid.UUID, entity *LoginCode) (err error)
		// Delete deletes an existing record by id
		Delete(ctx context.Context, id uuid.UUID) (err error)
		// DeleteByUsername deletes login codes by username.
		DeleteByUsername(ctx context.Context, username string) (err error)
	}
)

const (
	LoginCodeStatusPENDING LoginCodeStatus = "PENDING"
	LoginCodeStatusSUCCESS LoginCodeStatus = "SUCCESS"
	LoginCodeStatusFAILED  LoginCodeStatus = "FAILED"
)
