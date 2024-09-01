class ChangeEmailTypeInInquiries < ActiveRecord::Migration[7.0]
  def up
    change_column :inquiries, :email, :string
  end

  def down
    change_column :inquiries, :email, :float
  end
end
