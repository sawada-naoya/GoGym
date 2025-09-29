package mapper

import (
	"gogym-api/internal/adapter/repository/record"
	"gogym-api/internal/domain/user"

	"github.com/oklog/ulid/v2"
)

func FromUserEntity(u *user.User) *record.User {
	return &record.User{
		ID:           u.ID.String(),
		Email:        u.Email.String(),
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func ToUserEntity(r *record.User) *user.User {
	email, _ := user.NewEmail(r.Email)
	id, _ := ulid.Parse(r.ID)
	return &user.User{
		ID:           id,
		Email:        email,
		PasswordHash: r.PasswordHash,
		Name:         r.Name,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
