class LocationsController < ApplicationController
  def index
    # @locations = Location.all
    @gyms = Gym.includes(:location).all
  end
end
