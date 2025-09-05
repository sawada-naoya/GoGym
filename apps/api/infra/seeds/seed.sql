-- Seed data for development
-- Insert sample users
INSERT INTO users (email, password_hash, display_name) VALUES
('john@example.com', '$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$xyz123hash', 'John Doe'),
('jane@example.com', '$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$abc456hash', 'Jane Smith');

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

-- Insert sample gyms with spatial data (Tokyo locations)
INSERT INTO gyms (name, description, location, address, city, prefecture, postal_code) VALUES
(
    'フィットネスクラブ渋谷',
    '最新設備を備えた総合フィットネスクラブです。プールやサウナも完備。',
    ST_SRID(POINT(139.7016, 35.6598), 4326),
    '東京都渋谷区渋谷1-2-3 渋谷ビル5F',
    '渋谷区',
    '東京都',
    '150-0002'
),
(
    'ジム新宿24',
    '24時間営業のマシン特化型ジムです。駅から徒歩2分の好立地。',
    ST_SRID(POINT(139.7006, 35.6895), 4326),
    '東京都新宿区新宿3-1-1 新宿駅前ビル3F',
    '新宿区',
    '東京都',
    '160-0022'
),
(
    'パワーハウス池袋',
    'フリーウェイトに特化した本格派トレーニングジム。',
    ST_SRID(POINT(139.7111, 35.7295), 4326),
    '東京都豊島区南池袋1-28-1 西武百貨店12F',
    '豊島区',
    '東京都',
    '171-8569'
);

-- Associate gyms with tags
INSERT INTO gym_tags (gym_id, tag_id) VALUES
-- フィットネスクラブ渋谷
(1, 3), -- プール
(1, 4), -- サウナ
(1, 5), -- 駐車場あり
(1, 6), -- シャワー
(1, 7), -- 女性専用エリア
(1, 9), -- 駅近

-- ジム新宿24
(2, 1), -- 24時間営業
(2, 6), -- シャワー
(2, 8), -- マシン豊富
(2, 9), -- 駅近

-- パワーハウス池袋
(3, 6), -- シャワー
(3, 8), -- マシン豊富
(3, 9); -- 駅近

-- Insert sample reviews
INSERT INTO reviews (user_id, gym_id, rating, comment, photos) VALUES
(
    1, 
    1, 
    4, 
    '設備が充実していて清潔感もあります。プールでの水泳も楽しめました。',
    JSON_ARRAY('https://example.com/photos/gym1-1.jpg', 'https://example.com/photos/gym1-2.jpg')
),
(
    2, 
    1, 
    5, 
    '女性専用エリアがあるので安心してトレーニングできます。スタッフの対応も親切です。',
    JSON_ARRAY('https://example.com/photos/gym1-3.jpg')
),
(
    1, 
    2, 
    5, 
    '24時間営業なので仕事帰りでも利用できて便利です。マシンも新しくて使いやすい。',
    JSON_ARRAY()
),
(
    2, 
    3, 
    4, 
    'フリーウェイトが充実している本格派のジム。本気でトレーニングしたい人におすすめ。',
    JSON_ARRAY('https://example.com/photos/gym3-1.jpg')
);

-- Insert sample favorites
INSERT INTO favorites (user_id, gym_id) VALUES
(1, 1),
(1, 2),
(2, 1),
(2, 3);