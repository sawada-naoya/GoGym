class FavoritesController < ApplicationController
  # before_action :authenticate_user!

  def create
    @gym = Gym.find(params[:gym_id])
    current_user.favorite(@gym)
  end

  def destroy
    @gym = Gym.find(params[:gym_id])
    current_user.unfavorite(@gym)
  end
end
