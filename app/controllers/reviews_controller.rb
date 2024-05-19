class ReviewsController < ApplicationController
  before_action :set_gym, only: [:index, :new, :create]
  before_action :set_review, only: [:edit, :update, :destroy]

  def index
    @reviews = @gym.reviews.page(params[:page]).per(5)
    @average_rating = @gym.reviews.average(:rating).to_f.round(2)
  end

  def new
    @review = Review.new
    @average_rating = @gym.reviews.average(:rating).to_f.round(2)
  end

  def create
    @review = @gym.reviews.build(review_params.merge(user: current_user))
    if @review.save
      redirect_to gym_reviews_path(@gym)
      flash[:success] = t('flash.review_create_success')
    else
      flash.now[:danger] = t('flash.review_create_failure')
      render :new, status: :unprocessable_entity
    end
  end

  def show
    @gym = Gym.includes(:location).find(params[:id])
    @gyms = Gym.includes(:location).all
    @average_rating = @gym.reviews.average(:rating).to_f.round(2)
  end

  # データの編集画面を表示
  def edit; end

  def update
    if @review.update(review_params)
      redirect_to gym_reviews_path(@review.gym)
      flash[:success] = t('flash.review_update_success')
    else
      flash.now[:danger] = t('flash.review_update_failure')
      render :edit, status: :unprocessable_entity
    end
  end

  def destroy
    @review.destroy
  end

  def user_index
    @user = User.find(params[:user_id])
    @reviews = @user.reviews.page(params[:page]).per(5)
  end

  private

  def set_gym
    @gym = Gym.includes(:location).find(params[:gym_id])
  end

  def set_review
    @review = current_user.reviews.find(params[:id])
    @gym = @review.gym
  end

  def review_params
    params.require(:review).permit(:title, :content, :image, :rating)
  end

end
