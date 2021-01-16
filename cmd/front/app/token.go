package app

type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`
}

type Payload struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
	Role  string `json:"role"`
	Exp   int64  `json:"exp"`
}
