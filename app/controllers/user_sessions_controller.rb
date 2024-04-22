class UserSessionsController < ApplicationController
  # skip_before_action :require_login, only: %i[new create]
  def new; end

  def create
    @user = login(params[:email], params[:password])
    if @user
      redirect_to root_path
      flash[:success] = t('flash.login_success')
    else
      flash.now[:danger] = t('flash.login_failure')
      render :new, status: :unprocessable_entity
    end
  end

  def destroy
    logout
    redirect_to root_path, status: :see_other
    flash[:success] = t('flash.logout')
  end
end
