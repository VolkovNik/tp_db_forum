package domain

type User struct {
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	About    string `json:"about,omitempty"`
}

type Users []User

type UserRepository interface {
	Create(user User) (*User, error)
	SelectByEmailOrNickname(nickname string, email string) (Users, error)
	GetProfileInfo(nickname string) (User, error)
	UpdateProfileInfo(user *User) error
}

type UserUsecase interface {
	Create(user User) (*User, error)
	SelectByEmailOrNickname(nickname string, email string) (Users, error)
	GetProfileInfo(nickname string) (User, error)
	UpdateProfileInfo(user *User) error
}