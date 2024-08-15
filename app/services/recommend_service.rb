class RecommendService

  def initialize(user)
    @user = user
    @logger = Rails.logger
  end

  def call
    @logger.info "RecommendServiceを呼び出しています: ユーザーID=#{@user.id}"
    recommended_gyms
  end

  private

  #　コンテンツフィルタリング
    # ① 現在のユーザーが評価したジム（例: ゴールドジム渋谷店）をデータベースから取得します。
    # ② 取得したジムに関連するタグ（例: 24時間、サウナ、プロテイン、駐車場、パーソナル）を確認し、同様のタグを持つジムを検索します。検索されたジムに設定されているタグの中で、出現頻度が高い順に3つのタグ（例: プロテイン、24時間、サウナ）を選びます。
    # ③ 選定された3つのタグに基づいて、同様のタグを持つジムをデータベースからすべて取得します。

  def content_based_recommended_gyms
    @logger.info "コンテンツベースのフィルタリングを実行中: ユーザーID=#{@user.id}"
    user_gym_ids = @user.reviews.pluck(:gym_id)
    @logger.info "ユーザーが評価したジムのIDを取得: #{user_gym_ids}"
    tag_ids = find_popular_tags(user_gym_ids)
    @logger.info "人気のタグのIDを取得: #{tag_ids}"
    gyms = Gym.joins(:tags).where(tags: { id: tag_ids }).distinct
    @logger.info "タグに基づいてフィルタリングされたジム: #{gyms.pluck(:id)}"
    gyms
  end

  # 特定のジムに関連するタグをその出現頻度に基づいて取得
  def find_popular_tags(gym_ids)
    @logger.info "人気のタグを検索: ジムID=#{gym_ids}"
    tags = Tag.joins(:gyms)
              .where(gyms: { id: gym_ids }) # 指定されたジムIDに関連するタグをフィルタリング
              .group('tags.id') # タグIDごとにグルーピング
              .order('COUNT(tags.id) DESC') # タグの出現頻度を降順でソート
              .limit(3) # 上位3つのタグを取得
              .pluck(:id) # 取得したタグのIDを配列として返す
    @logger.info "人気のタグを取得: #{tags}"
    tags
  end


  #　ユーザー・ベースの協調フィルタリング
    # ① 自分と同じジムを3つ以上お気に入り登録しているユーザーを探します。
    # ② そのユーザーの平均評価が、自分の平均評価の±0.2以内のユーザーを探します。
    # ③ ②で見つけたユーザーが評価しているジムの中で、自身の平均評価よりも高く評価しているジムを探します。
    # ④ それを各ユーザーのビューに渡して、3つのおすすめのジムとして取得します。

  def user_based_recommendations
    @logger.info "ユーザーベースのフィルタリングを実行中: ユーザーID=#{@user.id}"
    similar_users = find_similar_users
    @logger.info "類似ユーザーを見つけました: ユーザーID=#{similar_users.map(&:id)}"
    gyms = find_highly_rated_gyms(similar_users).uniq.take(3)
    @logger.info "高評価のジムを見つけました: ジムID=#{gyms.map(&:id)}"
    gyms
  end

  def calculate_average_rating(user, gym_ids)
    ratings = user.reviews.where(gym_id: gym_ids).group(:gym_id).average(:rating)
    return 0 if ratings.empty?

    ratings.values.sum / ratings.size
  end

  # 現在のユーザーと類似した行動を持つユーザーを見つける
  def find_similar_users
    # 現在のユーザーがお気に入りに登録しているジムのIDを取得
    user_gym_ids = @user.favorites.pluck(:gym_id)
    @logger.info "ユーザーのお気に入りジムIDを取得: #{user_gym_ids}"

    # お気に入りにしているジムが3つ以上重複している類似ユーザーを見つける
    similar_users = User.joins(:favorites)
                        .where(favorites: { gym_id: user_gym_ids }) # 共通のジムを持つユーザーをフィルタリング
                        .group('users.id') # ユーザーごとにグループ化
                        .having('COUNT(favorites.gym_id) >= 3') # 3つ以上の共通ジムを持つユーザーを選ぶ
                        .where.not(id: @user.id) #　現在のユーザー自身を除外
    @logger.info "お気に入りが3つ以上重複しているユーザーを検索クエリ実行: #{similar_users.to_sql}"
    @logger.info "お気に入りが3つ以上重複しているユーザーを検索: 類似ユーザーID=#{similar_users.map(&:id)}"

    # 現在のユーザーのジムの平均評価を取得
    user_average_rating = calculate_average_rating(@user, user_gym_ids)
    @logger.info "現在のユーザーの平均評価: #{user_average_rating}"

    # 現在のユーザーの平均評価より±0.2以内の類似ユーザーを選ぶ
    selected_users = similar_users.select do |similar_user|
      similar_user_average_rating = calculate_average_rating(similar_user, user_gym_ids)
      @logger.info "類似ユーザーID: #{similar_user.id}, 平均評価: #{similar_user_average_rating}"

      next if similar_user_average_rating.nil?  # nilチェックを追加
      (user_average_rating - similar_user_average_rating).abs <= 0.2
    end.compact  # nilを除去

    @logger.info "評価が±0.2以内のユーザーを選択: #{selected_users.map(&:id)}"
    selected_users
  end

  def find_highly_rated_gyms(selected_users)
    similar_user_ids = selected_users.map(&:id)
    @logger.info "類似ユーザーのIDを取得: #{similar_user_ids}"
    user_gym_ids = @user.favorites.pluck(:gym_id)
    @logger.info "ユーザーのお気に入りジムIDを取得: #{user_gym_ids}"

    highly_rated_gyms = []

    selected_users.each do |user|
      user.reviews.each do |review|
        gym_id = review.gym_id
        gym_average_rating = calculate_average_rating(@user, [gym_id])
        @logger.info "ジムID: #{gym_id}, 類似ユーザー評価: #{review.rating}, ジム平均評価: #{gym_average_rating}"

        if review.rating > gym_average_rating
          highly_rated_gyms << gym_id
        end
      end
    end

    Gym.where(id: highly_rated_gyms.uniq)
  end

  # 最終的なおすすめジムの表示:
  # コンテンツベースフィルタリングとユーザーベースフィルタリングで取得されたジムを組み合わせ、重複したジムを最終的なおすすめジムとしてユーザーに表示します。

  def recommended_gyms
    user_based_gyms = user_based_recommendations
    content_based_gyms = content_based_recommended_gyms
    @logger.info "ユーザーベースの推奨ジム: #{user_based_gyms.pluck(:id)}"
    @logger.info "コンテンツベースの推奨ジム: #{content_based_gyms.pluck(:id)}"
    # 重複するジムを抽出
    duplicated_gyms = user_based_gyms & content_based_gyms

    # 重複がない、もしくは2つ未満の場合は、ユーザーベースのジムを優先して表示
    if duplicated_gyms.size >= 2
      @recommended_gyms = duplicated_gyms
    else
      @recommended_gyms = user_based_gyms
    end

    # ログに表示されるジムのIDを記録
    Rails.logger.info "推奨ジム: #{@recommended_gyms.pluck(:id)}"
    @recommended_gyms
  end
end
