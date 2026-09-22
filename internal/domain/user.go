package domain

type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

func NewUser(
	id int,
	email string,
	passwordHash string,
) *User {
	return &User{
		ID:           int64(id),
		Email:        email,
		PasswordHash: passwordHash,
	}
}
