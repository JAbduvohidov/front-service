package models

type UserRequestDTO struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	Avatar   string `json:"avatar"`
}

type UserResponseDTO struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

type TokenRequestDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TokenResponseDTO struct {
	Token string `json:"token"`
}

type ErrorDTO struct {
	Errors []string `json:"errors"`
}