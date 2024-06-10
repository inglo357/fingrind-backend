package requests

type UserLoginRequest struct{
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUsernameRequest struct{
	NewUsername string `json:"new_username" binding:"required"`
}

type UpdatePasswordRequest struct{
	NewPassword string `json:"new_password" binding:"required"`
}