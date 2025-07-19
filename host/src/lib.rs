#[allow(warnings)]
mod bindings;

use bindings::exports::wasi::http::incoming_handler::Guest;
use bindings::wasi::http::types::*;

struct Component;

impl Guest for Component {
    fn handle(request: IncomingRequest, response_out: ResponseOutparam) {
        // リクエストのパスを取得
        let path = request.path_with_query().unwrap_or_else(|| "/".to_string());

        // リクエストメソッドを取得
        let method = request.method();

        // ヘッダーの作成
        let headers = Headers::new();

        // レスポンスの作成
        let response = OutgoingResponse::new(headers);

        // ステータスコードとメッセージの設定
        let (status_code, message) = match (method, path.as_str()) {
            (Method::Get, "/") => (200, "Welcome to WASI HTTP Server!"),
            (Method::Get, "/hello") => (200, "Hello, World from WASM!"),
            (Method::Get, "/api/status") => (200, r#"{"status": "ok", "version": "1.0"}"#),
            _ => (404, "Not Found"),
        };

        response.set_status_code(status_code).unwrap();

        // レスポンスボディの作成
        let body = response.body().unwrap();
        let stream = body.write().unwrap();

        // データを書き込む
        stream.write(message.as_bytes()).unwrap();
        stream.flush().unwrap();
        drop(stream);

        // ボディを完了
        OutgoingBody::finish(body, None).unwrap();

        // レスポンスを送信（関連関数として呼び出す）
        ResponseOutparam::set(response_out, Ok(response));
    }
}

bindings::export!(Component with_types_in bindings);
