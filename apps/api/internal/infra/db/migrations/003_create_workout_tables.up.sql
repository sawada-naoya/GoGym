CREATE TABLE IF NOT EXISTS workout_parts (
  id          INT AUTO_INCREMENT PRIMARY KEY,
  name        VARCHAR(50) NOT NULL,
  is_default  TINYINT(1) NOT NULL DEFAULT 0,
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
  is_default        TINYINT(1) NOT NULL DEFAULT 0,
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

-- 部位（プリセット）
INSERT INTO workout_parts (name, is_default, user_id)
VALUES ('胸',1,NULL), ('腕',1,NULL), ('肩',1,NULL), ('背中',1,NULL), ('脚',1,NULL)
ON DUPLICATE KEY UPDATE name = VALUES(name);


-- 2) 種目（プリセット）
INSERT INTO workout_exercises (name, workout_part_id, is_default, user_id)
VALUES
  -- 胸
  ('ベンチプレス',                 (SELECT id FROM workout_parts WHERE name='胸' AND user_id IS NULL), 1, NULL),
  ('ダンベルプレス',               (SELECT id FROM workout_parts WHERE name='胸' AND user_id IS NULL), 1, NULL),
  ('インクラインベンチプレス',     (SELECT id FROM workout_parts WHERE name='胸' AND user_id IS NULL), 1, NULL),
  ('インクラインダンベルプレス',   (SELECT id FROM workout_parts WHERE name='胸' AND user_id IS NULL), 1, NULL),
  ('スミスマシンベンチプレス',     (SELECT id FROM workout_parts WHERE name='胸' AND user_id IS NULL), 1, NULL),
  ('スミスマシンインクラインベンチプレス', (SELECT id FROM workout_parts WHERE name='胸' AND user_id IS NULL), 1, NULL),
  ('ペックフライ',                 (SELECT id FROM workout_parts WHERE name='胸' AND user_id IS NULL), 1, NULL),
  ('ケーブルクロスオーバー',       (SELECT id FROM workout_parts WHERE name='胸' AND user_id IS NULL), 1, NULL),

  -- 脚
  ('スクワット',                   (SELECT id FROM workout_parts WHERE name='脚' AND user_id IS NULL), 1, NULL),
  ('ブルガリアンスクワット',       (SELECT id FROM workout_parts WHERE name='脚' AND user_id IS NULL), 1, NULL),
  ('レッグプレス',                 (SELECT id FROM workout_parts WHERE name='脚' AND user_id IS NULL), 1, NULL),
  ('レッグカール',                 (SELECT id FROM workout_parts WHERE name='脚' AND user_id IS NULL), 1, NULL),
  ('レッグエクステンション',       (SELECT id FROM workout_parts WHERE name='脚' AND user_id IS NULL), 1, NULL),

  -- 腕
  ('バーベルアームカール',         (SELECT id FROM workout_parts WHERE name='腕' AND user_id IS NULL), 1, NULL),
  ('ダンベルアームカール',         (SELECT id FROM workout_parts WHERE name='腕' AND user_id IS NULL), 1, NULL),
  ('プリチャーカール',             (SELECT id FROM workout_parts WHERE name='腕' AND user_id IS NULL), 1, NULL),
  ('インクラインダンベルカール',   (SELECT id FROM workout_parts WHERE name='腕' AND user_id IS NULL), 1, NULL),
  ('ケーブルプレスダウン',         (SELECT id FROM workout_parts WHERE name='腕' AND user_id IS NULL), 1, NULL),
  ('ケーブルアームカール',         (SELECT id FROM workout_parts WHERE name='腕' AND user_id IS NULL), 1, NULL),
  ('スカルクラッシャー',           (SELECT id FROM workout_parts WHERE name='腕' AND user_id IS NULL), 1, NULL),

  -- 背中
  ('チンニング',                   (SELECT id FROM workout_parts WHERE name='背中' AND user_id IS NULL), 1, NULL),
  ('ラットプルダウン',             (SELECT id FROM workout_parts WHERE name='背中' AND user_id IS NULL), 1, NULL),
  ('ベントオーバーロウ',           (SELECT id FROM workout_parts WHERE name='背中' AND user_id IS NULL), 1, NULL),
  ('ハーフデッドリフト',           (SELECT id FROM workout_parts WHERE name='背中' AND user_id IS NULL), 1, NULL),

  -- 肩
  ('サイドレイズ',                 (SELECT id FROM workout_parts WHERE name='肩' AND user_id IS NULL), 1, NULL),
  ('ダンベルショルダープレス',     (SELECT id FROM workout_parts WHERE name='肩' AND user_id IS NULL), 1, NULL),
  ('バーベルショルダープレス',     (SELECT id FROM workout_parts WHERE name='肩' AND user_id IS NULL), 1, NULL),
  ('スミスマシンショルダープレス', (SELECT id FROM workout_parts WHERE name='肩' AND user_id IS NULL), 1, NULL),
  ('リアデルト',                   (SELECT id FROM workout_parts WHERE name='肩' AND user_id IS NULL), 1, NULL)
ON DUPLICATE KEY UPDATE
  workout_part_id = VALUES(workout_part_id),
  is_default       = VALUES(is_default);
