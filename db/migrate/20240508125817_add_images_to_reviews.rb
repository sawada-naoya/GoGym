class AddImagesToReviews < ActiveRecord::Migration[7.0]
  def change
    add_column :reviews, :images, :json
  end
end
