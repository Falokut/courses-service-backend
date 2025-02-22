package entity

type User struct {
	Id       int64
	Fio      string
	Username string
	Password string
	RoleId   int32
	RoleName string
}

type RegisterUser struct {
	Fio          string
	Username     string
	PasswordHash string
	RoleId       int32
}
