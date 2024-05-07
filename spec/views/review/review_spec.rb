require 'rails_helper'

RSpec.describe '口コミ', type: :system do
  let(:me) { create(:user) }
  let(:gym) { create(:gym) }
  let!(:review_by_me) { create(:review, user: me, gym: gym) }
  let!(:review_by_others) { create(:review, gym: gym) }

  describe '口コミのCRUD' do
    before do
      login_as(me)
      click_on('掲示板')
      click_on('掲示板一覧')
      within "#gym-id-#{gym.id}" do
        page.find_link(gym.title, exact_text: true).click
      end
    end
    describe '口コミの一覧' do
      it '口コミの一覧が表示されること' do
        within '#table-review' do
          expect(page).to have_content(review_by_me.body), '口コミの本文が表示されていません'
          expect(page).to have_content(review_by_me.user.decorate.full_name), '口コミの投稿者のフルネームが表示されていません'
        end
      end
    end

    describe '口コミの作成' do
      it '口コミを作成できること' do
        fill_in '口コミ', with: '新規口コミ'
        click_on '投稿'
        sleep(0.5)
        review = review.last
        within "#review-#{review.id}" do
          expect(page).to have_content(me.decorate.full_name), '新規作成した口コミの投稿者のフルネームが表示されていません'
          expect(page).to have_content('新規口コミ'), '新規作成した口コミの本文が表示されていません'
        end
      end
      it '口コミの作成に失敗すること' do
      expect {
        fill_in '口コミ', with: ''
        click_on '投稿'
        sleep(0.5)
      }.to change { review.count }.by(0), '口コミが作成されています'
      end
    end

    # describe '口コミの削除' do
    #   it '口コミを削除できること' do
    #     within("#review-#{review_by_me.id}") do
    #       page.accept_confirm { find('.delete-review-link').click }
    #     end
    #     expect(page).not_to have_content(review_by_me.body), '口コミの削除が正しく機能していません'
    #   end
    # end

    # describe '口コミの編集' do
    #   context '他人の口コミの場合' do
    #     it '編集ボタン・削除ボタンが表示されないこと' do
    #       within "#review-#{review_by_others.id}" do
    #         expect(page).not_to have_selector('.edit-review-button'), '他人の口コミに対して編集ボタンが表示されてしまっています'
    #         expect(page).not_to have_selector('.delete-review-button'), '他人の口コミに対して削除ボタンが表示されてしまっています'
    #       end
    #     end
    #   end
    # end
  end
end
