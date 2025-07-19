# 📘 Todo API Documentation

Base URL: http://localhost:8080

### 🔍 Health Check

```sh
curl http://localhost:8080/health
```

GET /health
ヘルスチェック用のエンドポイント。

Response

- 200 OK: サーバーが稼働中であることを示します。

### 📝 Todos

#### 全ての Todo を取得

```sh
curl http://localhost:8080/todos
```

GET /todos

Response

- 200 OK:

  ```json

  [
  {
  "id": 1,
  "user_id": 1,
  "title": "買い物に行く",
  "description": "牛乳を買う",
  "completed": false
  },
  ...
  ]
  ```

#### 新しい Todo を作成

```sh
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "title": "宿題をやる",
    "description": "数学と英語を終わらせる",
    "completed": false
  }'
```

POST /todos

Request Body

```json
{
  "user_id": 1,
  "title": "宿題をやる",
  "description": "数学と英語",
  "completed": false
}
```

Response

- 201 Created: 空レスポンス（成功）
- 400 Bad Request: 不正な JSON
- 500 Internal Server Error: DB エラーなど

#### 既存の Todo を更新

```sh
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "title": "宿題をやった",
    "description": "理科も追加",
    "completed": true
  }'
```

PUT /todos/{id}

Path Parameter
id：更新対象の Todo ID

Request Body

```json
{
  "user_id": 1,
  "title": "宿題を終わらせる",
  "description": "理科を追加",
  "completed": true
}
```

Response

- 200 OK: 正常に更新された場合
- 400 Bad Request: ID またはリクエスト形式が不正
- 500 Internal Server Error: DB 更新失敗など

#### 指定された Todo を削除

```sh
curl -X DELETE http://localhost:8080/todos/1
```

DELETE /todos/{id}

Path Parameter
id：削除対象の Todo ID

Response

- 200 OK: 正常に削除
- 400 Bad Request: ID 不正
- 500 Internal Server Error: DB エラーなど

### 👤 Users

#### POST /users/{id}/todos

```sh

```

新しいユーザーを作成します。
（※/users/{id}/todos の URL ですが、実際の実装内容は ユーザー作成 です。設計ミスの可能性があるため後述します）

Request Body

```json
{
  "name": "Taro Yamada"
}
```

Response

- 201 Created:
  ```json
  {
    "user_id": 1
  }
  ```
- 400 Bad Request: 不正なリクエストボディ
- 500 Internal Server Error: DB エラーなど
