class TagsController < ApplicationController
  def show
    @tag = Tag.find(params[:id])
    @gyms = @tag.gyms.page(params[:page]).per(10)
    @average_ratings = calculate_average_ratings_for_gyms(@gyms)
    @gym_images = get_gym_images(@gyms)
  end

  private

  def calculate_average_ratings_for_gyms(gyms)
    gyms.each_with_object({}) do |gym, ratings|
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
