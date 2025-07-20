.PHONY: all build install clean todo-wasm host-wasm compose

# Default target
all: build

# Install required tools
install:
	cargo install wac-cli wit-deps-cli wkg
	cargo install wasm-tools --locked
	@echo "Please install TinyGo from: https://tinygo.org/getting-started/install/macos/"

# Build todo module WASM
todo-wasm: modules/todo/todo.wasm

modules/todo/todo.wasm: modules/todo/*.go modules/todo/wit/todo.wit
	cd modules/todo && tinygo build --target=wasip2 -o todo.wasm --wit-package bindings.wasm --wit-world api .

# Update WIT dependencies and build host
host-wasm: host/target/wasm32-wasip2/release/wasi_http_server.wasm

host/target/wasm32-wasip2/release/wasi_http_server.wasm: host/src/*.rs host/wit/*.wit host/Cargo.toml
	cd host && wit-deps update
	cd host && cargo component build --release --target wasm32-wasip2

# Compose the final WASM
compose: composed.wasm

composed.wasm: modules/todo/todo.wasm host/target/wasm32-wasip2/release/wasi_http_server.wasm
	wac plug host/target/wasm32-wasip2/release/wasi_http_server.wasm --plug modules/todo/todo.wasm -o composed.wasm

# Full build process
build: composed.wasm

# Clean build artifacts
clean:
	rm -f composed.wasm
	rm -f modules/todo/todo.wasm
	cd host && cargo clean
