class Review < ApplicationRecord
  belongs_to :gym
  belongs_to :user

  mount_uploader :image, ImageUploader

  validates :content, presence: true, length: { maximum: 101 }
end
