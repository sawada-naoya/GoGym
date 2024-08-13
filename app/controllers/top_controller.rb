class TopController < ApplicationController
  # skip_before_action :require_login, only: %i[index]

  def index
    @gyms = Gym.includes(:location).all
  end

  def search
    @gyms = @q.result
  end

end
