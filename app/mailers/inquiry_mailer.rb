class InquiryMailer < ApplicationMailer
  default from: 'noreply@example.com'
  layout 'mailer'

  def new_inquiry(inquiry)
    @inquiry = inquiry
    mail(from: inquiry.email, to: ENV['MAIL_ADDRESS'], subject: 'Webサイトより問い合わせが届きました') do |format|
      format.text
    end
  end

  def inquiry_confirmation(inquiry)
    @inquiry = inquiry
    mail(to: inquiry.email, subject: 'お問い合わせありがとうございます') do |format|
      format.text { render 'inquiry_mailer' }
    end
  end
end
