class ApplicationController < ActionController::Base

  private

  def not_authenticated
    redirect_to login_path
    flash[:warning] = t('flash.require_login')
  end
end
