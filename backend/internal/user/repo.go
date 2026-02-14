package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"piguy.nl/koopify/internal/db"
)

type HashFn func(string) (string, error)

type UserRepository interface {
	CreateUser(ctx context.Context, user CreateUserRequest) (*UserResponse, error)
	LoginUser(ctx context.Context, user LoginUserRequest) (*UserResponse, error)
}

type PGUserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) PGUserRepository {
	return PGUserRepository{queries: queries}
}

func (pgur PGUserRepository) CreateUser(
	ctx context.Context,
	user CreateUserRequest,
) (*UserResponse, error) {
	params, err := user.ToDbParams(hashFn)
	if err != nil {
		return nil, err
	}

	dbUserResp, err := pgur.queries.CreateUser(ctx, *params)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == PG_UNIQUE_VIOLATION {
			return nil, ErrUserExists
		}

		return nil, err
	}

	userResp := UserResponseFrom(dbUserResp)
	return &userResp, nil
}

func (pgur PGUserRepository) LoginUser(ctx context.Context, user LoginUserRequest) (*UserResponse, error) {
	dbUser, err := pgur.queries.GetUserWithEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	userResp := UserResponseFrom(dbUser)
	return &userResp, nil
}

// Hash the password
// Hashing is related to how persistance is handled, so it is part of the repo
func hashFn(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
