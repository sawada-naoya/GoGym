class UsersController < ApplicationController
    before_action :require_login, only: [:show, :favorites]
    before_action :set_user, only: [:show, :favorites, :edit_review, :update_review, :destroy_review]
    before_action :set_review, only: [:edit_review, :update_review, :destroy_review]
    protect_from_forgery with: :null_session

  def new
    @user = User.new
  end

  def create
    @user = User.new(user_params)
    if @user.save
      redirect_to root_path
      flash[:success] = t('flash.user_registration_success')
    else
      flash.now[:danger] = t('flash.user_registration_failure')
      render :new, status: :unprocessable_entity
    end
  end

  def show
    @reviews = @user.reviews.order(updated_at: :desc)
  end

  def favorites
    @q = @user.favorite_gyms.ransack(params[:q])
    @favorite_gyms = @q.result(distinct: true).order(created_at: :desc).page(params[:page]).per(10)
    @average_ratings = calculate_average_ratings_for_gyms(@favorite_gyms)
    @gym_images = get_gym_images(@favorite_gyms)
  end

  def edit_review
  end

  def update_review
    if @review.update(review_params)
      redirect_to user_user_reviews_path(@user)
      flash[:success] = t('flash.review_update_success')
    else
      flash.now[:danger] = t('flash.review_update_failure')
      render :edit_review, status: :unprocessable_entity
    end
  end

  def destroy_review
    @review.destroy
  end

  private

  def user_params
    params.require(:user).permit(:email, :name, :password, :password_confirmation)
  end

  def set_user
    @user = User.find(params[:id])
  end

  def set_review
    @review = current_user.reviews.find(params[:id])
  end

  def review_params
    params.require(:review).permit(:content, :rating, :image)
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

  def get_gym_images(gyms)
    images = {}
    gyms.each do |gym|
      review_with_image = gym.reviews.where.not(image: nil).first
      images[gym.id] = review_with_image&.image&.url || 'fake'
    end
    images
  end
end
