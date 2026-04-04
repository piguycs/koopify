package user

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"piguy.nl/koopify/internal/db"
)

type CreateUserRequest struct {
	DisplayName string `json:"displayName" validate:"required,lte=64,gte=4"`
	Email       string `json:"email" validate:"required,lte=128"`
	Password    string `json:"password" validate:"required,lte=72"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,lte=128"`
	Password string `json:"password" validate:"required,lte=72"`
}

type UpdateUserRequest struct {
	DisplayName *string `json:"displayName" validate:"omitempty,lte=64,gte=4"`
	Email       *string `json:"email" validate:"omitempty,lte=128,email"`
	Password    *string `json:"password" validate:"omitempty,lte=72,gte=8"`
}

type UserResponse struct {
	ID                  int64      `json:"id"`
	DisplayName         string     `json:"displayName"`
	Email               string     `json:"email"`
	Admin               bool       `json:"admin"`
	RequestedDeletionAt *time.Time `json:"requestedDeletionAt"`
	DeletionScheduledAt *time.Time `json:"deletionScheduledAt"`
}

func (cur CreateUserRequest) ToDbParams(hashFn HashFn) (*db.CreateUserParams, error) {
	password, err := hashFn(cur.Password)
	if err != nil {
		return nil, err
	}

	params := db.CreateUserParams{
		DisplayName: cur.DisplayName,
		Email:       cur.Email,
		Password:    password,
	}

	return &params, nil
}

func UserResponseFrom(dbu db.User) UserResponse {
	return UserResponse{
		ID:                  dbu.ID,
		DisplayName:         dbu.DisplayName,
		Email:               dbu.Email,
		Admin:               dbu.Admin,
		RequestedDeletionAt: timeFromTimestamptz(dbu.RequestedDeletionAt),
		DeletionScheduledAt: timeFromTimestamptz(dbu.DeletionScheduledAt),
	}
}

func timeFromTimestamptz(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	value := ts.Time
	return &value
}
