class AddGymIdToLocations < ActiveRecord::Migration[7.0]
  def change
    add_column :locations, :gym_id, :integer
  end
end
