require 'rails_helper'

RSpec.describe User, type: :model do
  # ここからテストスタート
  # user作成に対するテスト
  describe 'バリデーションチェック' do
    let(:user) { FactoryBot.build(:user) }

    it '設定したバリデーションが機能しているか' do
      expect(user).to be_valid
    end
    it '名前が未入力だとinvalidになるか' do
      user.name = nil
      expect(user).to be_invalid
    end
    it '名前が256文字以上だとinvalidになるか' do
      user.name = 'a'* 256
      expect(user).to be_invalid
    end
    it 'メールが未入力だとinvalidになるか' do
      user.email = nil
      expect(user).to be_invalid
    end
    it 'メールが重複しているとinvalidになるか' do
      user.save
      user2 = FactoryBot.build(:user, email: user.email)
      expect(user2).to be_invalid
    end
    it 'パスワードが未入力だとinvalidになるか' do
      user.password = nil
      expect(user).to be_invalid
    end
    it 'パスワードが3文字未満だとinvalidになるか' do
      user.password = 'a'* 2
      expect(user).to be_invalid
    end
    it 'パスワードが確認入力と一致しない場合invalidになるか' do
      user.password = 'password'
      user.password_confirmation = 'different_password'
      expect(user).to be_invalid
    end
    it 'パスワード確認が未入力だとinvalidになるか' do
      user.password_confirmation = nil
      expect(user).to be_invalid
    end
  end
end
