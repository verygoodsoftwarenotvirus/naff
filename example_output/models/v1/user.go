package models

import (
	"context"
	"net/http"
)

const (
	UserKey        ContextKey = "user"
	UserIDKey      ContextKey = "user_id"
	UserIsAdminKey ContextKey = "is_admin"
)

type (
	User struct {
		ID                    uint64
		Username              string
		HashedPassword        string
		Salt                  []byte
		TwoFactorSecret       string
		PasswordLastChangedOn *uint64
		IsAdmin               bool
		CreatedOn             uint64
		UpdatedOn             *uint64
		ArchivedOn            *uint64
	}
	UserList struct {
		Pagination
		Users []User
	}
	UserLoginInput struct {
		Username  string
		Password  string
		TOTPToken string
	}
	UserInput struct {
		Username        string
		Password        string
		TwoFactorSecret string
	}
	UserCreationResponse struct {
		ID                    uint64
		Username              string
		TwoFactorSecret       string
		PasswordLastChangedOn *uint64
		IsAdmin               bool
		CreatedOn             uint64
		UpdatedOn             *uint64
		ArchivedOn            *uint64
		TwoFactorQRCode       string
	}
	PasswordUpdateInput struct {
		NewPassword     string
		CurrentPassword string
		TOTPToken       string
	}
	TOTPSecretRefreshInput struct {
		CurrentPassword string
		TOTPToken       string
	}
	TOTPSecretRefreshResponse struct {
		TwoFactorSecret string
	}
	UserDataManager interface {
		GetUser(ctx context.Context, userID uint64) (*User, error)
		GetUserByUsername(ctx context.Context, username string) (*User, error)
		GetUserCount(ctx context.Context, filter *QueryFilter) (uint64, error)
		GetUsers(ctx context.Context, filter *QueryFilter) (*UserList, error)
		CreateUser(ctx context.Context, input *UserInput) (*User, error)
		UpdateUser(ctx context.Context, updated *User) error
		ArchiveUser(ctx context.Context, userID uint64) error
	}
	UserDataServer interface {
		UserInputMiddleware(next http.Handler) http.Handler
		PasswordUpdateInputMiddleware(next http.Handler) http.Handler
		TOTPSecretRefreshInputMiddleware(next http.Handler) http.Handler
		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		NewTOTPSecretHandler() http.HandlerFunc
		UpdatePasswordHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update accepts a User as input and merges those values if they're set
func (u *User) Update(input *User) {
	if input.Username != "" && input.Username != u.Username {
		u.Username = input.Username
	}
	if input.HashedPassword != "" && input.HashedPassword != u.HashedPassword {
		u.HashedPassword = input.HashedPassword
	}
	if input.TwoFactorSecret != "" && input.TwoFactorSecret != u.TwoFactorSecret {
		u.TwoFactorSecret = input.TwoFactorSecret
	}
}
