package user

const PasswordMinLength = 8
const PasswordPolicyMessage = "Password must be at least 8 characters."

type PasswordPolicyResponse struct {
	MinLength int    `json:"minLength"`
	Message   string `json:"message"`
}
