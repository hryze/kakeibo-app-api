# kakeibo-app-api

## Documentation for API Endpoints

### user-rest-service

| HTTP request | HTTP response status code | Description |
| :--- | :--- | :--- |
| **POST** /signup | 201 Created<br>400 Bad Request<br>409 Conflict<br>500 Internal Server Error | ユーザー新規登録 |
| **POST** /login | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | ユーザーログイン |
| **DELETE** /logout | 200 OK<br>400 Bad Request<br>500 Internal Server Error | ユーザーログアウト |
| **GET** /groups | 200 OK<br>401 Unauthorized<br>500 Internal Server Error | 承認グループ, 未承認グループ一覧取得 |
| **POST** /groups | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ作成 |
| **PUT** /groups/{group_id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ名更新 |
| **POST** /groups/{group_id}/users | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error | グループ招待 |
| **DELETE** /groups/{group_id}/users | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ退会 |
| **POST** /groups/{group_id}/users/approved | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ招待承認 |
| **DELETE** /groups/{group_id}/users/unapproved | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ招待拒否 |
| **GET** /groups/{group_id}/users/{user_id} | 200 OK<br>400 Bad Request<br>500 Internal Server Error | グループ所属確認 |

### account-rest-service

| HTTP request | HTTP response status code | Description |
| :--- | :--- | :--- |
| **GET** /categories | 200 OK<br>401 Unauthorized<br>500 Internal Server Error | カテゴリー一覧取得 |
| **POST** /categories/custom-categories | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error | カスタムカテゴリー追加 |
| **PUT** /categories/custom-categories/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error | カスタムカテゴリー更新 |
| **DELETE** /categories/custom-categories/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | カスタムカテゴリー削除 |
| **GET** /transactions/{year_month} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 月別家計簿トランザクション一覧取得 |
| **POST** /transactions | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 家計簿トランザクション追加 |
| **PUT** /transactions/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 家計簿トランザクション更新 |
| **DELETE** /transactions/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 家計簿トランザクション削除 |
| **GET** /transactions/search | 200 OK<br>401 Unauthorized<br>500 Internal Server Error | 家計簿トランザクション検索 |
| **POST** /standard-budgets | 201 Created<br>500 Internal Server Error | 家計簿標準予算初期値追加 |
| **GET** /standard-budgets | 200 OK<br>401 Unauthorized<br>500 Internal Server Error | 家計簿標準予算取得 |
| **PUT** /standard-budgets | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 家計簿標準予算更新 |
| **GET** /custom-budgets/{year_month} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 月別家計簿カスタム予算取得 |
| **POST** /custom-budgets/{year_month} | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 月別家計簿カスタム予算追加 |
| **PUT** /custom-budgets/{year_month} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 月別家計簿カスタム予算更新 |
| **DELETE** /custom-budgets/{year_month} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 月別家計簿カスタム予算削除 |
| **GET** /budgets/{year} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 年別家計簿予算一覧取得 |
| **GET** /groups/{group_id}/categories | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループカテゴリー一覧取得 |
| **POST** /groups/{group_id}/categories/custom-categories | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error | グループカスタムカテゴリー追加 |
| **PUT** /groups/{group_id}/categories/custom-categories/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error | グループカスタムカテゴリー更新 |
| **DELETE** /groups/{group_id}/categories/custom-categories/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループカスタムカテゴリー削除 |
| **GET** /groups/{group_id}/transactions/{year_month} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 月別グループ家計簿トランザクション一覧取得 |
| **POST** /groups/{group_id}/transactions | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ家計簿トランザクション追加 |
| **PUT** /groups/{group_id}/transactions/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ家計簿トランザクション更新 |
| **DELETE** /groups/{group_id}/transactions/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ家計簿トランザクション削除 |
| **GET** /groups/{group_id}/transactions/search | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ家計簿トランザクション検索 |
| **GET** /groups/{group_id}/transactions/{year_month}/account | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 月別グループ家計簿トランザクション会計データ取得 |
| **POST** /groups/{group_id}/transactions/{year_month}/account | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 月別グループ家計簿トランザクション会計データ追加 |
| **PUT** /groups/{group_id}/transactions/{year_month}/account | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 月別グループ家計簿トランザクション会計データ更新 |
| **DELETE** /groups/{group_id}/transactions/{year_month}/account | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 月別グループ家計簿トランザクション会計データ削除 |
| **POST** /groups/{group_id}/standard-budgets | 201 Created<br>400 Bad Request<br>500 Internal Server Error | グループ家計簿標準予算初期値追加 |
| **GET** /groups/{group_id}/standard-budgets | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ家計簿標準予算取得 |
| **PUT** /groups/{group_id}/standard-budgets | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ家計簿標準予算更新 |
| **GET** /groups/{group_id}/custom-budgets/{year_month} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ月別家計簿カスタム予算取得 |
| **POST** /groups/{group_id}/custom-budgets/{year_month} | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ月別家計簿カスタム予算追加 |
| **PUT** /groups/{group_id}/custom-budgets/{year_month} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ月別家計簿カスタム予算更新 |
| **DELETE** /groups/{group_id}/custom-budgets/{year_month} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ月別家計簿カスタム予算削除 |
| **GET** /groups/{group_id}/budgets/{year} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ年別家計簿予算一覧取得 |

### todo-rest-service

| HTTP request | HTTP response status code | Description |
| :--- | :--- | :--- |
| **GET** /todo-list/{date} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 日別実施予定todo, 締切予定todo一覧取得 |
| **GET** /todo-list/{year_month} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | 月別実施予定todo, 締切予定todo一覧取得 |
| **POST** /todo-list | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | todo追加 |
| **PUT** /todo-list/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | todo更新 |
| **DELETE** /todo-list/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | todo削除 |
| **GET** /todo-list/search | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | todo検索 |
| **GET** /groups/{group_id}/todo-list/{date} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ日別実施予定todo、グループ締切予定todo一覧取得 |
| **GET** /groups/{group_id}/todo-list/{year_month} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループ月別実施予定todo、グループ締切予定todo一覧取得 |
| **POST** /groups/{group_id}/todo-list | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループtodo追加 |
| **PUT** /groups/{group_id}/todo-list/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループtodo更新 |
| **DELETE** /groups/{group_id}/todo-list/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループtodo削除 |
| **GET** /groups/{group_id}/todo-list/search | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループtodo検索 |
| **GET** /groups/{group_id}/tasks/users | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | ユーザー別グループタスク一覧取得 |
| **POST** /groups/{group_id}/tasks/users | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error | グループタスクユーザー追加 |
| **GET** /groups/{group_id}/tasks | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループタスク一覧取得 |
| **POST** /groups/{group_id}/tasks | 201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループタスク追加 |
| **PUT** /groups/{group_id}/tasks/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループタスク更新 |
| **DELETE** /groups/{group_id}/tasks/{id} | 200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error | グループタスク削除 |
