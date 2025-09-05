class DropSocialProvidersTable < ActiveRecord::Migration[7.0]
  def change
    drop_table :social_providers
  end
  def down
    create_table :social_providers do |t|
      t.integer :string
      t.timestamps null: false
    end
  end
end
