-- ============================================
-- GoGym Database Schema (PostgreSQL)
-- ============================================

-- Users table: アプリケーションユーザー
CREATE TABLE users (
    id CHAR(26) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_email ON users(email);

-- Gyms table: ユーザーが登録したジム（最小限の情報）
CREATE TABLE gyms (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    normalized_name VARCHAR(255) NOT NULL,
    latitude  DECIMAL(10,7) NOT NULL,
    longitude DECIMAL(10,7) NOT NULL,
    source_url VARCHAR(1000) NOT NULL,
    primary_photo_url VARCHAR(1000) NULL,
    place_id VARCHAR(128) NULL,
    created_by CHAR(26) NOT NULL,
    updated_by CHAR(26) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT uq_gyms_place_id UNIQUE (place_id),
    CONSTRAINT fk_gyms_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_gyms_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_gyms_location ON gyms(latitude, longitude);
CREATE INDEX idx_gyms_name ON gyms(name);
CREATE UNIQUE INDEX uq_gyms_created_by_normalized_name ON gyms(created_by, normalized_name) WHERE deleted_at IS NULL;

-- Refresh tokens table: JWT リフレッシュトークン管理
CREATE TABLE refresh_tokens (
    jti CHAR(26) NOT NULL PRIMARY KEY,
    user_id CHAR(26) NOT NULL,
    revoked_at TIMESTAMP NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_refresh_tokens_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_active ON refresh_tokens(user_id, revoked_at, expires_at);

-- Workout parts table: トレーニング部位（胸、背中、脚など）
CREATE TABLE workout_parts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    user_id CHAR(26) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_workout_parts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_workout_parts_name_user UNIQUE (name, user_id)
);

-- Workout exercises table: 種目（ベンチプレス、スクワットなど）
CREATE TABLE workout_exercises (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    workout_part_id INT NULL,
    user_id CHAR(26) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_workout_exercises_part FOREIGN KEY (workout_part_id) REFERENCES workout_parts(id) ON DELETE SET NULL,
    CONSTRAINT fk_workout_exercises_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_workout_exercises_name_user UNIQUE (name, user_id)
);

-- Workout records table: トレーニング記録（日付・場所）
CREATE TABLE workout_records (
    id SERIAL PRIMARY KEY,
    user_id CHAR(26) NOT NULL,
    gym_id BIGINT NULL,
    performed_date DATE NOT NULL,
    started_at TIMESTAMP NULL,
    ended_at TIMESTAMP NULL,
    note TEXT NULL,
    condition_level SMALLINT NULL,
    duration_minutes INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_workout_records_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_workout_records_gym FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE SET NULL
);

CREATE INDEX idx_workout_records_user_date ON workout_records(user_id, performed_date);
CREATE INDEX idx_workout_records_gym ON workout_records(gym_id);

-- Workout sets table: セット詳細（重量・回数）
CREATE TABLE workout_sets (
    id SERIAL PRIMARY KEY,
    workout_record_id INT NOT NULL,
    workout_exercise_id INT NOT NULL,
    set_number INT NOT NULL,
    weight_kg DECIMAL(6,2) NOT NULL DEFAULT 0.00,
    reps INT NOT NULL DEFAULT 0,
    estimated_max DECIMAL(6,2) NULL,
    note VARCHAR(255) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_workout_sets_record FOREIGN KEY (workout_record_id) REFERENCES workout_records(id) ON DELETE CASCADE,
    CONSTRAINT fk_workout_sets_exercise FOREIGN KEY (workout_exercise_id) REFERENCES workout_exercises(id) ON DELETE CASCADE,
    CONSTRAINT uq_workout_sets_order UNIQUE (workout_record_id, workout_exercise_id, set_number)
);

CREATE INDEX idx_workout_sets_exercise_time ON workout_sets(workout_exercise_id, created_at);
CREATE INDEX idx_workout_sets_record ON workout_sets(workout_record_id, set_number);
