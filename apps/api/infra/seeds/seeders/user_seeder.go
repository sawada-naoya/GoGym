package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSeeder はユーザーデータをシードする
type UserSeeder struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewUserSeeder はUserSeederを作成する
func NewUserSeeder(db *gorm.DB, logger *slog.Logger) *UserSeeder {
	return &UserSeeder{
		db:     db,
		logger: logger,
	}
}

func (s *UserSeeder) Name() string {
	return "UserSeeder"
}

// UserSeedData はNDJSONから読み込むユーザーデータ構造
type UserSeedData struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio,omitempty"`
	Location    string `json:"location,omitempty"`
	BirthYear   *int   `json:"birth_year,omitempty"`
	Gender      string `json:"gender,omitempty"`
	IsVerified  bool   `json:"is_verified,omitempty"`
}

// User はDBに挿入するユーザー構造体
type User struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	DisplayName  string `gorm:"not null"`
	Bio          *string
	Location     *string
	BirthYear    *int
	Gender       *string
	IsVerified   bool `gorm:"default:false"`
}

func (s *UserSeeder) Run(ctx context.Context) error {
	// NDJSONファイルを開く
	file, err := os.Open("infra/seeds/data/users.ndjson")
	if err != nil {
		return fmt.Errorf("ユーザーシードファイルの読み込みに失敗: %w", err)
	}
	defer file.Close()

	// 既存データをクリア（開発用）
	if err := s.db.Exec("DELETE FROM users").Error; err != nil {
		return fmt.Errorf("ユーザーテーブルのクリアに失敗: %w", err)
	}
	
	scanner := bufio.NewScanner(file)
	count := 0
	
	for scanner.Scan() {
		var seedData UserSeedData
		if err := json.Unmarshal(scanner.Bytes(), &seedData); err != nil {
			s.logger.Warn("ユーザーデータのパースエラー", "line", count+1, "error", err)
			continue
		}

		// デフォルトパスワードをハッシュ化
		passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("パスワードハッシュ化エラー: %w", err)
		}

		user := User{
			Email:        seedData.Email,
			PasswordHash: string(passwordHash),
			DisplayName:  seedData.DisplayName,
			IsVerified:   seedData.IsVerified,
		}

		// オプショナルフィールドの設定
		if seedData.Bio != "" {
			user.Bio = &seedData.Bio
		}
		if seedData.Location != "" {
			user.Location = &seedData.Location
		}
		if seedData.BirthYear != nil {
			user.BirthYear = seedData.BirthYear
		}
		if seedData.Gender != "" {
			user.Gender = &seedData.Gender
		}

		// データベースに挿入
		if err := s.db.Create(&user).Error; err != nil {
			s.logger.Error("ユーザー作成エラー", "email", seedData.Email, "error", err)
			continue
		}

		count++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ファイル読み込みエラー: %w", err)
	}

	s.logger.Info("ユーザーシード完了", "作成数", count)
	return nil
}