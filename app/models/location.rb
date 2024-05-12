class Location < ApplicationRecord
  belongs_to :gym

  geocoded_by :address
  after_validation :geocode
  validates :address, presence: true
end
