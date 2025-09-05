class AddUserIdToGyms < ActiveRecord::Migration[7.0]
  def change
    add_reference :gyms, :user, foreign_key: true, null: true
  end
end
