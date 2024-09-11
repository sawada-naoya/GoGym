require 'rails_helper'

RSpec.describe Favorite, type: :model do
  describe 'バリデーションチェック' do
    let(:favorite) { FactoryBot.build(:favorite) }

    it '設定したバリデーションが機能しているか' do
      expect(favorite).to be_valid
    end
  end
end
