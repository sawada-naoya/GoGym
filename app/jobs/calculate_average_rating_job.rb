class CalculateAverageRatingJob < ApplicationJob
  queue_as :default

  def perform(gym_id)
    gym = Gym.find(gym_id)
    average_rating = gym.reviews.average(:rating).to_f.round(2)
    gym.update(average_rating: average_rating)
  end
end
