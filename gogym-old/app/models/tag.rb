class Tag < ApplicationRecord
  has_many :gym_tags
  has_many :gyms, through: :gym_tags
end
