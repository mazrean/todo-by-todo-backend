# host/Dockerfile
FROM ubuntu:22.04

# 非対話モード
ENV DEBIAN_FRONTEND=noninteractive
# 明示的に HOME を設定
ENV HOME=/root

WORKDIR /app

# curl と証明書を入れつつ、install.sh で wasmtime を ~/.wasmtime/bin にインストールし、
# さらに /usr/local/bin にシンボリックリンクを張る
RUN apt-get update \
    && apt-get install -y \
        curl \
        ca-certificates \
        xz-utils \       
    && rm -rf /var/lib/apt/lists/* \
    && curl https://wasmtime.dev/install.sh -sSf | bash \
    && ln -sf /root/.wasmtime/bin/wasmtime /usr/local/bin/wasmtime

# あなたの .wasm をコピー
COPY target/wasm32-wasip1/release/wasi_http_server.wasm .

EXPOSE 3001

# これでどこからでも `wasmtime` が呼べる
CMD ["wasmtime", "serve", "--addr", "0.0.0.0:3001", "-S", "cli,inherit-env", "wasi_http_server.wasm"]
