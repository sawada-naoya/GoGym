class CreateAdminUserAndUpdateGyms < ActiveRecord::Migration[7.0]
  def up
    # 管理者ユーザーを作成
    admin_user = User.find_or_create_by!(email: 'admin@example.com') do |user|
      user.name = 'Admin User'
      user.password = 'adminpassword'
      user.password_confirmation = 'adminpassword'
    end

    # 既存のGymデータのuser_idがnilの場合、管理者ユーザーIDを設定
    Gym.where(user_id: nil).update_all(user_id: admin_user.id)
  end
end
