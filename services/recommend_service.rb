class RecommendService

  def initialize(user)
    @user = user
  end

  def call
    recommended_gyms
  end

  private

  #　コンテンツフィルタリング
    # - 現在のユーザーが評価を行ったジムをデータベースから取得
    # - 取得したジムに関連するタグを集計し、出現頻度が高い順に3つのタグを選ぶ
    # - 選定したタグに基づいて、関連するジムをデータベースから取得

  def content_based_recommended_gyms
    # 現在のユーザーが評価したすべてのジムのIDを取得し、配列として返す
    user_gym_ids = @user.reviews.pluck(:gym_id)
    # ユーザーが評価したジムに関連する人気のタグのIDを取得
    tag_ids = find_popular_tags(user_gym_ids)
    # 取得したタグのIDに一致するタグを持つジムを重複を除いてフィルタリング
    Gym.joins(:tags).where(tags: { id: tag_ids }).distinct
  end

  # 特定のジムに関連するタグをその出現頻度に基づいて取得
  def find_popular_tags(gym_ids)
    # tagsテーブルとgymsテーブルを結合
    Tag.joins(:gyms)
      .where(gyms: { id: gym_ids }) # 指定されたジムIDに関連するタグをフィルタリング
      .group('tags.id') # タグIDごとにグルーピング
      .order('COUNT(tags.id) DESC') # タグの出現頻度を降順でソート
      .limit(3) # 上位3つのタグを取得
      .pluck(:id) # 取得したタグのIDを配列として返す
  end


  #　ユーザー・ベースの協調フィルタリング
    # - 自分と同じジムを３つ以上お気に入り登録しているユーザーを探す
    # - 各ジムに対する平均ratingが＋ー0.2以内の平均ratingを持つユーザーを探す
    # - 見つけたユーザーの評価しているジムの中で、ratingがそのジムの平均ratingよりも高く評価しているジムを３つ探す。
    # - それを各ユーザーのビューに渡して３つのおすすめのジムを表示させる。

  def user_based_recommendations
    similar_users = find_similar_users
    find_highly_rated_gyms(similar_users).uniq.take(3)
  end

  # 現在のユーザーと類似した行動を持つユーザーを見つける
  def find_similar_users
    # 現在のユーザーがお気に入りに登録しているジムのIDを取得
    user_gym_ids = @user.favorites.pluck(:gym_id)

    # お気に入りにしているジムが3つ以上重複している他のユーザーを見つける
    similar_users = User.joins(:favorites) # usersテーブルとfavoritesテーブルを結合
                        .where(favorites: { gym_id: user_gym_ids }) # 共通のジムを持つユーザーをフィルタリング
                        .group('users.id') # ユーザーごとにグループ化
                        .having('COUNT(favorites.gym_id) >= 3') # 3つ以上の共通ジムを持つユーザーを選ぶ
                        .where.not(id: @user.id) #　現在のユーザー自身を除外

    # 現在のユーザーのジム評価を取得
    user_gym_ratings = @user.reviews.where(gym_id: user_gym_ids).pluck(:gym_id, :rating).to_h
    # 評価が±0.2以内のユーザーを選ぶ
    similar_users.select do |similar_user|
      # 類似ユーザーの共通ジムの評価を取得し、all?メソッドを使ってすべてのレビューが±0.2以内の評価差であることを確認
      similar_user.reviews.where(gym_id: user_gym_ids).all? do |review|
        (user_gym_ratings[review.gym_id] - review.rating).abs <= 0.2
      end
    end
  end

  def find_highly_rated_gyms(similar_users)
    # 先に見つけた類似ユーザーのリストからユーザーのIDを配列として取得
    similar_user_ids = similar_users.map(&:id)
    # 現在のユーザーがお気に入りに登録しているジムのIDを取得
    user_gym_ids = @user.favorites.pluck(:gym_id)

    # 各ジムに対するレビュー情報を取得
    Gym.joins(:reviews)
      .where(reviews: { user_id: similar_user_ids }) # 類似ユーザーによるレビューのみを対象とする
      .where.not(id: user_gym_ids) # 現在のユーザーがお気に入りにしているジムを除外
      .group('gyms.id') # ジムごとにグルーピング
      .having('AVG(reviews.rating) > (SELECT AVG(rating) FROM reviews WHERE gym_id = gyms.id)') # 各ジムのレビューの平均評価が、そのジムの全体平均評価を上回っているジムを対象とし、対象レビューの平均評価が全体平均評価より高いジムのみをフィルタリング
      .order('AVG(reviews.rating) DESC') # 平均評価が高い順にソート
      .limit(3) # 上位3つのジムを取得
  end
end
