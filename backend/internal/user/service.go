package user

import "context"

type UserService struct {
	repository UserRepository
}

// Create a new user service
// Requires the user repository for persisting data
func NewUserService(repository UserRepository) UserService {
	return UserService{repository: repository}
}

// Create a new user on the database
func (us *UserService) CreateUser(ctx context.Context, user CreateUserRequest) (*UserResponse, error) {
	return us.repository.CreateUser(ctx, user)
}

func (us *UserService) LoginUser(ctx context.Context, user LoginUserRequest) (*UserResponse, error) {
	return us.repository.LoginUser(ctx, user)
}

func (us *UserService) GetUserByID(ctx context.Context, id int64) (*UserResponse, error) {
	return us.repository.GetUserByID(ctx, id)
}

func (us *UserService) ListUsers(ctx context.Context) ([]UserResponse, error) {
	return us.repository.ListUsers(ctx)
}

func (us *UserService) UpdateUserAdmin(ctx context.Context, userID int64, admin bool) (*UserResponse, error) {
	return us.repository.UpdateUserAdmin(ctx, userID, admin)
}

func (us *UserService) UpdateCurrentUser(
	ctx context.Context,
	userID int64,
	update UpdateUserRequest,
) (*UserResponse, error) {
	currentUser, err := us.repository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	newDisplayName := currentUser.DisplayName
	if update.DisplayName != nil {
		newDisplayName = *update.DisplayName
	}

	newEmail := currentUser.Email
	if update.Email != nil {
		newEmail = *update.Email
	}

	return us.repository.UpdateUser(ctx, userID, newDisplayName, newEmail)
}

func (us *UserService) RequestDeletion(ctx context.Context, userID int64) (*UserResponse, error) {
	delayHours, err := us.repository.GetDeletionPolicy(ctx)
	if err != nil {
		return nil, err
	}

	return us.repository.RequestUserDeletion(ctx, userID, delayHours)
}

func (us *UserService) CancelDeletion(ctx context.Context, userID int64) (*UserResponse, error) {
	return us.repository.CancelUserDeletion(ctx, userID)
}
