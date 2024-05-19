require 'rails_helper'

RSpec.feature "ReviewCRUD", type: :feature do
  let(:user) { create(:user) }
  let(:gym) { create(:gym) }
  let!(:review) { create(:review, gym: gym, user: user) }

  before do
    visit login_path
    fill_in 'email', with: user.email
    fill_in 'password', with: 'password'
    click_button 'ログイン'
    visit gym_path(gym)
  end

  describe '口コミのCRUD' do
    describe '口コミの一覧表示' do
      it '口コミ一覧に投稿された口コミが表示されること' do
        click_on '口コミ一覧'

        expect(page).to have_current_path gym_reviews_path(gym)
        expect(page).to have_content(review.user.name)
        expect(page).to have_content(review.title)
        expect(page).to have_content(review.content)
        expect(page).to have_content(review.rating)
      end
    end

    describe '口コミの作成' do
      before do
        click_on '口コミ投稿'
      end

      it '口コミの作成に成功する' do

        fill_in 'review_title', with: '素晴らしいジム!'
        fill_in 'review_content', with: 'とても清潔で設備も充実しています。'
        click_button '投稿'

        expect(page).to have_content '口コミが投稿されました'
        expect(page).to have_current_path gym_reviews_path(gym)
        expect(page).to have_content '素晴らしいジム!'
        expect(page).to have_content 'とても清潔で設備も充実しています。'
      end

      it '口コミの作成に失敗する' do
        click_on '口コミ投稿'

        fill_in 'review_title', with: ''
        fill_in 'review_content', with: ''
        click_button '投稿'

        expect(page).to have_content '口コミを投稿出来ませんでした'
      end
    end

    describe '口コミの編集' do
      before do
        click_on '口コミ一覧'
        within "#review-#{review.id}" do
          click_on '編集'
        end
      end

      it '口コミの編集に成功する' do
        fill_in 'review_title', with: '更新されたジム!'
        fill_in 'review_content', with: 'さらに清潔で設備も増えました。'
        click_button '更新'

        expect(page).to have_content '口コミを編集しました'
        expect(page).to have_current_path gym_reviews_path(gym)
        expect(page).to have_content '更新されたジム!'
        expect(page).to have_content 'さらに清潔で設備も増えました。'
      end

      it '口コミの編集に失敗する' do
        fill_in 'review_title', with: ''
        fill_in 'review_content', with: ''
        click_button '更新'

        expect(page).to have_content '口コミを編集できませんでした'
      end
    end
  end
end
