# Rust build stage
FROM rust:1.88 AS rust-builder

WORKDIR /app

# Add wasm32-wasip2 target
RUN rustup target add wasm32-wasip2

# Install WIT tools
RUN cargo install wit-deps-cli cargo-component wasm-tools --locked

# Copy Rust source code and required dependencies
COPY host/ ./host/
COPY modules/todo/wit/ ./modules/todo/wit/

# Build Rust WASM
RUN cd host && wit-deps update && \
    cargo component build --release --target wasm32-wasip2

# Go build stage  
FROM tinygo/tinygo:0.38.0 AS go-builder

WORKDIR /app

# Copy pre-built wasm-tools from a Rust stage
COPY --from=rust-builder /usr/local/cargo/bin/wasm-tools /usr/local/bin/wasm-tools

# Copy Go source code
COPY modules/ ./modules/

# Build Go WASM
RUN cd modules/todo && \
    tinygo build --target=wasip2 -o /app/todo.wasm --wit-package bindings.wasm --wit-world api .

# Composition stage
FROM rust:1.88 AS composer

WORKDIR /app

# Install wac-cli
RUN cargo install wac-cli --locked

# Copy built WASM files from previous stages
COPY --from=rust-builder /app/host/target/wasm32-wasip2/release/wasi_http_server.wasm ./host/target/wasm32-wasip2/release/
COPY --from=go-builder /app/todo.wasm ./modules/todo/todo.wasm

# Compose the final WASM
RUN wac plug host/target/wasm32-wasip2/release/wasi_http_server.wasm --plug modules/todo/todo.wasm -o composed.wasm

# Runtime stage
FROM ubuntu:22.04 AS runtime

ENV DEBIAN_FRONTEND=noninteractive
ENV HOME=/root

WORKDIR /app

ENV DB_FILEPATH=/app/store/db.json

# Install minimal runtime dependencies
RUN apt-get update && apt-get install -y \
    curl \
    ca-certificates \
    xz-utils \
    && rm -rf /var/lib/apt/lists/*

# Install wasmtime only
RUN curl https://wasmtime.dev/install.sh -sSf | bash \
    && ln -sf /root/.wasmtime/bin/wasmtime /usr/local/bin/wasmtime

# Copy the built WASM file from composer stage
COPY --from=composer /app/composed.wasm .

EXPOSE 8080

# Run the composed WASM
CMD ["wasmtime", "serve", "--wasi", "cli,inherit-env", "--dir", "/app/store", "composed.wasm"]
