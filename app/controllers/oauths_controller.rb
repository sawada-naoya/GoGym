class OauthsController < ApplicationController
  skip_before_action :require_login, raise: false

  def oauth
    #指定されたプロバイダの認証ページにユーザーをリダイレクトさせる
    login_at(params[:provider])
  end

  def callback
    provider = params[:provider]
    # 既存のユーザーをプロバイダ情報を元に検索し、存在すればログイン
    if @user = login_from(provider)
      redirect_to root_path, notice:"#{provider.titleize}アカウントでログインしました"
    else
      begin
        @user = create_from(provider)

        reset_session
        auto_login(@user)
        redirect_to root_path, notice:"#{provider.titleize}アカウントでログインしました"
      rescue
        redirect_to root_path, danger:"#{provider.titleize}アカウントでのログインに失敗しました"
      end
    end
  end

  private

  def auth_params
    params.permit(:code, :provider)
  end
end
