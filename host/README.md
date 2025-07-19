# WASM HTTP サーバー

## 環境構築

### プロジェクトをビルド

```sh
cargo component build --release
```

成功すると、target/wasm32-wasip1/release/wasi_http_server.wasm が生成される

### サーバーを起動

```sh
wasmtime serve \
  -S cli,inherit-env \
  --addr 127.0.0.1:3000 \
  target/wasm32-wasip1/release/wasi_http_server.wasm
```

ポートが埋まっている場合

```sh
lsof -i :3000
kill -9 <PID>
```
