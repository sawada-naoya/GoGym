package main

import (
	"fmt"
	"log"
	"os"

	"gogym-api/internal/configs"
	"gogym-api/internal/infra/db"

	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("APP_ENV") == "" {
		_ = os.Setenv("APP_ENV", "development")
	}
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load(".env.local")
		_ = godotenv.Overload(".env")
	}
}

func main() {
	config, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	database, err := db.NewDB(config.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, _ := database.DB()
	defer sqlDB.Close()

	// 既存のプリセットデータを削除
	result := database.Exec("DELETE FROM workout_parts WHERE user_id IS NULL")
	if result.Error != nil {
		log.Fatalf("Failed to delete existing parts: %v", result.Error)
	}
	fmt.Printf("Deleted %d existing workout parts\n", result.RowsAffected)

	// 正しいUTF-8でデータを再挿入
	parts := []struct {
		Name      string
		IsDefault bool
	}{
		{"胸", true},
		{"腕", true},
		{"肩", true},
		{"背中", true},
		{"脚", true},
	}

	for _, part := range parts {
		result := database.Exec(
			"INSERT INTO workout_parts (name, is_default, user_id) VALUES (?, ?, NULL)",
			part.Name,
			part.IsDefault,
		)
		if result.Error != nil {
			log.Fatalf("Failed to insert part '%s': %v", part.Name, result.Error)
		}
		fmt.Printf("✓ Inserted: %s\n", part.Name)
	}

	fmt.Println("\nSuccessfully fixed workout parts data!")

	// 確認のためデータを取得
	var results []struct {
		ID   int
		Name string
	}
	database.Raw("SELECT id, name FROM workout_parts WHERE user_id IS NULL ORDER BY id").Scan(&results)

	fmt.Println("\nCurrent data:")
	for _, r := range results {
		fmt.Printf("  ID=%d, Name=%s\n", r.ID, r.Name)
	}
}
