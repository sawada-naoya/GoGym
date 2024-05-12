class GymsController < ApplicationController
  # ユーザーがログインしているかどうかをチェック
  before_action :require_login, only: [:new, :create]

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
      render :new, status: :unprocessable_entity
      flash.now[:danger] = t('flash.gym_create_failure')
    end
  end

  def index
    @gyms = @q.result(distinct: true).page(params[:page]).per(5)
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
  def edit; end

  private

  def gym_params
    params.require(:gym).permit(:name, :membership_fee, :business_hours, :access, :remarks, :website,
    location_attributes: [:address])
  end

end
