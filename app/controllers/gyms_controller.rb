require './app/services/recommend_service'

class GymsController < ApplicationController
  before_action :require_login, only: [:new, :create, :edit, :update]
  before_action :set_gym, only: [:show, :edit, :update, :images_index]
  before_action :calculate_average_rating, only: [:show, :images_index]
  before_action :load_tags, only: [:new, :create, :edit, :update]

  # GET /gyms/new
  def new
    @gym = Gym.new
    @gym.build_location
  end

  # POST /gyms
  def create
    @gym = current_user.gyms.build(gym_params)
    if @gym.save
      redirect_to @gym
      flash[:success] = t('flash.gym_create_success')
    else
      flash.now[:danger] = t('flash.gym_create_failure')
      render :new, status: :unprocessable_entity
    end
  end

  def index
    if params[:tag_id]
      @tag = Tag.find(params[:tag_id])
      @gyms = @tag.gyms.page(params[:page]).per(5)
    else
      @gyms = @q.result(distinct: true).page(params[:page]).per(5)
    end
    @average_ratings = calculate_average_ratings_for_gyms(@gyms)
    @gym_images = get_gym_images(@gyms)

    if logged_in?
      begin
        Rails.logger.info "現在のユーザーID: #{current_user.id}"
        @recommended_gyms = RecommendService.new(current_user).call
        Rails.logger.info "推奨ジム: #{@recommended_gyms.pluck(:id)}"
      rescue => e
        Rails.logger.error "エラー: #{e.message}"
        Rails.logger.error e.backtrace.join("\n")
        @recommended_gyms = []
      end
    else
      @recommended_gyms = []
    end

    if @recommended_gyms.empty?
      @popular_gyms = Gym.order('view_count DESC').limit(3) # よく見られているジムを取得
    end
  end

  def images_index
    @reviews = @gym.reviews.where.not(image: nil).page(params[:page]).per(9)
  end

  def show
    @gym = Gym.includes(:location, :reviews).find(params[:id])
    @gyms = Gym.includes(:location).all
    @tags = @gym.tags
    increment_view_count
  end

  # データの編集画面を表示
  def edit
  end

  def update
    if @gym.update(gym_params)
      redirect_to @gym
      flash[:success] = t('flash.gym_update_success')
    else
      render :edit, status: :unprocessable_entity
      flash.now[:danger] = t('flash.gym_update_failure')
    end
  end

  def autocomplete
    @gyms = Gym.where("name like ?", "%#{params[:q]}%").limit(6)
    render partial: 'shared/autocomplete', locals: { gyms: @gyms }
  end

  private

  def increment_view_count
    @gym.increment!(:view_count)
  end

  def set_gym
    @gym = Gym.find(params[:id])
  end

  def calculate_average_rating
    @average_rating = @gym.reviews.average(:rating).to_f.round(2)
  end

  def calculate_average_ratings_for_gyms(gyms)
    # each_with_objectメソッドは、指定したオブジェクト（ここでは空のハッシュ{}）を使って各ジムを処理
    gyms.each_with_object({}) do |gym, ratings|
      # ratings[gym.id]にジムの平均点を設定することで、ジムのIDをキーとして、平均評価を値とするハッシュを構築
      ratings[gym.id] = gym.reviews.average(:rating).to_f.round(2)
    end
  end

  def gym_params
    params.require(:gym).permit(:name, :membership_fee, :business_hours, :access, :remarks, :website, tag_ids: [],
    location_attributes: [:address])
  end

  def get_gym_images(gyms)
    images = {}
    gyms.each do |gym|
      review_with_image = gym.reviews.where.not(image: nil).first
      images[gym.id] = review_with_image&.image&.url || 'fake'
    end
    Rails.logger.debug { "Gym Images: #{images.inspect}" }
    images
  end

  def load_tags
    @tags = Tag.all
  end
end
