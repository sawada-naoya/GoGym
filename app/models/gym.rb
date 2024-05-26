class Gym < ApplicationRecord
  has_one :location, dependent: :destroy
  has_many :reviews
  # 親モデルのフォームから子モデルの属性を直接受け取ることができるようにする
  accepts_nested_attributes_for :location
  validates :name, :access, presence: true

  # addressメソッドでLocationモデルのaddressを返す
  def address
    location&.address
  end
end
