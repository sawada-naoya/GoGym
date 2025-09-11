package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"gorm.io/gorm"
)

// GymSeeder はジムデータをシードする
type GymSeeder struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewGymSeeder(db *gorm.DB, logger *slog.Logger) *GymSeeder {
	return &GymSeeder{
		db:     db,
		logger: logger,
	}
}

func (s *GymSeeder) Name() string {
	return "GymSeeder"
}

// 位置情報構造体
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// 営業時間構造体
type GymHour struct {
	DayOfWeek int     `json:"day_of_week"`
	IsClosed  bool    `json:"is_closed"`
	OpenTime  *string `json:"open_time,omitempty"`
	CloseTime *string `json:"close_time,omitempty"`
}

// 料金プラン構造体
type PricingPlan struct {
	PlanName     string `json:"plan_name"`
	PlanType     string `json:"plan_type"`
	Price        int    `json:"price"`
	Description  string `json:"description,omitempty"`
	IsPopular    bool   `json:"is_popular,omitempty"`
	DisplayOrder int    `json:"display_order,omitempty"`
}

// GymSeedData はNDJSONから読み込むジムデータ構造
type GymSeedData struct {
	Name           string        `json:"name"`
	Description    string        `json:"description,omitempty"`
	Location       Location      `json:"location"`
	Address        string        `json:"address"`
	City           string        `json:"city,omitempty"`
	Prefecture     string        `json:"prefecture,omitempty"`
	PostalCode     string        `json:"postal_code,omitempty"`
	PhoneNumber    string        `json:"phone_number,omitempty"`
	Website        string        `json:"website,omitempty"`
	AccessInfo     string        `json:"access_info,omitempty"`
	ParkingInfo    string        `json:"parking_info,omitempty"`
	Amenities      []string      `json:"amenities,omitempty"`
	Capacity       *int          `json:"capacity,omitempty"`
	OperatorName   string        `json:"operator_name,omitempty"`
	BrandName      string        `json:"brand_name,omitempty"`
	PriceRangeMin  *int          `json:"price_range_min,omitempty"`
	PriceRangeMax  *int          `json:"price_range_max,omitempty"`
	Slug           string        `json:"slug"`
	MetaTitle      string        `json:"meta_title,omitempty"`
	MetaDescription string       `json:"meta_description,omitempty"`
	MainPhotoURL   string        `json:"main_photo_url,omitempty"`
	Photos         []string      `json:"photos,omitempty"`
	Hours          []GymHour     `json:"hours,omitempty"`
	PricingPlans   []PricingPlan `json:"pricing_plans,omitempty"`
}

// データベースエンティティ
type Gym struct {
	ID              int64   `gorm:"primaryKey;autoIncrement"`
	Name            string  `gorm:"not null"`
	Description     *string
	Location        string  `gorm:"type:POINT"`  // WKT形式で格納
	Address         string  `gorm:"not null"`
	City            *string
	Prefecture      *string
	PostalCode      *string
	PhoneNumber     *string
	Website         *string
	AccessInfo      *string
	ParkingInfo     *string
	Amenities       *string `gorm:"type:JSON"`
	Capacity        *int
	OperatorName    *string
	BrandName       *string
	PriceRangeMin   *int
	PriceRangeMax   *int
	AverageRating   float64 `gorm:"default:0.00"`
	ReviewCount     int     `gorm:"default:0"`
	IsActive        bool    `gorm:"default:true"`
	Slug            *string `gorm:"unique"`
	MetaTitle       *string
	MetaDescription *string
	MainPhotoURL    *string
	Photos          *string `gorm:"type:JSON"`
}

type GymHourDB struct {
	ID        int64   `gorm:"primaryKey;autoIncrement"`
	GymID     int64   `gorm:"not null"`
	DayOfWeek int     `gorm:"not null"`
	OpenTime  *string `gorm:"type:TIME"`
	CloseTime *string `gorm:"type:TIME"`
	IsClosed  bool    `gorm:"default:false"`
}

type GymPricingPlan struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	GymID        int64  `gorm:"not null"`
	PlanName     string `gorm:"not null"`
	PlanType     string `gorm:"not null"`
	Price        int    `gorm:"not null"`
	Description  *string
	IsPopular    bool `gorm:"default:false"`
	DisplayOrder int  `gorm:"default:0"`
	IsActive     bool `gorm:"default:true"`
}

func (s *GymSeeder) Run(ctx context.Context) error {
	file, err := os.Open("internal/infra/db/seeds/data/gyms.ndjson")
	if err != nil {
		return fmt.Errorf("ジムシードファイルの読み込みに失敗: %w", err)
	}
	defer file.Close()

	// 既存データクリア
	s.db.Exec("DELETE FROM gym_pricing_plans")
	s.db.Exec("DELETE FROM gym_hours") 
	s.db.Exec("DELETE FROM gyms")

	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {
		var seedData GymSeedData
		if err := json.Unmarshal(scanner.Bytes(), &seedData); err != nil {
			s.logger.Warn("ジムデータのパースエラー", "line", count+1, "error", err)
			continue
		}

		// トランザクション開始
		tx := s.db.Begin()

		// ジム基本データ作成
		gym := Gym{
			Name:    seedData.Name,
			Address: seedData.Address,
			IsActive:  true,
		}

		// オプショナルフィールド設定
		if seedData.Description != "" {
			gym.Description = &seedData.Description
		}
		if seedData.City != "" {
			gym.City = &seedData.City
		}
		if seedData.Prefecture != "" {
			gym.Prefecture = &seedData.Prefecture
		}
		if seedData.PostalCode != "" {
			gym.PostalCode = &seedData.PostalCode
		}
		if seedData.PhoneNumber != "" {
			gym.PhoneNumber = &seedData.PhoneNumber
		}
		if seedData.Website != "" {
			gym.Website = &seedData.Website
		}
		if seedData.AccessInfo != "" {
			gym.AccessInfo = &seedData.AccessInfo
		}
		if seedData.ParkingInfo != "" {
			gym.ParkingInfo = &seedData.ParkingInfo
		}
		if seedData.OperatorName != "" {
			gym.OperatorName = &seedData.OperatorName
		}
		if seedData.BrandName != "" {
			gym.BrandName = &seedData.BrandName
		}
		if seedData.PriceRangeMin != nil {
			gym.PriceRangeMin = seedData.PriceRangeMin
		}
		if seedData.PriceRangeMax != nil {
			gym.PriceRangeMax = seedData.PriceRangeMax
		}
		if seedData.Slug != "" {
			gym.Slug = &seedData.Slug
		}
		if seedData.MetaTitle != "" {
			gym.MetaTitle = &seedData.MetaTitle
		}
		if seedData.MetaDescription != "" {
			gym.MetaDescription = &seedData.MetaDescription
		}
		if seedData.MainPhotoURL != "" {
			gym.MainPhotoURL = &seedData.MainPhotoURL
		}
		if seedData.Capacity != nil {
			gym.Capacity = seedData.Capacity
		}

		// JSON配列をJSONストリングに変換
		if len(seedData.Amenities) > 0 {
			amenitiesJSON, _ := json.Marshal(seedData.Amenities)
			amenitiesStr := string(amenitiesJSON)
			gym.Amenities = &amenitiesStr
		}
		if len(seedData.Photos) > 0 {
			photosJSON, _ := json.Marshal(seedData.Photos)
			photosStr := string(photosJSON)
			gym.Photos = &photosStr
		}

		// Raw SQLでジムデータを直接挿入（POINT型を含む、SRID 4326指定）
		insertSQL := `INSERT INTO gyms (name, description, location, address, city, prefecture, postal_code, phone_number, website, access_info, parking_info, amenities, capacity, operator_name, brand_name, price_range_min, price_range_max, is_active, slug, meta_title, meta_description, main_photo_url, photos) VALUES (?, ?, ST_GeomFromText(?, 4326), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		
		pointWKT := fmt.Sprintf("POINT(%f %f)", seedData.Location.Lat, seedData.Location.Lng)
		
		result := tx.Exec(insertSQL,
			gym.Name, gym.Description, pointWKT, gym.Address, gym.City, gym.Prefecture,
			gym.PostalCode, gym.PhoneNumber, gym.Website, gym.AccessInfo, gym.ParkingInfo,
			gym.Amenities, gym.Capacity, gym.OperatorName, gym.BrandName, 
			gym.PriceRangeMin, gym.PriceRangeMax, gym.IsActive, gym.Slug,
			gym.MetaTitle, gym.MetaDescription, gym.MainPhotoURL, gym.Photos)
			
		if result.Error != nil {
			tx.Rollback()
			s.logger.Error("ジム作成エラー", "name", seedData.Name, "error", result.Error)
			continue
		}
		
		// 挿入されたジムのIDを取得
		var insertedGym struct {
			ID int64 `gorm:"column:id"`
		}
		if err := tx.Raw("SELECT id FROM gyms WHERE slug = ? OR (name = ? AND address = ?)", 
			gym.Slug, gym.Name, gym.Address).Scan(&insertedGym).Error; err != nil {
			tx.Rollback()
			s.logger.Error("挿入されたジムID取得エラー", "name", seedData.Name, "error", err)
			continue
		}
		gym.ID = insertedGym.ID

		// 営業時間データ挿入
		for _, hour := range seedData.Hours {
			gymHour := GymHourDB{
				GymID:     gym.ID,
				DayOfWeek: hour.DayOfWeek,
				IsClosed:  hour.IsClosed,
			}
			if hour.OpenTime != nil {
				gymHour.OpenTime = hour.OpenTime
			}
			if hour.CloseTime != nil {
				gymHour.CloseTime = hour.CloseTime
			}
			
			if err := tx.Table("gym_hours").Create(&gymHour).Error; err != nil {
				tx.Rollback()
				s.logger.Error("営業時間作成エラー", "gym", seedData.Name, "error", err)
				break
			}
		}

		// 料金プランデータ挿入
		for _, plan := range seedData.PricingPlans {
			pricingPlan := GymPricingPlan{
				GymID:        gym.ID,
				PlanName:     plan.PlanName,
				PlanType:     plan.PlanType,
				Price:        plan.Price,
				IsPopular:    plan.IsPopular,
				DisplayOrder: plan.DisplayOrder,
				IsActive:     true,
			}
			if plan.Description != "" {
				pricingPlan.Description = &plan.Description
			}

			if err := tx.Create(&pricingPlan).Error; err != nil {
				tx.Rollback()
				s.logger.Error("料金プラン作成エラー", "gym", seedData.Name, "error", err)
				break
			}
		}

		// トランザクションコミット
		if err := tx.Commit().Error; err != nil {
			s.logger.Error("ジムデータコミットエラー", "name", seedData.Name, "error", err)
			continue
		}

		count++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ファイル読み込みエラー: %w", err)
	}

	s.logger.Info("ジムシード完了", "作成数", count)
	return nil
}