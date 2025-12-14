package gym

import (
	gymUsecase "gogym-api/internal/usecase/gym"

	"gorm.io/gorm"
)

type gymRepository struct {
	db *gorm.DB
}

func NewGymRepository(db *gorm.DB) gymUsecase.Repository {
	return &gymRepository{db: db}
}
