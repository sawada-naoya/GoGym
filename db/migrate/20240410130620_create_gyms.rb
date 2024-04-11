class CreateGyms < ActiveRecord::Migration[7.0]
  def change
    create_table :gyms do |t|
      t.string :name, null: false
      t.string :membership_fee
      t.string :business_hours
      t.string :access
      t.string :remarks
      t.string :photos
      t.string :website

      t.timestamps
    end
  end
end
