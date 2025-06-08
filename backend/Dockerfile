FROM golang:1.23-alpine AS builder

WORKDIR /app

# システムの依存関係をインストール
RUN apk add --no-cache gcc musl-dev

# 依存関係をコピーしてインストール
COPY go.mod go.sum* ./
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o tickethub-api .

# 最終的な軽量イメージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# バイナリをコピー
COPY --from=builder /app/tickethub-api .

# 必要なディレクトリの作成
RUN mkdir -p /app/data

# 環境変数の設定
ENV GIN_MODE=release
ENV TZ=Asia/Tokyo

EXPOSE 8080

CMD ["./tickethub-api"]
