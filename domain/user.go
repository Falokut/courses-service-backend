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

type DeleteUserRequest struct {
	UserId int32 `validate:"required"`
}
