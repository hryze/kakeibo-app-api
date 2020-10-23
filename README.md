# kakeibo-app-api

## Documentation for API Endpoints

### user-rest-service

| HTTP request | HTTP response<br>status code | Description |
| :--- | :--- | :--- |
| **POST**<br>&emsp;/signup | <pre>201 Created<br>400 Bad Request<br>409 Conflict<br>500 Internal Server Error</pre> | ユーザー新規登録 |
| **POST**<br>&emsp;/login | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | ユーザーログイン |
| **DELETE**<br>&emsp;/logout | <pre>200 OK<br>400 Bad Request<br>500 Internal Server Error</pre> | ユーザーログアウト |
| **GET**<br>&emsp;/groups | <pre>200 OK<br>401 Unauthorized<br>500 Internal Server Error</pre> | 承認グループ,<br>未承認グループ一覧取得 |
| **POST**<br>&emsp;/groups | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ作成 |
| **PUT**<br>&emsp;/groups/{group_id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ名更新 |
| **POST**<br>&emsp;/groups/{group_id}/users | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error</pre> | グループ招待 |
| **DELETE**<br>&emsp;/groups/{group_id}/users | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ退会 |
| **POST**<br>&emsp;/groups/{group_id}/users/approved | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ招待承認 |
| **DELETE**<br>&emsp;/groups/{group_id}/users/unapproved | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ招待拒否 |
| **GET**<br>&emsp;/groups/{group_id}/users/{user_id} | <pre>200 OK<br>400 Bad Request<br>500 Internal Server Error</pre> | グループ所属確認 |

### account-rest-service

| HTTP request | HTTP response<br>status code | Description |
| :--- | :--- | :--- |
| **GET**<br>&emsp;/categories | <pre>200 OK<br>401 Unauthorized<br>500 Internal Server Error</pre> | カテゴリー一覧取得 |
| **POST**<br>&emsp;/categories/custom-categories | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error</pre> | カスタムカテゴリー追加 |
| **PUT**<br>&emsp;/categories/custom-categories/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error</pre> | カスタムカテゴリー更新 |
| **DELETE**<br>&emsp;/categories/custom-categories/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | カスタムカテゴリー削除 |
| **GET**<br>&emsp;/transactions/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別家計簿トランザクション一覧取得 |
| **POST**<br>&emsp;/transactions | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 家計簿トランザクション追加 |
| **PUT**<br>&emsp;/transactions/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 家計簿トランザクション更新 |
| **DELETE**<br>&emsp;/transactions/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 家計簿トランザクション削除 |
| **GET**<br>&emsp;/transactions/search | <pre>200 OK<br>401 Unauthorized<br>500 Internal Server Error</pre> | 家計簿トランザクション検索 |
| **POST**<br>&emsp;/standard-budgets | <pre>201 Created<br>500 Internal Server Error</pre> | 家計簿標準予算初期値追加 |
| **GET**<br>&emsp;/standard-budgets | <pre>200 OK<br>401 Unauthorized<br>500 Internal Server Error</pre> | 家計簿標準予算取得 |
| **PUT**<br>&emsp;/standard-budgets | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 家計簿標準予算更新 |
| **GET**<br>&emsp;/custom-budgets/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別家計簿カスタム予算取得 |
| **POST**<br>&emsp;/custom-budgets/{year_month} | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別家計簿カスタム予算追加 |
| **PUT**<br>&emsp;/custom-budgets/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別家計簿カスタム予算更新 |
| **DELETE**<br>&emsp;/custom-budgets/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別家計簿カスタム予算削除 |
| **GET**<br>&emsp;/budgets/{year} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 年別家計簿予算一覧取得 |
| **GET**<br>&emsp;/groups/{group_id}/categories | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループカテゴリー一覧取得 |
| **POST**<br>&emsp;/groups/{group_id}/categories/custom-categories | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error</pre> | グループカスタムカテゴリー追加 |
| **PUT**<br>&emsp;/groups/{group_id}/categories/custom-categories/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error</pre> | グループカスタムカテゴリー更新 |
| **DELETE**<br>&emsp;/groups/{group_id}/categories/custom-categories/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループカスタムカテゴリー削除 |
| **GET**<br>&emsp;/groups/{group_id}/transactions/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別グループ家計簿トランザクション一覧取得 |
| **POST**<br>&emsp;/groups/{group_id}/transactions | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿トランザクション追加 |
| **PUT**<br>&emsp;/groups/{group_id}/transactions/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿トランザクション更新 |
| **DELETE**<br>&emsp;/groups/{group_id}/transactions/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿トランザクション削除 |
| **GET**<br>&emsp;/groups/{group_id}/transactions/search | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿トランザクション検索 |
| **GET**<br>&emsp;/groups/{group_id}/transactions/{year_month}/account | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別グループ家計簿トランザクション会計データ取得 |
| **POST**<br>&emsp;/groups/{group_id}/transactions/{year_month}/account | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別グループ家計簿トランザクション会計データ追加 |
| **PUT**<br>&emsp;/groups/{group_id}/transactions/{year_month}/account | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別グループ家計簿トランザクション会計データ更新 |
| **DELETE**<br>&emsp;/groups/{group_id}/transactions/{year_month}/account | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別グループ家計簿トランザクション会計データ削除 |
| **POST**<br>&emsp;/groups/{group_id}/standard-budgets | <pre>201 Created<br>400 Bad Request<br>500 Internal Server Error</pre> | グループ家計簿標準予算初期値追加 |
| **GET**<br>&emsp;/groups/{group_id}/standard-budgets | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿標準予算取得 |
| **PUT**<br>&emsp;/groups/{group_id}/standard-budgets | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿標準予算更新 |
| **GET**<br>&emsp;/groups/{group_id}/custom-budgets/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ月別家計簿カスタム予算取得 |
| **POST**<br>&emsp;/groups/{group_id}/custom-budgets/{year_month} | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ月別家計簿カスタム予算追加 |
| **PUT**<br>&emsp;/groups/{group_id}/custom-budgets/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ月別家計簿カスタム予算更新 |
| **DELETE**<br>&emsp;/groups/{group_id}/custom-budgets/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ月別家計簿カスタム予算削除 |
| **GET**<br>&emsp;/groups/{group_id}/budgets/{year} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ年別家計簿予算一覧取得 |

### todo-rest-service

| HTTP request | HTTP response<br>status code | Description |
| :--- | :--- | :--- |
| **GET**<br>&emsp;/todo-list/{date} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 日別実施予定todo,<br>締切予定todo一覧取得 |
| **GET**<br>&emsp;/todo-list/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別実施予定todo,<br>締切予定todo一覧取得 |
| **GET**<br>&emsp;/todo-list/expired | <pre>200 OK<br>401 Unauthorized<br>500 Internal Server Error</pre> | 期限切れtodo一覧取得 |
| **POST**<br>&emsp;/todo-list | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | todo追加 |
| **PUT**<br>&emsp;/todo-list/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | todo更新 |
| **DELETE**<br>&emsp;/todo-list/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | todo削除 |
| **GET**<br>&emsp;/todo-list/search | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | todo検索 |
| **GET**<br>&emsp;/groups/{group_id}/todo-list/{date} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ日別実施予定todo,<br>グループ締切予定todo一覧取得 |
| **GET**<br>&emsp;/groups/{group_id}/todo-list/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ月別実施予定todo,<br>グループ締切予定todo一覧取得 |
| **GET**<br>&emsp;/groups/{group_id}/todo-list/expired | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ期限切れtodo一覧取得 |
| **POST**<br>&emsp;/groups/{group_id}/todo-list | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループtodo追加 |
| **PUT**<br>&emsp;/groups/{group_id}/todo-list/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループtodo更新 |
| **DELETE**<br>&emsp;/groups/{group_id}/todo-list/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループtodo削除 |
| **GET**<br>&emsp;/groups/{group_id}/todo-list/search | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループtodo検索 |
| **GET**<br>&emsp;/groups/{group_id}/tasks/users | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | ユーザー別グループタスク一覧取得 |
| **POST**<br>&emsp;/groups/{group_id}/tasks/users | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error</pre> | グループタスクユーザー追加 |
| **DELETE**<br>&emsp;/groups/{group_id}/tasks/users | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループタスクユーザー削除 |
| **GET**<br>&emsp;/groups/{group_id}/tasks | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループタスク一覧取得 |
| **POST**<br>&emsp;/groups/{group_id}/tasks | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループタスク追加 |
| **PUT**<br>&emsp;/groups/{group_id}/tasks/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループタスク更新 |
| **DELETE**<br>&emsp;/groups/{group_id}/tasks/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループタスク削除 |
