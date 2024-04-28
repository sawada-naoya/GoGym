class ApplicationController < ActionController::Base
  before_action :set_q

  private

  def set_q
    @q = Gym.ransack(params[:q])
  end

  def not_authenticated
    redirect_to login_path
    flash[:warning] = t('flash.require_login')
  end
end
