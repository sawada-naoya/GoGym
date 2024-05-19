class GymsController < ApplicationController
  before_action :require_login, only: [:new, :create, :edit, :update]
  before_action :set_gym, only: [:show, :edit, :update, :images_index]
  before_action :calculate_average_rating, only: [:show, :images_index]

  # GET /gyms/new
  def new
    @gym = Gym.new
    @gym.build_location
  end

  # POST /gyms
  def create
    @gym = Gym.new(gym_params)
    if @gym.save
      redirect_to @gym
      flash[:success] = t('flash.gym_create_success')
    else
      flash.now[:danger] = t('flash.gym_create_failure')
      render :new, status: :unprocessable_entity
    end
  end

  def index
    # if params[:q][:location_address_cont].present?
    #   # まず、地理的な検索を行う
    #   @gyms = Gym.near(params[:q][:location_address_cont], 10)
    #   @q = @gyms.ransack(params[:q])
    # else
    #   # 地理的な検索がない場合、通常の検索を行う
    #   @q = Gym.ransack(params[:q])
    # end
    # 検索結果をページネーションで区切る
    @gyms = @q.result(distinct: true).page(params[:page]).per(5)
    @average_ratings = calculate_average_ratings_for_gyms(@gyms)
    @gym_images = get_gym_images(@gyms)
  end

  def images_index
    @gym = Gym.find(params[:id])
    @reviews = @gym.reviews.where.not(image: nil).page(params[:page]).per(9)
  end

  def show
    @gym = Gym.includes(:location, :reviews).find(params[:id])
    @gyms = Gym.includes(:location).all
  end

  # データの編集画面を表示
  def edit
    @gym = Gym.find(params[:id])
  end

  def update
    @gym = Gym.find(params[:id])
    if @gym.update(gym_params)
      redirect_to @gym
      flash[:success] = t('flash.gym_update_success')
    else
      render :edit, status: :unprocessable_entity
      flash.now[:danger] = t('flash.gym_update_failure')
    end
  end

  private

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
    params.require(:gym).permit(:name, :membership_fee, :business_hours, :access, :remarks, :website,
    location_attributes: [:address])
  end

  def require_login
    unless logged_in?
      message = if action_name == 'edit'
                  "編集を行うにはログインが必要です。"
                elsif action_name == 'new' || action_name == 'create'
                  "登録するにはログインが必要です。"
                else
                  "この操作を行うにはログインが必要です。"
                end
      flash[:danger] = message
      redirect_to login_path
    end
  end

  def get_gym_images(gyms)
    images = {}
    gyms.each do |gym|
      review_with_image = gym.reviews.where.not(image: nil).first
      images[gym.id] = review_with_image&.image || 'fake'
    end
    images
  end
end
