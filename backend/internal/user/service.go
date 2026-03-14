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
