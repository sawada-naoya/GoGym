require 'rails_helper'

RSpec.describe Review, type: :model do
  describe 'バリデーションチェック' do
    let(:review) { FactoryBot.build(:review) }

    it '設定したバリデーションが機能しているか' do
      expect(review).to be_valid
    end
    it '内容が未入力だとinvalidになるか' do
      review.content = nil
      expect(review).to be_invalid
    end
    it '内容が102文字以上だとinvalidになるか' do
      review.content = 'a' * 102
      expect(review).to be_invalid
    end
    it '評価が未入力だとinvalidになるか' do
      review.rating = nil
      expect(review).to be_invalid
    end
    it '評価が1未満または5を超える値の場合invalidになるか' do
      review.rating = 0.9
      expect(review).to be_invalid

      review.rating = 5.1
      expect(review).to be_invalid
    end
  end
end
