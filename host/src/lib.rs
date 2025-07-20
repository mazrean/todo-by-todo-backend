#[allow(warnings)]
mod bindings;

use bindings::exports::wasi::http::incoming_handler::Guest;
use bindings::wasi::http::types::*;
use bindings::wasi::io::streams::InputStream;

wit_bindgen::generate!({
    world: "host",
    generate_all,
});

// todo APIの関数をimport
use crate::todo::api::todo_api::{list_todos, create_todo, update_todo, delete_todo, create_user, ApiError, Todo, TodoRequest, UserRequest};

struct Component;

impl Component {
    fn handle_api_error(error: ApiError) -> (u16, String) {
        match error {
            ApiError::InternalError(msg) => (500, format!("Internal Server Error: {}", msg)),
            ApiError::InvalidId(id) => (400, format!("Invalid ID: {}", id)),
            ApiError::InvalidRequest(msg) => (400, format!("Invalid Request: {}", msg)),
            ApiError::NotFound(id) => (404, format!("Todo Not Found: {}", id)),
        }
    }

    fn todos_to_json(todos: &[Todo]) -> String {
        // 簡単なJSON形式に変換（serde無しで）
        let mut json = String::from("[");
        for (i, todo) in todos.iter().enumerate() {
            if i > 0 {
                json.push(',');
            }
            json.push_str(&format!(
                r#"{{"id":{},"user_id":{},"title":"{}","description":{},"completed":{},"created_at":{},"updated_at":{}}}"#,
                todo.id,
                todo.user_id,
                todo.title.replace('"', r#"\""#),
                todo.description.as_ref().map_or("null".to_string(), |d| format!(r#""{}""#, d.replace('"', r#"\""#))),
                todo.completed.unwrap_or(false),
                todo.created_at.as_ref().map_or("null".to_string(), |d| format!(r#""{}""#, d.replace('"', r#"\""#))),
                todo.updated_at.as_ref().map_or("null".to_string(), |d| format!(r#""{}""#, d.replace('"', r#"\""#)))
            ));
        }
        json.push(']');
        json
    }

    fn parse_request_body_todo(_request: &IncomingRequest) -> Result<TodoRequest, String> {
        // リクエストボディの読み取り（簡易版）
        // 実際の実装では、リクエストボディを読み取ってJSONをパースする必要がある
        // ここでは例として固定値を返す
        Ok(TodoRequest {
            user_id: 1,
            title: "Parsed Todo".to_string(),
            description: Some("Parsed from request body".to_string()),
            completed: false,
        })
    }

    fn parse_request_body_user(_request: &IncomingRequest) -> Result<UserRequest, String> {
        // リクエストボディの読み取り（簡易版）
        // 実際の実装では、リクエストボディを読み取ってJSONをパースする必要がある
        // ここでは例として固定値を返す
        Ok(UserRequest {
            name: "Parsed User".to_string(),
        })
    }
}

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
            (Method::Get, "/api/todos") => {
                // Todoリストの取得
                let todos_result = list_todos();
                if let Some(error) = todos_result.error {
                    Self::handle_api_error(error)
                } else {
                    // JSONシリアライゼーションは簡単な形式で行う
                    let todos_json = Self::todos_to_json(&todos_result.todos);
                    (200, todos_json)
                }
            }
            (Method::Post, "/api/todos") => {
                // リクエストボディからTodoRequestを解析
                match Self::parse_request_body_todo(&request) {
                    Ok(todo_request) => {
                        let create_result = create_todo(&todo_request);
                        if let Some(error) = create_result.error {
                            Self::handle_api_error(error)
                        } else {
                            (201, "Todo created successfully".to_string())
                        }
                    }
                    Err(err) => (400, format!("Invalid request body: {}", err)),
                }
            }
            (Method::Put, path) if path.starts_with("/api/todos/") => {
                // パスからIDを抽出
                if let Some(id_str) = path.strip_prefix("/todos/") {
                    if let Ok(id) = id_str.parse::<u64>() {
                        match Self::parse_request_body_todo(&request) {
                            Ok(todo_request) => {
                                if let Some(error) = update_todo(id, &todo_request) {
                                    Self::handle_api_error(error)
                                } else {
                                    (200, "Todo updated successfully".to_string())
                                }
                            }
                            Err(err) => (400, format!("Invalid request body: {}", err)),
                        }
                    } else {
                        (400, "Invalid ID".to_string())
                    }
                } else {
                    (400, "Invalid path".to_string())
                }
            }
            (Method::Delete, path) if path.starts_with("/api/todos/") => {
                // パスからIDを抽出
                if let Some(id_str) = path.strip_prefix("/todos/") {
                    if let Ok(id) = id_str.parse::<u64>() {
                        if let Some(error) = delete_todo(id) {
                            Self::handle_api_error(error)
                        } else {
                            (200, "Todo deleted successfully".to_string())
                        }
                    } else {
                        (400, "Invalid ID".to_string())
                    }
                } else {
                    (400, "Invalid path".to_string())
                }
            }
            (Method::Post, "/api/users") => {
                // リクエストボディから UserRequest をパース
                match Self::parse_request_body_user(&request) {
                    Ok(user_request) => {
                        // create_user は Option<ApiError> を返す想定
                        if let Some(err) = create_user(&user_request) {
                            Self::handle_api_error(err)
                        } else {
                            (201, "User created successfully".to_string())
                        }
                    }
                    Err(err) => (400, format!("Invalid request body: {}", err)),
                }
            }
            (Method::Get, "/health") => (200, "OK".to_string()),
            _ => (404, "Not Found".to_string()),
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
