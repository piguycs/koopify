package user

import (
	"context"
	"errors"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"piguy.nl/koopify/internal/db"
)

type HashFn func(string) (string, error)

type UserRepository interface {
	CreateUser(ctx context.Context, user CreateUserRequest) (*UserResponse, error)
	LoginUser(ctx context.Context, user LoginUserRequest) (*UserResponse, error)
	GetUserByID(ctx context.Context, id int64) (*UserResponse, error)
	ListUsers(ctx context.Context) ([]UserResponse, error)
	UpdateUser(ctx context.Context, id int64, displayName string, email string) (*UserResponse, error)
	UpdateUserPassword(ctx context.Context, id int64, password string) (*UserResponse, error)
	UpdateUserAdmin(ctx context.Context, id int64, admin bool) (*UserResponse, error)
	GetDeletionPolicy(ctx context.Context) (int32, error)
	RequestUserDeletion(ctx context.Context, id int64, delayHours int32) (*UserResponse, error)
	CancelUserDeletion(ctx context.Context, id int64) (*UserResponse, error)
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if dbUser.DeletionScheduledAt.Valid {
		now := time.Now()
		if !dbUser.DeletionScheduledAt.Time.After(now) {
			return nil, ErrAccountDeletionScheduled
		}
	}

	userResp := UserResponseFrom(dbUser)
	return &userResp, nil
}

func (pgur PGUserRepository) GetUserByID(ctx context.Context, id int64) (*UserResponse, error) {
	dbUser, err := pgur.queries.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	userResp := UserResponseFrom(dbUser)
	return &userResp, nil
}

func (pgur PGUserRepository) ListUsers(ctx context.Context) ([]UserResponse, error) {
	dbUsers, err := pgur.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]UserResponse, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = UserResponseFrom(dbUser)
	}
	return users, nil
}

func (pgur PGUserRepository) UpdateUser(
	ctx context.Context,
	id int64,
	displayName string,
	email string,
) (*UserResponse, error) {
	params := db.UpdateUserParams{
		ID:          id,
		DisplayName: displayName,
		Email:       email,
	}

	dbUser, err := pgur.queries.UpdateUser(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == PG_UNIQUE_VIOLATION {
			return nil, ErrUserExists
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	userResp := UserResponseFrom(dbUser)
	return &userResp, nil
}

func (pgur PGUserRepository) UpdateUserPassword(ctx context.Context, id int64, password string) (*UserResponse, error) {
	hashedPassword, err := hashFn(password)
	if err != nil {
		log.Error("Unable to run hashFn on password due to", "error", err)
		return nil, err
	}

	params := db.UpdateUserPasswordParams{
		ID:       id,
		Password: hashedPassword,
	}

	dbUser, err := pgur.queries.UpdateUserPassword(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == PG_UNIQUE_VIOLATION {
			return nil, ErrUserExists
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	userResp := UserResponseFrom(dbUser)
	return &userResp, nil
}

func (pgur PGUserRepository) UpdateUserAdmin(ctx context.Context, id int64, admin bool) (*UserResponse, error) {
	dbUser, err := pgur.queries.UpdateUserAdmin(ctx, db.UpdateUserAdminParams{
		ID:    id,
		Admin: admin,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	userResp := UserResponseFrom(dbUser)
	return &userResp, nil
}

func (pgur PGUserRepository) GetDeletionPolicy(ctx context.Context) (int32, error) {
	policy, err := pgur.queries.GetDeletionPolicy(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrDeletionPolicyNotFound
		}
		return 0, err
	}
	return policy.DeletionDelayHours, nil
}

func (pgur PGUserRepository) RequestUserDeletion(
	ctx context.Context,
	id int64,
	delayHours int32,
) (*UserResponse, error) {
	dbUser, err := pgur.queries.RequestUserDeletion(ctx, db.RequestUserDeletionParams{
		ID:    id,
		Hours: delayHours,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	userResp := UserResponseFrom(dbUser)
	return &userResp, nil
}

func (pgur PGUserRepository) CancelUserDeletion(
	ctx context.Context,
	id int64,
) (*UserResponse, error) {
	dbUser, err := pgur.queries.CancelUserDeletion(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
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
