class Inquiry < ApplicationRecord
    validates :name, presence: true, length: { maximum: 20 }
    # VALID_EMAIL_REGEX = /\A[\w+\-.]+@[a-z\d\-.]+\.[a-z]+\z/i,format: { with: VALID_EMAIL_REGEX }

    validates :email, presence: true, length: { maximum: 30 }
    validates :inquiry_content, presence: true, length: { maximum: 500 }
end
