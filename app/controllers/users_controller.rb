class UsersController < ApplicationController
    # ユーザーがログインしているかどうかをチェック
    before_action :require_login, only: [:show]

  def new
    @user = User.new
  end

  def create
    @user = User.new(user_params)
    if @user.save
      redirect_to root_path
      flash[:success] = t('flash.user_registration_success')
    else
      flash.now[:danger] = t('flash.user_registration_failure')
      render :new, status: :unprocessable_entity
    end
  end

  def show
    @user = User.find(params[:id])
    # @favorite_gyms = @user.favorite_gyms
    @reviews = @user.reviews
  end

  private

  def user_params
    params.require(:user).permit(:email, :name, :password, :password_confirmation)
  end
end
