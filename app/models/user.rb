class User < ApplicationRecord
  authenticates_with_sorcery!
  has_many :reviews
  has_many :favorites, dependent: :destroy
  has_many :favorite_gyms, through: :favorites, source: :gym
  has_many :authentications, dependent: :destroy

  validates :name, presence: { message: 'を入力してください' }, length: { maximum: 255 }
  validates :email, presence: { message: 'を入力してください' }, uniqueness: true
  validates :password, presence: { message: 'を入力してください' }, length: { minimum: 3, message: "は3文字以上で入力してください" }, confirmation: true, if: -> { new_record? || changes[:crypted_password] }
  validates :password_confirmation, presence: true, if: -> { new_record? || changes[:crypted_password] }
  # validates :reset_password_token, presence: true, uniqueness: true, allow_nil: true

  # 指定されたジム (gym) をユーザーのお気に入りリストに追加
  def favorite(gym)
    favorite_gyms << gym unless favorite?(gym)
  end

  # 指定されたジム (gym) をユーザーのお気に入りリストから削除
  def unfavorite(gym)
    favorite_gyms.delete(gym) if favorite?(gym)
  end

  # 指定されたジム (gym) がユーザーのお気に入りリストに含まれているかどうかを確認
  def favorite?(gym)
    favorite_gyms.include?(gym)
  end
end
