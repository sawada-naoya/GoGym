class GymsController < ApplicationController

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

end
