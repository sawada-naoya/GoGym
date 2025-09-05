require 'rails_helper'

RSpec.describe Location, type: :model do
  describe 'バリデーションチェック' do
    let(:location) { FactoryBot.build(:location) }

    it '設定したバリデーションが機能しているか' do
      expect(location).to be_valid
    end
    it 'アドレスが空の場合invalidになるか' do
      location.address = nil
      expect(location).to be_invalid
    end
  end
end
