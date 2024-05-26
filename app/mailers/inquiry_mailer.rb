class InquiryMailer < ApplicationMailer
  # default from: "zetianzhimi7@gmail.com"
  layout "mailer"

  def new_inquiry(inquiry)
    @inquiry = inquiry
    mail(from: @inquiry.email, to: 'zetianzhimi7@gmail.com', subject: 'Webサイトより問い合わせが届きました') do |format|
      format.text
    end
  end
end
