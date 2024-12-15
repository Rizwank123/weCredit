package domain

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type (
	// UserRole represents the role a user in the system.
	UserRole string
)
type (
	// User defines the module for User
	User struct {
		Base
		UserName string `db:"user_name" json:"user_name,omitempty" example:"+919876543210"`
		Role     string `db:"role" json:"role,omitempty"  example:"USER"`
		FullName string `db:"full_name" json:"full_name,omitempty" example:"John Doe"`
		BaseAudit
	} // @name User

)

type (
	// CreateUserInput define the module for CreateUser
	RegisterUserInput struct {
		FullName string   `json:"full_name" example:"John Doe"`
		UserName string   `json:"user_name" example:"+919876543210"`
		Role     UserRole `json:"role" example:"RECEPTIONIST"`
	} // @name CreateUserInput
	// UpdateUserInput define the module for the UpdateUserInput
	UpdateUserInput struct {
		RegisterUserInput
	} // @name UpdateUserInput
	// InitLoginInput define the module for the InitLoginInput
	InitLoginInput struct {
		UserName string `json:"username" example:"+919876543210"`
	} // @name InitLoginInput
	// LoginInput  define the module for the LoginInput
	LoginInput struct {
		UserName string `json:"username" example:"+919876543210"`
		Otp      string `json:"otp" example:"123456"`
	} // @name LoginInput
	// LoginOutput define the module for the LoginOutput
	LoginOutput struct {
		Token     string `json:"token"`
		ExpiresIn int64  `json:"expires_in"`
	} // @name LoginOutput
)

type (
	// UserRepository defines the methods that any use repository should implements
	UserRepository interface {
		// FindByID return the user by id
		FindByID(ctx context.Context, id uuid.UUID) (result User, err error)
		// FindByUserName return the user by username
		FindByUserName(ctx context.Context, username string) (result User, err error)
		// CreateUser creates a new user
		CreateUser(ctx context.Context, entity *User) (err error)
		// UpdateUser updates the user
		UpdateUser(ctx context.Context, entity *User) (err error)
		// DeleteUser deletes the user
		DeleteUser(ctx context.Context, id uuid.UUID) (err error)
	}

	// UserService defines the methods that any use service should implements
	UserService interface {
		// Login login the user
		Login(input LoginInput) (result LoginOutput, err error)
		// InitLogin init the login
		InitLogin(input InitLoginInput) (err error)
		// RegisterUser register a new user
		RegisterUser(input RegisterUserInput) (result User, err error)
		// FindByUserName find the user by username
		FindByUserName(username string) (result User, err error)
		// FindByID find the user by id
		FindByID(id uuid.UUID) (result User, err error)
	}
)
