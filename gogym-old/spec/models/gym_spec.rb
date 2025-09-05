require 'rails_helper'

RSpec.describe Gym, type: :model do
  describe 'バリデーションチェック' do
    let(:user) { FactoryBot.create(:user) }
    let(:gym) { FactoryBot.build(:gym, user: user) }

    it '設定したバリデーションが機能しているか' do
      expect(gym).to be_valid
    end
    it '名前が未入力だとinvalidになるか' do
      gym.name = nil
      expect(gym).to be_invalid
    end
    it '住所が未入力だとinvalidになるか' do
      gym.access = nil
      expect(gym).to be_invalid
    end
  end
end
