class ReviewsController < ApplicationController
  def index
    @gym = Gym.find(params[:gym_id])
    @reviews = @gym.reviews.page(params[:page]).per(5)
  end

  def new
    @gym = Gym.includes(:location).find(params[:gym_id])
    @review = Review.new
  end

  def create
    @gym = Gym.find(params[:gym_id])
    @review = @gym.reviews.build(review_params.merge(user: current_user))
    if @review.save
      redirect_to gym_reviews_path
      flash[:success] = t('flash.review_create_success')
    else
      flash.now[:danger] = t('flash.review_create_failure')
      render :new, status: :unprocessable_entity
    end
  end

  def show
    @gym = Gym.includes(:location).find(params[:id])
    @gyms = Gym.includes(:location).all
  end

  # データの編集画面を表示
  def edit; end

  def update
    # もし、gym_paramsで許可された属性情報に基づいて掲示板が更新できた場合、
    # つまり、新しい属性情報がデータベースに保存された場合は以下のコードを実行する。
    if @gym.update(gym_params)
      redirect_to gym_path(@gym)
      flash[:success] = t('flash.gym_update_success')
    else
      flash.now[:danger] = t('flash.gym_update_failure')
      render :edit, status: :unprocessable_entity
    end
  end

  def destroy
    @gym.destroy!
    redirect_to gyms_path, status: :see_other
    flash[:success] = t('flash.gym_delete_success')
  end

  private

  def review_params
    params.require(:review).permit(:title, :content, :image)
  end

end
