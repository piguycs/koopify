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

	// skip updating non password fields if they have not been changed
	// I wanted to make the password update field be part of the same PATCH interface as other fields
	// therefore this little hack. As I do not wish for password updating to be part of the patch functions
	// in the backend. Therefore the two `UpdateUser` and `UpdateUserPassword` fns
	nonPasswdChanged := false

	newDisplayName := currentUser.DisplayName
	if update.DisplayName != nil {
		nonPasswdChanged = true
		newDisplayName = *update.DisplayName
	}

	newEmail := currentUser.Email
	if update.Email != nil {
		nonPasswdChanged = true
		newEmail = *update.Email
	}

	// this is so hacky that I forgot how to declare typed variables in Go lmao
	// for a moment muscle memory kicked in and I did `let user: UserResponse`
	var user *UserResponse = nil
	if update.Password != nil {
		user, err = us.repository.UpdateUserPassword(ctx, userID, *update.Password)
		if err != nil {
			return nil, err
		}
	}

	if nonPasswdChanged {
		return us.repository.UpdateUser(ctx, userID, newDisplayName, newEmail)
	} else {
		// just return the password's return value here
		return user, nil
	}

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

func (us *UserService) CancelUserDeletionAdmin(ctx context.Context, userID int64) (*UserResponse, error) {
	return us.repository.CancelUserDeletion(ctx, userID)
}

func (us *UserService) UpdateUserDetailsAdmin(ctx context.Context, userID int64, update UpdateUserRequest) (*UserResponse, error) {
	// First, get the target user to check if they are an admin
	targetUser, err := us.repository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Cannot edit admin users
	if targetUser.Admin {
		return nil, ErrCannotEditAdmin
	}

	// Apply updates
	newDisplayName := targetUser.DisplayName
	if update.DisplayName != nil {
		newDisplayName = *update.DisplayName
	}

	newEmail := targetUser.Email
	if update.Email != nil {
		newEmail = *update.Email
	}

	return us.repository.UpdateUser(ctx, userID, newDisplayName, newEmail)
}
