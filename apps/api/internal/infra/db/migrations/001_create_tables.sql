-- +goose Up
-- 初期スキーマ（MySQL 8.0 / InnoDB / utf8mb4）
-- SRID 4326 を前提に地理検索・日本語検索(簡易)・正規化を両立

-- Users
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE COLLATE utf8mb4_0900_ai_ci,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    -- 追加プロフィール
    profile_photo_url VARCHAR(500) COMMENT 'プロフィール写真URL',
    bio TEXT COMMENT '自己紹介文',
    location VARCHAR(100) COMMENT '居住地域',
    birth_year YEAR NULL COMMENT '生年（統計用途）',
    gender ENUM('male', 'female', 'other', 'not_specified') DEFAULT 'not_specified' COMMENT '性別',
    is_verified BOOLEAN DEFAULT FALSE COMMENT '認証済みフラグ',
    last_login_at DATETIME NULL COMMENT '最終ログイン日時（TZ依存を避けDATETIME）',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_email (email),
    INDEX idx_users_location (location),
    INDEX idx_users_last_login (last_login_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Gyms
CREATE TABLE gyms (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    -- 位置情報: POINT(SRID=4326)
    location POINT NOT NULL SRID 4326,
    address VARCHAR(500) NOT NULL,
    city VARCHAR(100),
    prefecture VARCHAR(100),
    postal_code VARCHAR(10),

    -- 検索用 FTS（生成列は STORED にして FULLTEXT 対応を確実化）
    fts TEXT GENERATED ALWAYS AS (CONCAT(COALESCE(name,''),' ',COALESCE(description,''))) STORED,

    -- 追加メタ
    phone_number VARCHAR(32) COMMENT '電話番号（E.164考慮）',
    website VARCHAR(2048) COMMENT '公式サイトURL',
    -- opening_hours JSON は二重管理回避のため採用しない（詳細は gym_hours に正規化）
    price_range_min INT UNSIGNED COMMENT '最低料金（円）',
    price_range_max INT UNSIGNED COMMENT '最高料金（円）',
    access_info TEXT COMMENT 'アクセス（最寄り駅等）',
    parking_info TEXT COMMENT '駐車場情報',
    amenities JSON COMMENT '設備（配列）',
    capacity INT UNSIGNED COMMENT '定員',
    age_restrictions TEXT COMMENT '年齢制限・条件',
    operator_name VARCHAR(100) COMMENT '運営会社',
    brand_name VARCHAR(150) COMMENT 'ブランド名（チェーン）',

    -- 集計（トリガで整合性維持）
    average_rating DECIMAL(3,2) NOT NULL DEFAULT 0.00 COMMENT '平均評価',
    review_count INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'レビュー数',

    -- ステータス
    is_active BOOLEAN NOT NULL DEFAULT TRUE COMMENT '公開/非公開',
    is_temporarily_closed BOOLEAN NOT NULL DEFAULT FALSE COMMENT '一時休業',
    closure_reason TEXT COMMENT '休業理由',

    -- SEO
    slug VARCHAR(255) UNIQUE COMMENT 'SEOスラッグ',
    meta_title VARCHAR(100) COMMENT 'SEOタイトル',
    meta_description VARCHAR(300) COMMENT 'SEO説明',

    -- 画像（メイン/JSONは継続。将来は gym_photos に移行可能）
    main_photo_url VARCHAR(500) COMMENT 'メイン写真URL',
    photos JSON COMMENT '追加写真URL配列',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Index
    SPATIAL INDEX idx_gyms_location (location),
    FULLTEXT INDEX idx_gyms_fts (fts),
    INDEX idx_gyms_city (city),
    INDEX idx_gyms_prefecture (prefecture),
    INDEX idx_gyms_price_range (price_range_min, price_range_max),
    INDEX idx_gyms_brand (brand_name),
    INDEX idx_gyms_active (is_active, is_temporarily_closed),
    INDEX idx_gyms_slug (slug),
    INDEX idx_gyms_rating (average_rating),
    INDEX idx_gyms_review_count (review_count),
    INDEX idx_gyms_active_region (is_active, prefecture, city),
    CONSTRAINT chk_gyms_price_range CHECK (
      price_range_min IS NULL OR price_range_max IS NULL OR price_range_min <= price_range_max
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tags
CREATE TABLE tags (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tags_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Gym-Tag map
CREATE TABLE gym_tags (
    gym_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (gym_id, tag_id),
    FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    INDEX idx_gym_tags_tag_id (tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Reviews
CREATE TABLE reviews (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    gym_id BIGINT NOT NULL,
    rating TINYINT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    photos JSON,
    -- 詳細評価
    helpful_count INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '参考になった数',
    visit_date DATE COMMENT '訪問日',
    visit_purpose ENUM('trial', 'regular', 'day_pass', 'other') COMMENT '訪問目的',
    cleanliness_rating TINYINT COMMENT '清潔さ 1-5',
    staff_rating TINYINT COMMENT 'スタッフ 1-5',
    equipment_rating TINYINT COMMENT '設備 1-5',
    value_rating TINYINT COMMENT 'コスパ 1-5',
    is_verified_visit BOOLEAN NOT NULL DEFAULT FALSE COMMENT '訪問確認済み',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE,

    INDEX idx_reviews_user_id (user_id),
    INDEX idx_reviews_gym_id (gym_id),
    INDEX idx_reviews_rating (rating),
    INDEX idx_reviews_created_at (created_at),

    -- 詳細評価・ユースケース別索引
    INDEX idx_reviews_helpful (helpful_count),
    INDEX idx_reviews_visit_date (visit_date),
    INDEX idx_reviews_cleanliness (cleanliness_rating),
    INDEX idx_reviews_staff (staff_rating),
    INDEX idx_reviews_equipment (equipment_rating),
    INDEX idx_reviews_value (value_rating),
    INDEX idx_reviews_gym_helpful (gym_id, helpful_count),

    UNIQUE KEY unique_user_gym_review (user_id, gym_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Favorites
CREATE TABLE favorites (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    gym_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE,
    INDEX idx_favorites_user_id (user_id),
    INDEX idx_favorites_gym_id (gym_id),
    UNIQUE KEY unique_user_gym_favorite (user_id, gym_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Refresh tokens
CREATE TABLE refresh_tokens (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_refresh_tokens_user_id (user_id),
    INDEX idx_refresh_tokens_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 「参考になった」明細（重複投票防止）
CREATE TABLE review_helpfuls (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    review_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
    INDEX idx_review_helpfuls_user_id (user_id),
    INDEX idx_review_helpfuls_review_id (review_id),
    UNIQUE KEY unique_user_review_helpful (user_id, review_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ジムの営業時間（正規化）
CREATE TABLE gym_hours (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    gym_id BIGINT NOT NULL,
    day_of_week TINYINT NOT NULL COMMENT '0=Sun .. 6=Sat',
    open_time TIME NULL,
    close_time TIME NULL,
    is_closed BOOLEAN NOT NULL DEFAULT FALSE,
    notes VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE,
    INDEX idx_gym_hours_gym_day (gym_id, day_of_week),
    UNIQUE KEY unique_gym_day (gym_id, day_of_week)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 料金プラン（検索は plan_type + price が軸）
CREATE TABLE gym_pricing_plans (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    gym_id BIGINT NOT NULL,
    plan_name VARCHAR(100) NOT NULL,
    plan_type ENUM('monthly', 'annual', 'day_pass', 'trial') NOT NULL,
    price INT UNSIGNED NOT NULL,
    description TEXT,
    benefits JSON,
    is_popular BOOLEAN NOT NULL DEFAULT FALSE,
    display_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE,
    INDEX idx_pricing_gym_id (gym_id),
    INDEX idx_pricing_type_price (plan_type, price),
    INDEX idx_pricing_popular (is_popular),
    INDEX idx_pricing_order (display_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 写真（将来の並べ替え/モデレーションに備える）
CREATE TABLE gym_photos (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    gym_id BIGINT NOT NULL,
    url VARCHAR(500) NOT NULL,
    width INT UNSIGNED NULL,
    height INT UNSIGNED NULL,
    caption VARCHAR(200) NULL,
    display_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (gym_id) REFERENCES gyms(id) ON DELETE CASCADE,
    INDEX idx_gym_photos_gym (gym_id),
    INDEX idx_gym_photos_active (is_active, display_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Reviews 集計の一貫性維持トリガは権限の問題で後ほど追加
-- TODO: 本番環境では適切な権限設定後にトリガーを追加

-- +goose Down
-- 逆順で安全に削除
-- トリガーは作成していないのでスキップ

DROP TABLE IF EXISTS gym_photos;
DROP TABLE IF EXISTS gym_pricing_plans;
DROP TABLE IF EXISTS gym_hours;
DROP TABLE IF EXISTS review_helpfuls;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS favorites;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS gym_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS gyms;
DROP TABLE IF EXISTS users;
