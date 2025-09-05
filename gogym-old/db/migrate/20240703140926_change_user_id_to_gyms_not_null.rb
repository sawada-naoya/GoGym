class ChangeUserIdToGymsNotNull < ActiveRecord::Migration[7.0]
  def change
    change_column_null :gyms, :user_id, false
  end
end
