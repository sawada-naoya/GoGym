crumb :root do
  link "Home", root_path
end

crumb :user_new do
  link "会員登録", new_user_path
  parent :root
end

crumb :user_session do
  link "ログイン", login_path
  parent :root
end

crumb :user_show do |user|
  link "#{user.name}さんのマイページ", user_path(user)
  parent :root
end

crumb :gyms do
  link "ジム一覧", gyms_path
  parent :root
end

crumb :gym do |gym|
  link gym.name, gym_path(gym)
  parent :gyms
end

crumb :new_gym do
  link "ジムの登録", new_gym_path
  parent :gyms
end

crumb :gym_reviews do |gym|
  link "レビュー一覧", gym_reviews_path(gym_id: gym.id)
  parent :gym, gym
end

crumb :gym_reviews do |gym|
  link "口コミ一覧", gym_reviews_path(gym_id: gym.id)
  parent :gym, gym
end

crumb :new_review do |gym|
  link "口コミ投稿", new_gym_review_path(gym)
  parent :gym, gym
end

crumb :edit_review do |gym, review|
  link "口コミ編集", edit_review_path(gym, review)
  parent :gym, gym
end

crumb :edit_gym do |gym|
  link "ジムの編集", edit_gym_path(gym)
  parent :gym, gym
end

# Handling Locations
crumb :locations do
  link "位置情報検索", locations_path
  parent :gyms
end

crumb :inquiry do
  link "お問い合わせ", new_inquiry_path
  parent :root
end

crumb :terms do
  link "利用規約", terms_path
  parent :root
end

crumb :privacy do
  link "プライバシーポリシー", privacy_path
  parent :root
end

# タグの詳細ページ
crumb :tag do |tag|
  link tag.name, tag_path(tag)
  parent :gyms
end
