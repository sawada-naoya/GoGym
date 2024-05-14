class InquiriesController < ApplicationController
  def new
    @inquiry = Inquiry.new
  end

  def create
    @inquiry = Inquiry.new(inquiry_params)
    if @inquiry.save
      InquiryMailer.new_inquiry(@inquiry).deliver_now
      InquiryMailer.inquiry_mailer(@inquiry).deliver_now
      flash[:success] = t('flash.inquiry_send_success')
      redirect_to root_path
    else
      flash.now[:danger] = t('flash.inquiry_send_failure')
      render :new
    end
  end

  private

  def inquiry_params
    params.require(:inquiry).permit(:name, :email, :inqiry_content)
  end
end
