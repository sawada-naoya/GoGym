class GymsController < ApplicationController

  def index
    @gyms = @q.result(distinct: true).page(params[:page])
  end

  # def search
  #   @gyms = @q.result
  # end

  def show
    # URL パラメーターから受け取った ID に基づいて、指定された ID の掲示板をデータベースから検索
    @gym = Gym.find(params[:id])
    # ログインしているユーザーに関連付けられた新しいコメントオブジェクトを生成
    # @comment = current_user.comments.new
    # # 指定された掲示板に紐付くコメントを検索
    # # @gym.comments は、関連するコメントのコレクションを取得
    # # includes(:user) は、コメントに紐付くユーザー情報も同時に読み込むように指示
    # # order(created_at: :desc) は、コメントを作成された日時を降順に並べ替え
    # @comments = @gym.comments.includes(:user).order(created_at: :desc)
  end

  # データの編集画面を表示
  def edit; end

  def update
    # もし、gym_paramsで許可された属性情報に基づいて掲示板が更新できた場合、
    # つまり、新しい属性情報がデータベースに保存された場合は以下のコードを実行する。
    if @gym.update(gym_params)
      redirect_to gym_path(@gym)
      flash[:success] = t('flash.gym_update_success')
    else
      flash.now[:danger] = t('flash.gym_update_failure')
      render :edit, status: :unprocessable_entity
    end
  end

  def destroy
    @gym.destroy!
    redirect_to gyms_path, status: :see_other
    flash[:success] = t('flash.gym_delete_success')
  end

  private

  def set_q
    @q = Gym.ransack(params[:q])
  end
end
