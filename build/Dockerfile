FROM golang:1.18.1-alpine
# アップデート
RUN apk update && \
    apk --no-cache add git
# appディレクトリの作成
RUN mkdir /go/src/app
# ワーキングディレクトリの設定
WORKDIR /go/src/app/
# ホストのファイルをコンテナの作業ディレクトリに移行
ADD . /go/src/app

RUN go mod tidy
RUN go install github.com/cosmtrek/air@v1.29.0

CMD ["air", "-c", ".air.toml"]