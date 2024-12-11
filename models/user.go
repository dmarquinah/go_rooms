package models

type User struct {
	UserId   string
	UserRole string
}

func NewUserInstance(userId string, role string) *User {
	return &User{
		UserId:   userId,
		UserRole: role,
	}
}
