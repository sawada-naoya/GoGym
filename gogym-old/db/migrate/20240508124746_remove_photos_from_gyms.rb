class RemovePhotosFromGyms < ActiveRecord::Migration[7.0]
  def change
    remove_column :gyms, :photos, :string
  end
end
