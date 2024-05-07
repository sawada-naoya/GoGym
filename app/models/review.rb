class Review < ApplicationRecord
  belongs_to :gym
  belongs_to :user

  validates :content, presence: true, length: { maximum: 101 }
end
