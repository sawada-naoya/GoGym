class User < ApplicationRecord

  validates :name, presence: { message: 'を入力してください' }, length: { maximum: 255 }
  validates :email, presence: { message: 'を入力してください' }, uniqueness: true
  validates :password, presence: { message: 'を入力してください' }, length: { minimum: 3, message: "は3文字以上で入力してください" }

end
