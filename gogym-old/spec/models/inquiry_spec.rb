require 'rails_helper'

RSpec.describe Inquiry, type: :model do
  describe 'バリデーションチェック' do
    let(:inquiry) { FactoryBot.build(:inquiry) }

    it '設定したバリデーションが機能しているか' do
      expect(inquiry).to be_valid
    end
    it '名前が未入力だとinvalidになるか' do
      inquiry.name = nil
      expect(inquiry).to be_invalid
    end
    it '名前が21文字以上だとinvalidになるか' do
      inquiry.name = 'a' * 21
      expect(inquiry).to be_invalid
    end
    it 'メールが未入力だとinvalidになるか' do
      inquiry.email = nil
      expect(inquiry).to be_invalid
    end
    it 'メールが31文字以上だとinvalidになるか' do
      inquiry.email = 'a' * 32
      expect(inquiry).to be_invalid
    end
    it 'お問い合せ内容が未入力だとinvalidになるか' do
      inquiry.inquiry_content = nil
      expect(inquiry).to be_invalid
    end
    it 'お問い合せ内容が501文字以上だとinvalidになるか' do
      inquiry.inquiry_content = 'a' * 502
      expect(inquiry).to be_invalid
    end
  end
end
