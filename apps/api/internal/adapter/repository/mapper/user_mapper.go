package mapper

import (
	"gogym-api/internal/adapter/repository/record"
	"gogym-api/internal/domain/user"
)

func FromUserEntity(u *user.User) *record.User {
	return &record.User{
		ID:           string(u.ID),
		Email:        u.Email.String(),
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func ToUserEntity(r *record.User) *user.User {
	email, _ := user.NewEmail(r.Email)
	return &user.User{
		ID:           r.ID,
		Email:        email,
		PasswordHash: r.PasswordHash,
		Name:         r.Name,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
