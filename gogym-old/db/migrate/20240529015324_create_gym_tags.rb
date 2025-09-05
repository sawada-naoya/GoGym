class CreateGymTags < ActiveRecord::Migration[7.0]
  def change
    create_table :gym_tags do |t|
      t.references :gym, null: false, foreign_key: true
      t.references :tag, null: false, foreign_key: true

      t.timestamps
    end
    add_index :gym_tags, [:gym_id, :tag_id], unique: true
  end
end
