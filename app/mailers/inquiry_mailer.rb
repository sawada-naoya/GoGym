class InquiryMailer < ApplicationMailer
  default from: "from@example.com"
  layout "mailer"

  def new_inquiry(inquiry)
    @inquiry = inquiry
    mail(from: @inquiry.email, to: 'example@gmail.com', subject: 'Webサイトより問い合わせが届きました') do |format|
      format.text
    end
  end
end
