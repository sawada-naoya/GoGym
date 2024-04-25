class User < ApplicationRecord
  authenticates_with_sorcery!

  validates :name, presence: { message: 'を入力してください' }, length: { maximum: 255 }
  validates :email, presence: { message: 'を入力してください' }, uniqueness: true
  validates :password, presence: { message: 'を入力してください' }, length: { minimum: 3, message: "は3文字以上で入力してください" }, confirmation: true, if: -> { new_record? || changes[:crypted_password] }
  validates :password_confirmation, presence: true, if: -> { new_record? || changes[:crypted_password] }
  # validates :reset_password_token, presence: true, uniqueness: true, allow_nil: true

end
