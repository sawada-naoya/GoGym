class Location < ApplicationRecord
  geocoded_by :address
  after_validation :geocode
  belongs_to :gym
end
