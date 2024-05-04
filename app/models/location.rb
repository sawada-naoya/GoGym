class Location < ApplicationRecord
  belongs_to :gym

  # 保存前に経度と緯度をセットする
  before_save :geocode_endpoints

  private

  def geocode_endpoints
    # addressが変更されている場合のみgeocodeを実行
    if address_changed?
      geocoded = Geocoder.search(address).first
      if geocoded
        self.latitude = geocoded.latitude
        self.longitude = geocoded.longitude
      end
    end
  end
end
