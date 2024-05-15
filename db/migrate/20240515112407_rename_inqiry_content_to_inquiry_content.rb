class RenameInqiryContentToInquiryContent < ActiveRecord::Migration[7.0]
  def change
    rename_column :inquiries, :inqiry_content, :inquiry_content
  end
end
