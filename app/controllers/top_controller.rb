class TopController < ApplicationController
  # skip_before_action :require_login, only: %i[index]

  def index
    @q = Gym.ransack(params[:q])
    @gyms = @q.result(distinct: true).page(params[:page])
  end

  def search
    @gyms = @q.result
  end

  private

  def set_q
    @q = Gym.ransack(params[:q])
  end

end
