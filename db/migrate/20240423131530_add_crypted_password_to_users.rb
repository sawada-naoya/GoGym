class AddCryptedPasswordToUsers < ActiveRecord::Migration[7.0]
  def up
    add_column :users, :crypted_password, :string
  end

  def down
    remove_column :users, :crypted_password
  end
end
