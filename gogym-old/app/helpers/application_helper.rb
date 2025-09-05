module ApplicationHelper
  def page_title(title = '')
    base_title = 'GoGym'
    title.present? ? "#{title} | #{base_title}" : base_title
  end

  def active_if(path)
    path == controller_path ? 'active' : ''
  end

  def default_meta_tags
    {
      site: 'GoGym',
      title: 'ジムの検索サービス',
      reverse: true,
      charset: 'utf-8',
      description: 'GoGymはジムの検索サービスです。あなたのベストジムを見つけてみませんか？ Find the best gyms with GoGym!',
      keywords: 'トレーニング,training,ジム,スポーツ,スポーツ施設,gym,スポーツジム,筋トレ,プロテイン,運動,ダイエット,減量,ボディメイク,bodymake',
      canonical: Rails.env.production? ? 'https://gogym-m40u.onrender.com' : request.original_url,
      separator: '|',
      og: {
        site_name: :site,
        title: :title,
        description: :description,
        type: 'website',
        url: Rails.env.production? ? 'https://gogym-m40u.onrender.com' : request.original_url,
        image: image_url('og_image.png'),
        locale: 'ja-JP'
      },
      # Twitter用の設定を個別で設定する
      twitter: {
        card: 'summary_large_image', # Twitterで表示する場合は大きいカードにする
        site: '@your_twitter_handle', # アプリの公式Twitterアカウントがあれば、アカウント名を書く
        image: image_url('og_image.png') # 配置するパスやファイル名によって変更すること
      }
    }
  end
end
