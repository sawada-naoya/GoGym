require 'pycall/import'
include PyCall::Import

pyimport 'lib.tasks.recommend', as: :recommender

class RecommendationService
  # 指定されたユーザーに対するレコメンデーションを提供するクラスメソッド
  def self.recommended_gyms_for_user(user)
      # ユーザーがレビューしたジムのIDを取得
    reviewed_gym_ids = user.reviews.pluck(:gym_id)
    # 各ジムIDについて、Pythonのrecommender.recommendメソッドを呼び出して、ユーザーに対するレコメンデーションを生成。結果をフラットにして重複を除去。
    recommended_gym_ids = reviewed_gym_ids.map do |gym_id|
      recommender.recommend(user.id, gym_id)
    end.flatten.uniq
    # レコメンドされたジムIDに対応するジムをデータベースから取得し、最大３件まで返す
    Gym.where(id: recommended_gym_ids).limit(3)
  end

  def self.popular_gyms
    Gym.joins(:reviews)
      .group('gyms.id')
      .order('COUNT(reviews.id) DESC')
      .limit(3)
  end
end
