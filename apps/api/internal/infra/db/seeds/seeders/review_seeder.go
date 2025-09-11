package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/gorm"
)

// ReviewSeeder はレビューデータをシードする
type ReviewSeeder struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewReviewSeeder(db *gorm.DB, logger *slog.Logger) *ReviewSeeder {
	return &ReviewSeeder{
		db:     db,
		logger: logger,
	}
}

func (s *ReviewSeeder) Name() string {
	return "ReviewSeeder"
}

// ReviewSeedData はNDJSONから読み込むレビューデータ構造
type ReviewSeedData struct {
	UserEmail         string  `json:"user_email"`
	GymSlug           string  `json:"gym_slug"`
	Rating            int     `json:"rating"`
	Comment           string  `json:"comment,omitempty"`
	VisitDate         string  `json:"visit_date,omitempty"`       // YYYY-MM-DD形式
	VisitPurpose      string  `json:"visit_purpose,omitempty"`
	CleanlinessRating *int    `json:"cleanliness_rating,omitempty"`
	StaffRating       *int    `json:"staff_rating,omitempty"`
	EquipmentRating   *int    `json:"equipment_rating,omitempty"`
	ValueRating       *int    `json:"value_rating,omitempty"`
	IsVerifiedVisit   bool    `json:"is_verified_visit,omitempty"`
}

// Review はDBに挿入するレビュー構造体
type Review struct {
	ID                int64      `gorm:"primaryKey;autoIncrement"`
	UserID            int64      `gorm:"not null"`
	GymID             int64      `gorm:"not null"`
	Rating            int        `gorm:"not null"`
	Comment           *string
	HelpfulCount      int        `gorm:"default:0"`
	VisitDate         *time.Time
	VisitPurpose      *string
	CleanlinessRating *int
	StaffRating       *int
	EquipmentRating   *int
	ValueRating       *int
	IsVerifiedVisit   bool       `gorm:"default:false"`
}

// ユーザー、ジムのIDを取得するための構造体
type UserLookup struct {
	ID    int64  `gorm:"column:id"`
	Email string `gorm:"column:email"`
}

type GymLookup struct {
	ID   int64   `gorm:"column:id"`
	Slug *string `gorm:"column:slug"`
}

func (s *ReviewSeeder) Run(ctx context.Context) error {
	file, err := os.Open("internal/infra/db/seeds/data/reviews.ndjson")
	if err != nil {
		return fmt.Errorf("レビューシードファイルの読み込みに失敗: %w", err)
	}
	defer file.Close()

	// 既存データクリア
	if err := s.db.Exec("DELETE FROM reviews").Error; err != nil {
		return fmt.Errorf("レビューテーブルのクリアに失敗: %w", err)
	}

	// ユーザーとジムのマッピング作成
	userMap, err := s.createUserMap()
	if err != nil {
		return fmt.Errorf("ユーザーマッピング作成エラー: %w", err)
	}

	gymMap, err := s.createGymMap()
	if err != nil {
		return fmt.Errorf("ジムマッピング作成エラー: %w", err)
	}

	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {
		var seedData ReviewSeedData
		if err := json.Unmarshal(scanner.Bytes(), &seedData); err != nil {
			s.logger.Warn("レビューデータのパースエラー", "line", count+1, "error", err)
			continue
		}

		// ユーザーIDとジムIDの取得
		userID, exists := userMap[seedData.UserEmail]
		if !exists {
			s.logger.Warn("ユーザーが見つかりません", "email", seedData.UserEmail)
			continue
		}

		gymID, exists := gymMap[seedData.GymSlug]
		if !exists {
			s.logger.Warn("ジムが見つかりません", "slug", seedData.GymSlug)
			continue
		}

		review := Review{
			UserID:          userID,
			GymID:           gymID,
			Rating:          seedData.Rating,
			IsVerifiedVisit: seedData.IsVerifiedVisit,
		}

		// オプショナルフィールドの設定
		if seedData.Comment != "" {
			review.Comment = &seedData.Comment
		}
		
		if seedData.VisitDate != "" {
			visitDate, err := time.Parse("2006-01-02", seedData.VisitDate)
			if err == nil {
				review.VisitDate = &visitDate
			} else {
				s.logger.Warn("訪問日のパースエラー", "date", seedData.VisitDate, "error", err)
			}
		}

		if seedData.VisitPurpose != "" {
			review.VisitPurpose = &seedData.VisitPurpose
		}

		if seedData.CleanlinessRating != nil {
			review.CleanlinessRating = seedData.CleanlinessRating
		}
		if seedData.StaffRating != nil {
			review.StaffRating = seedData.StaffRating
		}
		if seedData.EquipmentRating != nil {
			review.EquipmentRating = seedData.EquipmentRating
		}
		if seedData.ValueRating != nil {
			review.ValueRating = seedData.ValueRating
		}

		// レビュー作成
		if err := s.db.Create(&review).Error; err != nil {
			s.logger.Error("レビュー作成エラー", 
				"user_email", seedData.UserEmail, 
				"gym_slug", seedData.GymSlug, 
				"error", err)
			continue
		}

		count++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ファイル読み込みエラー: %w", err)
	}

	s.logger.Info("レビューシード完了", "作成数", count)
	return nil
}

// createUserMap はメールアドレスからユーザーIDへのマッピングを作成
func (s *ReviewSeeder) createUserMap() (map[string]int64, error) {
	var users []UserLookup
	if err := s.db.Table("users").Select("id, email").Find(&users).Error; err != nil {
		return nil, err
	}

	userMap := make(map[string]int64)
	for _, user := range users {
		userMap[user.Email] = user.ID
	}
	return userMap, nil
}

// createGymMap はスラッグからジムIDへのマッピングを作成
func (s *ReviewSeeder) createGymMap() (map[string]int64, error) {
	var gyms []GymLookup
	if err := s.db.Table("gyms").Select("id, slug").Find(&gyms).Error; err != nil {
		return nil, err
	}

	gymMap := make(map[string]int64)
	for _, gym := range gyms {
		if gym.Slug != nil {
			gymMap[*gym.Slug] = gym.ID
		}
	}
	return gymMap, nil
}