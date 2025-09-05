class ChangeImagesToImageInReviews < ActiveRecord::Migration[7.0]
  def change
    rename_column :reviews, :images, :image
    change_column :reviews, :image, :string
  end
end
