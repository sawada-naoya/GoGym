class Favorite < ApplicationRecord
  belongs_to :user
  belongs_to :gym

  validates :user_id, uniqueness: { scope: :gym_id }
end
