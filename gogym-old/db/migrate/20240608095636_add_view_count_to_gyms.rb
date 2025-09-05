class AddViewCountToGyms < ActiveRecord::Migration[7.0]
  def change
    add_column :gyms, :view_count, :integer, default: 0, null: false
  end
end
