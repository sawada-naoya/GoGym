-- Seed data for development (simplified schema)
-- Insert sample users
INSERT INTO users (email, password_hash, display_name) VALUES
('test1@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '山田太郎'),
('test2@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '佐藤花子'),
('test3@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '田中健一'),
('admin@gogym.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '管理者');

-- Insert sample tags
INSERT INTO tags (name) VALUES
('24時間営業'),
('パーソナルトレーニング'),
('プール'),
('サウナ'),
('駐車場あり'),
('シャワー'),
('女性専用エリア'),
('マシン豊富'),
('駅近');

-- Insert sample gyms (simplified schema - no spatial data for now)
INSERT INTO gyms (name, description, location_latitude, location_longitude, address, city, prefecture, postal_code, is_active) VALUES
(
    'ゴールドジム 原宿東京',
    '本格的なウェイトトレーニングの聖地。豊富なフリーウェイトとマシンで、初心者から上級者まで満足できる設備。24時間営業で忙しい方にも最適。',
    35.6705,
    139.7026,
    '東京都渋谷区神宮前6-31-17 ベロックスビル B1F',
    '渋谷区',
    '東京都',
    '150-0001',
    true
),
(
    'ゴールドジム 原宿ANNEX',
    '原宿東京の姉妹店として、よりスペシャライズされたトレーニング環境を提供。コンパクトながら充実した設備で集中してトレーニングが可能。',
    35.6698,
    139.7038,
    '東京都渋谷区神宮前6-28-6 Q PLAZA HARAJUKU 4F',
    '渋谷区',
    '東京都',
    '150-0001',
    true
),
(
    'ゴールドジム 表参道東京',
    '表参道の中心地にある都市型フィットネスクラブ。洗練された空間で本格的なトレーニングを。ビジネスマンや美意識の高い女性に人気。',
    35.6654,
    139.7106,
    '東京都港区北青山3-6-19 三和実業表参道ビル B1F',
    '港区',
    '東京都',
    '107-0061',
    true
),
(
    'ゴールドジム 渋谷東京',
    '渋谷の中心部にありアクセス抜群。若い世代から年配の方まで幅広い層に愛される老舗店舗。充実した設備とフレンドリーなスタッフが魅力。',
    35.6581,
    139.7016,
    '東京都渋谷区渋谷1-23-16 cocoti SHIBUYA 6F',
    '渋谷区',
    '東京都',
    '150-0002',
    true
),
(
    'ゴールドジム 南青山東京',
    '青山エリアの高級感あふれる立地で、質の高いトレーニング環境を提供。女性向け設備も充実し、美意識の高い会員が多く集まる。',
    35.6627,
    139.7172,
    '東京都港区南青山2-27-20 植竹ビル B1F',
    '港区',
    '東京都',
    '107-0062',
    true
);

-- Associate gyms with tags
INSERT INTO gym_tags (gym_id, tag_id) VALUES
-- ゴールドジム 原宿東京
(1, 1), -- 24時間営業
(1, 8), -- マシン豊富
(1, 9), -- 駅近

-- ゴールドジム 原宿ANNEX
(2, 8), -- マシン豊富
(2, 9), -- 駅近

-- ゴールドジム 表参道東京
(3, 2), -- パーソナルトレーニング
(3, 6), -- シャワー
(3, 7), -- 女性専用エリア
(3, 9), -- 駅近

-- ゴールドジム 渋谷東京
(4, 6), -- シャワー
(4, 8), -- マシン豊富
(4, 9), -- 駅近

-- ゴールドジム 南青山東京
(5, 2), -- パーソナルトレーニング
(5, 6), -- シャワー
(5, 7), -- 女性専用エリア
(5, 9); -- 駅近

-- Insert sample reviews
INSERT INTO reviews (user_id, gym_id, rating, comment) VALUES
(1, 1, 5, '設備が充実していて、トレーナーの方も親切です。フリーウェイトエリアが広くて使いやすい。'),
(2, 1, 4, '施設は綺麗ですが、混雑時は待ち時間が長いことがあります。スタッフの対応は良いです。'),
(3, 2, 4, '24時間利用できるのが便利。マシンの種類も豊富で満足しています。'),
(1, 2, 3, 'マシン中心なのでフリーウェイトをしたい人には物足りないかも。でも清潔で使いやすいです。');

-- Insert sample favorites
INSERT INTO favorites (user_id, gym_id) VALUES
(1, 1),
(1, 2),
(2, 1),
(3, 3);