package domain

type GetRoleResponse struct {
	RoleName string
}

type User struct {
	Id       int64
	Fio      string
	Username string
	RoleId   int32
	RoleName string
}

type UserProfile struct {
	Id       int64
	Fio      string
	Username string
	RoleId   int32
	RoleName string
}

type EditUserRequest struct {
	UserId   int64  `validate:"required"`
	Username string `validate:"min=4,max=50"`
	Fio      string `validate:"min=6,max=60"`
	Password string `validate:"min=6,max=20"`
	RoleId   int32  `validate:"required"`
}

type DeleteUserRequest struct {
	UserId int32 `validate:"required"`
}
