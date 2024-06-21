class OauthsController < ApplicationController
  skip_before_action :require_login, raise: false

  def oauth
    Rails.logger.debug "プロバイダー: #{auth_params[:provider]}"
    #指定されたプロバイダの認証ページにユーザーをリダイレクトさせる
    login_at(auth_params[:provider])
  end

  def callback
    provider = auth_params[:provider]
    Rails.logger.debug "コールバック プロバイダー: #{provider}"
    # 既存のユーザーをプロバイダ情報を元に検索し、存在すればログイン
    if @user = login_from(provider)
      Rails.logger.debug "既存のユーザーが見つかりました: #{@user.inspect}"
      redirect_to root_path
      flash[:success] = "#{provider.titleize}アカウントでログインしました"
    else
      begin
        @user = create_from(provider)
        Rails.logger.debug "新しいユーザーが作成されました: #{@user.inspect}"

        reset_session
        auto_login(@user)
        Rails.logger.debug "新しいユーザーが自動的にログインされました: #{@user.inspect}"
        redirect_to root_path
        flash[:success] = "#{provider.titleize}アカウントでログインしました"
      rescue => e
        Rails.logger.error "ログインに失敗しました: #{e.message}"
        flash.now[:danger] = "#{provider.titleize}アカウントでのログインに失敗しました"
        redirect_to root_path
      end
    end
  end

  private

  def auth_params
    params.permit(:code, :provider, :scope, :authuser, :prompt)
  end
end
