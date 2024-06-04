# 環境変数 RUBY_VERSION を定義
ARG RUBY_VERSION=ruby:3.2.2

# 環境変数 NODE_VERSION を定義
ARG NODE_VERSION=18

FROM $RUBY_VERSION
ARG RUBY_VERSION
ARG NODE_VERSION
ENV LANG C.UTF-8
ENV TZ Asia/Tokyo

# 必要なパッケージのインストール
RUN curl -sL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash - \
&& wget --quiet -O - /tmp/pubkey.gpg https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - \
&& echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list \
&& apt-get update -qq \
&& apt-get install -y build-essential nodejs yarn libvips libpq-dev python3 python3-pip python3-venv \
&& apt-get clean \
&& rm -rf /var/lib/apt/lists/*

# Python仮想環境の作成
RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# コンテナ内に作業ディレクトリ GoGym を作成
# RUN mkdir /GoGym

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

ENTRYPOINT ["./entrypoints.sh"]

EXPOSE 3000

# ポート3000でアプリケーションを起動
CMD ["rails", "server", "-b", "0.0.0.0", "-p", "3000"]
