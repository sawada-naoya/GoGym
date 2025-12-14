-- ============================================
-- GoGym Database Schema
-- ============================================

-- Users table: アプリケーションユーザー
CREATE TABLE users (
    id CHAR(26) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Gyms table: ユーザーが登録したジム（最小限の情報）
CREATE TABLE gyms (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    latitude  DECIMAL(10,7) NOT NULL,
    longitude DECIMAL(10,7) NOT NULL,
    source_url VARCHAR(1000) NOT NULL,
    primary_photo_url VARCHAR(1000) NULL,
    place_id VARCHAR(128) NULL,
    created_by CHAR(26) NOT NULL,
    updated_by CHAR(26) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_gyms_location (latitude, longitude),
    INDEX idx_gyms_name (name),
    -- place_id は Google Places 導入時の重複防止用
    -- NULL は複数許容する前提
    UNIQUE KEY uq_gyms_place_id (place_id),
    CONSTRAINT fk_gyms_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_gyms_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Refresh tokens table: JWT リフレッシュトークン管理
CREATE TABLE refresh_tokens (
    jti CHAR(26) NOT NULL PRIMARY KEY,
    user_id CHAR(26) NOT NULL,
    revoked_at DATETIME NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_active (user_id, revoked_at, expires_at),
    CONSTRAINT fk_refresh_tokens_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Workout parts table: トレーニング部位（胸、背中、脚など）
CREATE TABLE workout_parts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    user_id CHAR(26) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_workout_parts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uq_workout_parts_name_user (name, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Workout exercises table: 種目（ベンチプレス、スクワットなど）
CREATE TABLE workout_exercises (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    workout_part_id INT NULL,
    user_id CHAR(26) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_workout_exercises_part FOREIGN KEY (workout_part_id) REFERENCES workout_parts(id) ON DELETE SET NULL,
    CONSTRAINT fk_workout_exercises_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uq_workout_exercises_name_user (name, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Workout records table: トレーニング記録（日付・場所）
CREATE TABLE workout_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id CHAR(26) NOT NULL,
    gym_id BIGINT NULL,
    performed_date DATE NOT NULL,
    started_at DATETIME NULL,
    ended_at DATETIME NULL,
    note TEXT NULL,
    condition_level TINYINT NULL,
    duration_minutes INT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_workout_records_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_workout_records_gym FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE SET NULL,
    INDEX idx_workout_records_user_date (user_id, performed_date),
    INDEX idx_workout_records_gym (gym_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Workout sets table: セット詳細（重量・回数）
CREATE TABLE workout_sets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    workout_record_id INT NOT NULL,
    workout_exercise_id INT NOT NULL,
    set_number INT NOT NULL,
    weight_kg DECIMAL(6,2) NOT NULL DEFAULT 0.00,
    reps INT NOT NULL DEFAULT 0,
    estimated_max DECIMAL(6,2) NULL,
    note VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_workout_sets_record FOREIGN KEY (workout_record_id) REFERENCES workout_records(id) ON DELETE CASCADE,
    CONSTRAINT fk_workout_sets_exercise FOREIGN KEY (workout_exercise_id) REFERENCES workout_exercises(id) ON DELETE CASCADE,
    UNIQUE KEY uq_workout_sets_order (workout_record_id, workout_exercise_id, set_number),
    KEY idx_workout_sets_exercise_time (workout_exercise_id, created_at),
    KEY idx_workout_sets_record (workout_record_id, set_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
