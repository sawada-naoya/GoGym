package gorm

import (
	"gogym-api/internal/adapter/db/gorm/record"
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
