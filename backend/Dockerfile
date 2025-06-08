# Goのバージョンを1.24に更新
FROM golang:1.24-alpine AS builder

WORKDIR /app

# システムの依存関係をインストール
RUN apk add --no-cache gcc musl-dev

# air をインストール
RUN go install github.com/cosmtrek/air@latest

# 依存関係をコピーしてインストール
COPY go.mod go.sum* ./
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド (開発時はairが担当するため、本番ビルドはコメントアウトまたは削除も可)
# RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o tickethub-api .

# 最終的な軽量イメージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# air とビルド済みバイナリをコピー
COPY --from=builder /go/bin/air /usr/local/bin/
# COPY --from=builder /app/tickethub-api . # 開発時はairがビルドするので不要になる場合がある

# .air.toml をコピー
COPY .air.toml .

# 必要なディレクトリの作成
RUN mkdir -p /app/data
RUN mkdir -p /app/tmp # air が使用するtmpディレクトリ

# 環境変数の設定
ENV GIN_MODE=debug
ENV TZ=Asia/Tokyo

EXPOSE 8080

# air を実行
CMD ["air", "-c", ".air.toml"]
