class GymsController < ApplicationController

  def search
    @gyms = @q.result
  end

  private

  def set_q
    @q = Gym.ransack(params[:q])
  end
end
