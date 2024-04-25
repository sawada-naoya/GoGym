require 'rails_helper'

RSpec.describe User, type: :model do
  it '名前、メールがあり、パスワードは3文字以上であれば有効であること' do
    user = build(:user, password: 'password', password_confirmation: 'password')
    expect(user).to be_valid
  end

  it 'メールはユニークであること' do
    user = create(:user, password: 'password', password_confirmation: 'password')
    user2 = build(:user, email: user.email)
    expect(user2).not_to be_valid
    expect(user2.errors[:email]).to include('はすでに存在します')
  end

  it 'メールアドレス名前は必須項目であること' do
    user = build(:user, password: 'password', password_confirmation: 'password')
    user.email = nil
    user.name = nil
    user.valid?
    expect(user.errors[:email]).to include('を入力してください')
    expect(user.errors[:name]).to include('を入力してください')
  end

  it '名は255文字以下であること' do
    user = build(:user, password: 'password', password_confirmation: 'password')
    user.name = 'a' * 256
    user.valid?
    expect(user.errors[:name]).to include('は255文字以内で入力してください')
  end
end
