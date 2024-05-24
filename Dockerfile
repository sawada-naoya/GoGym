# 環境変数 RUBY_VERSION を定義
ARG RUBY_VERSION=ruby:3.2.2

# 環境変数 NODE_VERSION を定義
ARG NODE_VERSION=18

FROM $RUBY_VERSION
ARG RUBY_VERSION
ARG NODE_VERSION
ENV LANG C.UTF-8
ENV TZ Asia/Tokyo
ENV RAILS_ENV=production

# 必要なパッケージのインストール
RUN curl -sL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash - \
&& wget --quiet -O - /tmp/pubkey.gpg https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - \
&& echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list \
&& apt-get update -qq \
&& apt-get install -y build-essential nodejs yarn libvips libpq-dev

# コンテナ内に作業ディレクトリ GoGym を作成
RUN mkdir /GoGym

# 作業ディレクトリを設定
WORKDIR /GoGym

# bundlerのインストール
RUN gem install bundler:2.3.17

# GemfileとGemfile.lock、yarn.lockをコピーしてbundle installを実行
COPY Gemfile /GoGym/Gemfile
COPY Gemfile.lock /GoGym/Gemfile.lock
COPY yarn.lock /GoGym/yarn.lock
RUN bundle install

# yarnをインストール
RUN yarn install

# ローカルのGoGym配下のファイルをコンテナ内のGoGym配下にコピー
COPY . /GoGym

# アセットのプリコンパイル
RUN bundle exec rails assets:precompile

# プリコンパイルとマイグレーションのためのエントリポイントスクリプトを追加
COPY entrypoint.sh /usr/bin/
RUN chmod +x /usr/bin/entrypoint.sh

# CMDをエントリーポイントスクリプトに変更
ENTRYPOINT ["entrypoint.sh"]

# ポート3001でアプリケーションを起動
CMD ["rails", "server", "-b", "0.0.0.0", "-p", "3001"]
