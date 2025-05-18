FROM node:20-alpine AS build

WORKDIR /app

# 依存関係をコピーしてインストール
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN yarn build

# 本番環境ステージ
FROM node:20-alpine AS production

WORKDIR /app

# ビルド結果とパッケージ設定をコピー
COPY --from=build /app/.output /app/.output
COPY --from=build /app/package.json /app/package.json

# アプリケーションを実行
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]
