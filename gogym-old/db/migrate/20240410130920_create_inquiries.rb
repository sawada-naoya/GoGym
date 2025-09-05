class CreateInquiries < ActiveRecord::Migration[7.0]
  def change
    create_table :inquiries do |t|
      t.string :name, null: false
      t.float :email, null: false
      t.text :inqiry_content, null: false

      t.timestamps
    end
  end
end
