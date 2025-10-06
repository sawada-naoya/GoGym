CREATE TABLE exercises (
    id          CHAR(26) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    is_builtin  TINYINT(1)   NOT NULL DEFAULT 0,
    created_by  CHAR(26)     NULL, -- users.id（NULLは内蔵プリセット）
    created_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP    NULL,
    CONSTRAINT fk_exercises_user
        FOREIGN KEY (created_by) REFERENCES users(id),
    UNIQUE KEY uq_exercise_name_user (name, created_by)
) ENGINE=InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE workouts (
    id                CHAR(26)  PRIMARY KEY,
    user_id           CHAR(26)  NOT NULL, -- users.id
    performed_at      DATETIME  NOT NULL,
    note              TEXT      NULL,
    condition         TINYINT   NULL, -- 1〜5
    place             VARCHAR(100) NULL,
    started_at        DATETIME  NULL,
    ended_at          DATETIME  NULL,
    duration_minutes  INT
      GENERATED ALWAYS AS (
        CASE
          WHEN started_at IS NULL OR ended_at IS NULL THEN NULL
          ELSE TIMESTAMPDIFF(MINUTE, started_at, ended_at)
        END
      ) STORED, -- 開始/終了から自動算出
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at        TIMESTAMP NULL,
    CONSTRAINT fk_workouts_user
      FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT chk_workouts_condition_range
      CHECK (condition IS NULL OR (condition BETWEEN 1 AND 5)),
    CONSTRAINT chk_workouts_time_order
      CHECK (started_at IS NULL OR ended_at IS NULL OR started_at < ended_at),
    KEY idx_workouts_user_time (user_id, performed_at),
    KEY idx_workouts_time_range (started_at, ended_at)
) ENGINE=InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;


CREATE TABLE workout_sets (
    id          CHAR(26)     PRIMARY KEY,
    workout_id  CHAR(26)     NOT NULL,
    exercise_id CHAR(26)     NOT NULL,
    set_number  INT          NOT NULL,       -- 同一workout内の順序（1,2,3...）
    weight      DECIMAL(6,2) NOT NULL,       -- kgで保存
    rep         INT          NOT NULL,
    max         DECIMAL(6,2) NULL,           -- 推定1RM（weightとrepから保存時に算出）
    note        VARCHAR(255) NULL,
    created_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP    NULL,
    CONSTRAINT fk_sets_workout
      FOREIGN KEY (workout_id)  REFERENCES workouts(id)  ON DELETE CASCADE,
    CONSTRAINT fk_sets_exercise
      FOREIGN KEY (exercise_id) REFERENCES exercises(id),
    CONSTRAINT chk_positive_rep
      CHECK (rep > 0),
    CONSTRAINT chk_nonneg_weight
      CHECK (weight >= 0),
    UNIQUE KEY uq_set_number (workout_id, exercise_id, set_number),
    KEY idx_sets_exercise_time (exercise_id, created_at),
    KEY idx_sets_workout (workout_id, set_number)
) ENGINE=InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
