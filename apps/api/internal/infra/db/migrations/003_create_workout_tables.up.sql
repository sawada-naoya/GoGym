CREATE TABLE IF NOT EXISTS workout_parts (
  id          INT AUTO_INCREMENT PRIMARY KEY,
  name        VARCHAR(50) NOT NULL,
  user_id     CHAR(26) NULL,                          -- users.id（NULL=プリセット）
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at  TIMESTAMP NULL,
  CONSTRAINT fk_workout_parts_user
    FOREIGN KEY (user_id) REFERENCES users(id),
  UNIQUE KEY uq_workout_parts_name_user (name, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS workout_exercises (
  id                INT AUTO_INCREMENT PRIMARY KEY,
  name              VARCHAR(100) NOT NULL,
  workout_part_id  INT NULL,
  user_id           CHAR(26) NULL,                    -- users.id（NULL=プリセット）
  created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at        TIMESTAMP NULL,
  CONSTRAINT fk_workout_exercises_part
    FOREIGN KEY (workout_part_id) REFERENCES workout_parts(id) ON DELETE SET NULL,
  CONSTRAINT fk_workout_exercises_user
    FOREIGN KEY (user_id) REFERENCES users(id),
  UNIQUE KEY uq_workout_exercises_name_user (name, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS workout_records (
  id               INT AUTO_INCREMENT PRIMARY KEY,
  user_id          CHAR(26) NOT NULL,                 -- users.id
  performed_date   DATE NOT NULL,
  started_at       DATETIME NULL,
  ended_at         DATETIME NULL,
  place            VARCHAR(100) NULL,
  note             TEXT NULL,
  condition_level  TINYINT NULL,
  duration_minutes INT NULL,                          -- アプリ側で計算して保存
  created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at       TIMESTAMP NULL,
  CONSTRAINT fk_workout_records_user
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS workout_sets (
  id                     INT AUTO_INCREMENT PRIMARY KEY,
  workout_record_id     INT NOT NULL,
  workout_exercise_id   INT NOT NULL,
  set_number             INT NOT NULL,
  weight_kg              DECIMAL(6,2) NOT NULL DEFAULT 0.00,
  reps                   INT NOT NULL DEFAULT 0,
  estimated_max          DECIMAL(6,2) NULL,
  note                   VARCHAR(255) NULL,
  created_at             TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at             TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at             TIMESTAMP NULL,
  CONSTRAINT fk_workout_sets_record
    FOREIGN KEY (workout_record_id)  REFERENCES workout_records(id)  ON DELETE CASCADE,
  CONSTRAINT fk_workout_sets_exercise
    FOREIGN KEY (workout_exercise_id) REFERENCES workout_exercises(id),
  UNIQUE KEY uq_workout_sets_order (workout_record_id, workout_exercise_id, set_number),
  KEY idx_workout_sets_exercise_time (workout_exercise_id, created_at),
  KEY idx_workout_sets_record (workout_record_id, set_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
