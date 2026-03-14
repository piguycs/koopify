package user

import "piguy.nl/koopify/internal/db"

type CreateUserRequest struct {
	DisplayName string `json:"displayName" validate:"required,lte=64,gte=4"`
	Email       string `json:"email" validate:"required,lte=128"`
	Password    string `json:"password" validate:"required,lte=72,gte=4"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,lte=128"`
	Password string `json:"password" validate:"required,lte=72,gte=4"`
}

type UserResponse struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Admin       bool   `json:"admin"`
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
		ID:          dbu.ID,
		DisplayName: dbu.DisplayName,
		Email:       dbu.Email,
		Admin:       dbu.Admin,
	}
}
