require 'rails_helper'

RSpec.feature "ReviewPosting", type: :feature do
  let(:user) { create(:user) }
  let(:gym) { create(:gym) }

  before do
    # ユーザーログインの処理をここに書く
    visit login_path
    fill_in 'Email', with: user.email
    fill_in 'Password', with: 'password'
    click_button 'ログイン'
    visit gym_path(gym)
    click_on '口コミ投稿'
  end

  scenario "口コミの作成に成功する" do
    fill_in 'Title', with: '素晴らしいジム!'
    fill_in 'Content', with: 'とても清潔で設備も充実しています。'
    attach_file 'Image', "#{Rails.root}/spec/fixtures/test_image.jpg" 
    click_button '投稿'

    expect(page).to have_content '口コミを投稿しました'
    expect(page).to have_current_path gym_reviews_path(gym)
    # 他にも成功時の確認が必要ならここに追加
  end

  scenario "口コミの作成に失敗する" do
    fill_in 'Title', with: '' # タイトルを空にして失敗条件を作る
    fill_in 'Content', with: '' # 内容も空に
    click_button '投稿'

    expect(page).to have_content '口コミの投稿に失敗しました'
    expect(page).to have_current_path gym_path(gym), status: :unprocessable_entity
  end
end
