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

  def require_login
    unless logged_in?
      message = case action_name
                when 'edit'
                  "編集を行うにはログインが必要です。"
                when 'new', 'create'
                  "登録するにはログインが必要です。"
                when 'favorite', 'unfavorite'
                  "お気に入り登録するにはログインが必要です。"
                else
                  "この操作を行うにはログインが必要です。"
                end
      flash[:danger] = message
      redirect_to login_path
    end
  end

  def logged_in?
    !!current_user
  end
end
