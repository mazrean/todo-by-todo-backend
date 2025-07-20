#[allow(warnings)]
mod bindings;

use bindings::exports::wasi::http::incoming_handler::Guest;
use bindings::wasi::http::types::*;
use serde::{Deserialize, Serialize};
use serde_json;

wit_bindgen::generate!({
    world: "host",
    generate_all,
});

#[derive(Serialize, Deserialize)]
struct SerializableTodo {
    id: u64,
    user_id: u64,
    title: String,
    description: Option<String>,
    completed: Option<bool>,
    created_at: Option<String>,
    updated_at: Option<String>,
}

#[derive(Serialize, Deserialize)]
struct SerializableTodoRequest {
    user_id: u64,
    title: String,
    description: Option<String>,
    completed: bool,
}

impl From<&Todo> for SerializableTodo {
    fn from(todo: &Todo) -> Self {
        Self {
            id: todo.id,
            user_id: todo.user_id,
            title: todo.title.clone(),
            description: todo.description.clone(),
            completed: todo.completed,
            created_at: todo.created_at.clone(),
            updated_at: todo.updated_at.clone(),
        }
    }
}

impl From<SerializableTodoRequest> for TodoRequest {
    fn from(req: SerializableTodoRequest) -> Self {
        Self {
            user_id: req.user_id,
            title: req.title,
            description: req.description,
            completed: req.completed,
        }
    }
}

// todo APIの関数をimport
use crate::todo::api::todo_api::{
    ApiError, Todo, TodoRequest, create_todo, delete_todo, list_todos, update_todo,
};

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
        // serde_jsonを使用してJSONに変換
        let serializable_todos: Vec<SerializableTodo> = todos.iter().map(|t| t.into()).collect();
        match serde_json::to_string(&serializable_todos) {
            Ok(json) => json,
            Err(_) => "[]".to_string(), // エラー時は空配列を返す
        }
    }

    fn parse_request_body(request: &IncomingRequest) -> Result<TodoRequest, String> {
        // リクエストボディを取得
        let body = request.consume().map_err(|_| "Failed to get request body")?;
        let input_stream = body.stream().map_err(|_| "Failed to get input stream")?;
        
        // ボディのデータを読み取り
        let mut body_bytes = Vec::new();
        loop {
            match input_stream.read(8192) {
                Ok(chunk) => {
                    if chunk.is_empty() {
                        break;
                    }
                    body_bytes.extend_from_slice(&chunk);
                }
                Err(_) => break,
            }
        }
        
        // バイト列を文字列に変換
        let body_str = std::str::from_utf8(&body_bytes)
            .map_err(|_| "Invalid UTF-8 in request body")?;
        
        // serde_jsonを使用してJSONをパース
        serde_json::from_str::<SerializableTodoRequest>(body_str)
            .map(|req| req.into())
            .map_err(|e| format!("JSON parse error: {}", e))
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
                match Self::parse_request_body(&request) {
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
                        match Self::parse_request_body(&request) {
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
            // (Method::Post, "/api/users") => {
            //     // リクエストボディから UserRequest をパース
            //     match Self::parse_request_body_user(&request) {
            //         Ok(user_request) => {
            //             // create_user は Option<ApiError> を返す想定
            //             if let Some(err) = create_user(&user_request) {
            //                 Self::handle_api_error(err)
            //             } else {
            //                 (201, "User created successfully".to_string())
            //             }
            //         }
            //         Err(err) => (400, format!("Invalid request body: {}", err)),
            //     }
            // }
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