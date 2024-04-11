class CreateSocialProviders < ActiveRecord::Migration[7.0]
  def change
    create_table :social_providers do |t|
      t.integer :string

      t.timestamps
    end
  end
end
