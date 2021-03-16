# Tukecholl

自身の転職活動用ポートフォリオとして制作した家計簿アプリです。  
Go, React, Kubernetesを用いたマイクロサービスアーキテクチャです。

## WebサイトURL

https://www.shakepiper.com

※ 火・木曜日11:00 ~ 17:00稼働  
※ ログインページのゲストユーザーログインより、「 **郷ひろみ** 」として簡単ログインできます。

## 制作背景

現在友人と4人でシェアハウスしており、家賃、食費、日用品等、全員で費用を出し合って生活しています。

そんな中シェアハウス内での支出を管理したいという事で家計簿アプリを探しましたが、個人に特化した家計簿アプリがほとんどで、グループ機能があったとしても私達の要件を満たせるようなアプリが見つかりませんでした。

このような経緯から、家計簿アプリの定番機能である収支管理、予算管理等は実装した上で、グループ内での月末精算、買い物予定リスト、TODOリスト、料理当番や洗濯当番などのシフト管理等、シェアハウスに特化した機能があると便利だと思い制作しました。

また、シェアハウスしている方以外にも、家族や夫婦、カップル等、同じような不便を感じられてる方々に使用して頂けるサービスを作りたいと考えました。

## 開発形態

**<ins>開発者</ins>**

- 平 侑祐
- 安樂 亮佑（共同開発者）
- 古澤 宏弥（共同開発者）
- 伊藤 稜悟（共同開発者）

**<ins>制作物 / 担当</ins>**

- API / 平 侑祐  
  https://github.com/paypay3/kakeibo-app-api

- Terraform / 平 侑祐  
  https://github.com/paypay3/kakeibo-app-terraform

- Kubernetes / 平 侑祐  
  https://github.com/paypay3/kakeibo-app-kubernetes

- Frontend / 安樂 亮佑（共同開発者）, 古澤 宏弥（共同開発者）, 伊藤 稜悟（共同開発者）  
  https://github.com/ryo-wens/kakeibo-front

**<ins>開発手法</ins>**

- アジャイル開発（スクラム）

**<ins>コミュニケーションツール</ins>**

- Slack
- Trello
- Googleスプレッドシート

## 使用技術

### 【 _Frontend_ 】

**<ins>Language</ins>**

- TypeScript v4.2.2
- Sass

**<ins>Library/Framework</ins>**

- React v17.0.1
- Redux v4.0.5

### 【 _Backend_ 】

**<ins>Language</ins>**

- Go v1.16.2

### 【 _Infrastructure_ 】

**<ins>Cloud Service</ins>**

- AWS

**<ins>Infrastructure as Code</ins>**

- Terraform v0.14.7
    - VPC
    - Subnet
    - Route Table
    - Internet Gateway
    - NAT Gateway
    - Security Group
    - EKS
    - ECR
    - S3
    - CloudFront
    - ELB
    - EC2
    - Route53
    - ACM
    - RDS(MySQL v8.0.20)
    - ElastiCache(Redis v5.0.6)
    - Secrets Manager
    - IAM

**<ins>Container</ins>**

- docker v20.10.2
- docker-compose v1.27.4（開発環境Database用）

**<ins>Container Orchestration</ins>**

- Kubernetes v1.18
    - api × 3
    - argocd
    - aws-load-balancer-controller
    - cert-manager
    - cluster-autoscaler
    - external-dns
    - external-secrets
    - initialize-rds-job
    - metrics-server
- Kustomize v4.0.4

**<ins>CI/CD</ins>**

- GitHub Actions
- ArgoCD

## インフラ構成図

![infra](https://user-images.githubusercontent.com/59386359/101023134-ad6ebb00-35b5-11eb-8fb8-58bc6e64eb5d.png)

## CI/CD pipeline

![ci-cd](https://user-images.githubusercontent.com/59386359/101271226-2e9b9d00-37c4-11eb-8842-8020fb66b11c.png)

## ER図

![kakeibo-er](https://user-images.githubusercontent.com/59386359/103619157-90425900-4f74-11eb-8f1b-b5c45297eb60.png)

## 機能一覧

### 【 個人利用機能 】

**<ins>ユーザー機能</ins>**

- ユーザー新規登録
- ユーザーログイン
- ユーザーログアウト

**<ins>家計簿機能</ins>**

- 月別家計簿取引一覧取得
- 家計簿取引最終更新10件取得
- 家計簿取引追加
- 家計簿取引更新
- 家計簿取引削除
- 家計簿取引検索

**<ins>家計簿予算機能</ins>**

- 標準予算取得
- 標準予算更新
- 月別カスタム予算取得
- 月別カスタム予算追加
- 月別カスタム予算更新
- 月別カスタム予算削除
- 年別予算一覧取得

**<ins>カテゴリー機能</ins>**

- カテゴリー一覧取得
- カスタムカテゴリー追加
- カスタムカテゴリー更新
- カスタムカテゴリー削除

**<ins>todo機能</ins>**

- 日別実施予定todo, 締切予定todo一覧取得
- 月別実施予定todo, 締切予定todo一覧取得
- 期限切れtodo一覧取得
- todo追加
- todo更新
- todo削除
- todo検索

**<ins>買い物リスト機能</ins>**

- 定期買い物リスト, 日別買い物リスト取得
- 定期買い物リスト, 月別買い物リスト取得
- 期限切れ買い物リスト取得
- 買い物リスト更新, 家計簿トランザクション自動追加/自動削除
- 買い物リスト削除
- 定期買い物リスト追加
- 定期買い物リスト更新
- 定期買い物リスト削除

### 【 グループ利用機能 】

**<ins>グループ機能</ins>**

- 承認グループ, 未承認グループ一覧取得
- グループ作成
- グループ名更新
- グループ招待
- グループ招待承認
- グループ招待拒否
- グループ退会

**<ins>グループ家計簿機能</ins>**

- 月別家計簿取引一覧取得
- 家計簿取引最終更新10件取得
- 家計簿取引追加
- 家計簿取引更新
- 家計簿取引削除
- 家計簿取引検索
- 年別家計簿取引会計状況一覧取得
- 月別家計簿取引自動会計
- 月別家計簿取引会計データ取得
- 月別家計簿取引会計データ更新
- 月別家計簿取引会計データ削除

**<ins>グループ予算機能</ins>**

- 標準予算取得
- 標準予算更新
- 月別カスタム予算取得
- 月別カスタム予算追加
- 月別カスタム予算更新
- 月別カスタム予算削除
- 年別予算一覧取得

**<ins>グループカテゴリー機能</ins>**

- カテゴリー一覧取得
- カスタムカテゴリー追加
- カスタムカテゴリー更新
- カスタムカテゴリー削除

**<ins>グループtodo機能</ins>**

- 日別実施予定todo, 締切予定todo一覧取得
- 月別実施予定todo, 締切予定todo一覧取得
- 期限切れtodo一覧取得
- todo追加
- todo更新
- todo削除
- todo検索

**<ins>グループ買い物リスト機能</ins>**

- 定期買い物リスト, 日別買い物リスト取得
- 定期買い物リスト, 月別買い物リスト取得
- 期限切れ買い物リスト取得
- 買い物リスト更新, 家計簿トランザクション自動追加/自動削除
- 買い物リスト削除
- 定期買い物リスト追加
- 定期買い物リスト更新
- 定期買い物リスト削除

**<ins>グループシフト管理機能</ins>**

- ユーザー別シフト一覧取得
- シフト機能用ユーザー追加
- シフト機能用ユーザー削除
- シフト機能用タスク一覧取得
- シフト機能用タスク追加
- シフト機能用タスク更新
- シフト機能用タスク削除

## 課題と今後実装したい機能

- コードの品質保証ができていない。
- [ ] バリデーション処理や複雑な処理をしている部分等を優先にテストを充実させていく。


- kubernetesにて複数環境に対応したmanifest管理ができていない。
- [x] kustomizeを導入し、環境毎に設定を定義。


- terraformにて複数環境に対応したリソース管理ができていない。
- [ ] 共通リソースをmodule化し、環境毎パラメータ設定ファイルを追加。

## Documentation for API Endpoints

### user-rest-service

| HTTP request | HTTP response<br>status code | Description |
| :--- | :--- | :--- |
| **POST**<br>&emsp;/signup | <pre>201 Created<br>400 Bad Request<br>409 Conflict<br>500 Internal Server Error</pre> | ユーザー新規登録 |
| **POST**<br>&emsp;/login | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | ユーザーログイン |
| **DELETE**<br>&emsp;/logout | <pre>200 OK<br>400 Bad Request<br>500 Internal Server Error</pre> | ユーザーログアウト |
| **GET**<br>&emsp;/user | <pre>200 OK<br>401 Unauthorized<br>404 Not Found<br>500 Internal Server Error</pre> | ユーザー情報取得 |
| **GET**<br>&emsp;/groups | <pre>200 OK<br>401 Unauthorized<br>500 Internal Server Error</pre> | 承認グループ,<br>未承認グループ一覧取得 |
| **POST**<br>&emsp;/groups | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ作成 |
| **PUT**<br>&emsp;/groups/{group_id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ名更新 |
| **POST**<br>&emsp;/groups/{group_id}/users | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error</pre> | グループ招待 |
| **DELETE**<br>&emsp;/groups/{group_id}/users | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ退会 |
| **POST**<br>&emsp;/groups/{group_id}/users/approved | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ招待承認 |
| **DELETE**<br>&emsp;/groups/{group_id}/users/unapproved | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ招待拒否 |
| **GET**<br>&emsp;/readyz | <pre>200 OK<br>503 Service Unavailable</pre> | Kubernetesヘルスチェック<br>※ Backend専用 |
| **GET**<br>&emsp;/groups/{group_id}/users | <pre>200 OK<br>400 Bad Request<br>500 Internal Server Error</pre> | グループユーザーIDリスト取得<br>※ Backend専用 |
| **GET**<br>&emsp;/groups/{group_id}/users/{user_id}/verify | <pre>200 OK<br>400 Bad Request<br>500 Internal Server Error</pre> | グループ所属確認<br>※ Backend専用 |
| **GET**<br>&emsp;/groups/{group_id}/users/verify | <pre>200 OK<br>400 Bad Request<br>500 Internal Server Error</pre> | ユーザーリストのグループ所属確認<br>※ Backend専用 |

### account-rest-service

| HTTP request | HTTP response<br>status code | Description |
| :--- | :--- | :--- |
| **GET**<br>&emsp;/categories | <pre>200 OK<br>401 Unauthorized<br>500 Internal Server Error</pre> | カテゴリー一覧取得 |
| **POST**<br>&emsp;/categories/custom-categories | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error</pre> | カスタムカテゴリー追加 |
| **PUT**<br>&emsp;/categories/custom-categories/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error</pre> | カスタムカテゴリー更新 |
| **DELETE**<br>&emsp;/categories/custom-categories/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>404 Not Found<br>500 Internal Server Error</pre> | カスタムカテゴリー削除 |
| **GET**<br>&emsp;/transactions/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別家計簿トランザクション一覧取得 |
| **GET**<br>&emsp;/transactions/latest | <pre>200 OK<br>401 Unauthorized<br>500 Internal Server Error</pre> | 家計簿トランザクション最終更新10件取得 |
| **POST**<br>&emsp;/transactions | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 家計簿トランザクション追加 |
| **PUT**<br>&emsp;/transactions/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 家計簿トランザクション更新 |
| **DELETE**<br>&emsp;/transactions/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 家計簿トランザクション削除 |
| **GET**<br>&emsp;/transactions/search | <pre>200 OK<br>401 Unauthorized<br>500 Internal Server Error</pre> | 家計簿トランザクション検索 |
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
| **DELETE**<br>&emsp;/groups/{group_id}/categories/custom-categories/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>404 Not Found<br>500 Internal Server Error</pre> | グループカスタムカテゴリー削除 |
| **GET**<br>&emsp;/groups/{group_id}/transactions/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 月別グループ家計簿トランザクション一覧取得 |
| **GET**<br>&emsp;/groups/{group_id}/transactions/latest | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿トランザクション最終更新10件取得 |
| **POST**<br>&emsp;/groups/{group_id}/transactions | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿トランザクション追加 |
| **PUT**<br>&emsp;/groups/{group_id}/transactions/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿トランザクション更新 |
| **DELETE**<br>&emsp;/groups/{group_id}/transactions/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿トランザクション削除 |
| **GET**<br>&emsp;/groups/{group_id}/transactions/search | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿トランザクション検索 |
| **GET**<br>&emsp;/groups/{group_id}/transactions/{year}/account | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 年別グループ家計簿トランザクション会計状況一覧取得 |
| **GET**<br>&emsp;/groups/{group_id}/transactions/{year_month}/account | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>404 Not Found<br>500 Internal Server Error</pre> | 月別グループ家計簿トランザクション会計データ取得 |
| **POST**<br>&emsp;/groups/{group_id}/transactions/{year_month}/account | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>404 Not Found<br>500 Internal Server Error</pre> | 月別グループ家計簿トランザクション会計データ追加 |
| **PUT**<br>&emsp;/groups/{group_id}/transactions/{year_month}/account/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>404 Not Found<br>500 Internal Server Error</pre> | 月別グループ家計簿トランザクション会計データ更新 |
| **DELETE**<br>&emsp;/groups/{group_id}/transactions/{year_month}/account | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>404 Not Found<br>500 Internal Server Error</pre> | 月別グループ家計簿トランザクション会計データ削除 |
| **GET**<br>&emsp;/groups/{group_id}/standard-budgets | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿標準予算取得 |
| **PUT**<br>&emsp;/groups/{group_id}/standard-budgets | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ家計簿標準予算更新 |
| **GET**<br>&emsp;/groups/{group_id}/custom-budgets/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ月別家計簿カスタム予算取得 |
| **POST**<br>&emsp;/groups/{group_id}/custom-budgets/{year_month} | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ月別家計簿カスタム予算追加 |
| **PUT**<br>&emsp;/groups/{group_id}/custom-budgets/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ月別家計簿カスタム予算更新 |
| **DELETE**<br>&emsp;/groups/{group_id}/custom-budgets/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ月別家計簿カスタム予算削除 |
| **GET**<br>&emsp;/groups/{group_id}/budgets/{year} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ年別家計簿予算一覧取得 |
| **GET**<br>&emsp;/readyz | <pre>200 OK<br>503 Service Unavailable</pre> | Kubernetesヘルスチェック<br>※ Backend専用 |
| **POST**<br>&emsp;/standard-budgets | <pre>201 Created<br>500 Internal Server Error</pre> | 家計簿標準予算初期値追加<br>※ Backend専用 |
| **GET**<br>&emsp;/categories/name | <pre>200 OK<br>500 Internal Server Error</pre> | カテゴリー名取得<br>※ Backend専用 |
| **GET**<br>&emsp;/categories/names | <pre>200 OK<br>500 Internal Server Error</pre> | カテゴリー名リスト取得<br>※ Backend専用 |
| **GET**<br>&emsp;/transactions/related-shopping-list | <pre>200 OK<br>500 Internal Server Error</pre> | 買い物リストに関連する家計簿トランザクション取得<br>※ Backend専用 |
| **POST**<br>&emsp;/groups/{group_id}/standard-budgets | <pre>201 Created<br>400 Bad Request<br>500 Internal Server Error</pre> | グループ家計簿標準予算初期値追加<br>※ Backend専用 |
| **GET**<br>&emsp;/groups/{group_id}/categories/name | <pre>200 OK<br>500 Internal Server Error</pre> | グループカテゴリー名取得<br>※ Backend専用 |
| **GET**<br>&emsp;/groups/{group_id}/categories/names | <pre>200 OK<br>500 Internal Server Error</pre> | グループカテゴリー名リスト取得<br>※ Backend専用 |
| **GET**<br>&emsp;/groups/{group_id}/transactions/related-shopping-list | <pre>200 OK<br>500 Internal Server Error</pre> | グループ買い物リストに関連する家計簿トランザクション取得<br>※ Backend専用 |

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
| **GET**<br>&emsp;/shopping-list/{date}/daily | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 定期買い物リスト,<br>日別日順買い物リスト取得 |
| **GET**<br>&emsp;/shopping-list/{date}/categories | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 定期買い物リスト,<br>日別カテゴリー順買い物リスト取得 |
| **GET**<br>&emsp;/shopping-list/{year_month}/daily | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 定期買い物リスト,<br>月別日順買い物リスト取得 |
| **GET**<br>&emsp;/shopping-list/{year_month}/categories | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 定期買い物リスト,<br>月別カテゴリー順買い物リスト取得 |
| **GET**<br>&emsp;/shopping-list/expired | <pre>200 OK<br>401 Unauthorized<br>500 Internal Server Error</pre> | 期限切れ買い物リスト取得 |
| **POST**<br>&emsp;/shopping-list | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 買い物リスト追加 |
| **PUT**<br>&emsp;/shopping-list/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 買い物リスト更新,<br>家計簿トランザクション自動追加/自動削除 |
| **DELETE**<br>&emsp;/shopping-list/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 買い物リスト削除 |
| **POST**<br>&emsp;/shopping-list/regular | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 定期買い物リスト追加 |
| **PUT**<br>&emsp;/shopping-list/regular/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 定期買い物リスト更新 |
| **DELETE**<br>&emsp;/shopping-list/regular/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | 定期買い物リスト削除 |
| **GET**<br>&emsp;/groups/{group_id}/todo-list/{date} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ日別実施予定todo,<br>グループ締切予定todo一覧取得 |
| **GET**<br>&emsp;/groups/{group_id}/todo-list/{year_month} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ月別実施予定todo,<br>グループ締切予定todo一覧取得 |
| **GET**<br>&emsp;/groups/{group_id}/todo-list/expired | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ期限切れtodo一覧取得 |
| **POST**<br>&emsp;/groups/{group_id}/todo-list | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループtodo追加 |
| **PUT**<br>&emsp;/groups/{group_id}/todo-list/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループtodo更新 |
| **DELETE**<br>&emsp;/groups/{group_id}/todo-list/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループtodo削除 |
| **GET**<br>&emsp;/groups/{group_id}/todo-list/search | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループtodo検索 |
| **GET**<br>&emsp;/groups/{group_id}/shopping-list/{date}/daily | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ定期買い物リスト,<br>日別日順買い物リスト取得 |
| **GET**<br>&emsp;/groups/{group_id}/shopping-list/{date}/categories | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ定期買い物リスト,<br>日別カテゴリー順買い物リスト取得 |
| **GET**<br>&emsp;/groups/{group_id}/shopping-list/{year_month}/daily | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ定期買い物リスト,<br>月別日順買い物リスト取得 |
| **GET**<br>&emsp;/groups/{group_id}/shopping-list/{year_month}/categories | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ定期買い物リスト,<br>月別カテゴリー順買い物リスト取得 |
| **GET**<br>&emsp;/groups/{group_id}/shopping-list/expired | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ期限切れ買い物リスト取得 |
| **POST**<br>&emsp;/groups/{group_id}/shopping-list | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ買い物リスト追加 |
| **PUT**<br>&emsp;/groups/{group_id}/shopping-list/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ買い物リスト更新,<br>家計簿トランザクション自動追加/自動削除 |
| **DELETE**<br>&emsp;/groups/{group_id}/shopping-list/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ買い物リスト削除 |
| **POST**<br>&emsp;/groups/{group_id}/shopping-list/regular | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ定期買い物リスト追加 |
| **PUT**<br>&emsp;/groups/{group_id}/shopping-list/regular/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ定期買い物リスト更新 |
| **DELETE**<br>&emsp;/groups/{group_id}/shopping-list/regular/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループ定期買い物リスト削除 |
| **GET**<br>&emsp;/groups/{group_id}/tasks/users | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | ユーザー別グループタスク一覧取得 |
| **POST**<br>&emsp;/groups/{group_id}/tasks/users | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>409 Conflict<br>500 Internal Server Error</pre> | グループタスクユーザー追加 |
| **DELETE**<br>&emsp;/groups/{group_id}/tasks/users | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループタスクユーザー削除 |
| **GET**<br>&emsp;/groups/{group_id}/tasks | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループタスク一覧取得 |
| **POST**<br>&emsp;/groups/{group_id}/tasks | <pre>201 Created<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループタスク追加 |
| **PUT**<br>&emsp;/groups/{group_id}/tasks/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループタスク更新 |
| **DELETE**<br>&emsp;/groups/{group_id}/tasks/{id} | <pre>200 OK<br>400 Bad Request<br>401 Unauthorized<br>500 Internal Server Error</pre> | グループタスク削除 |
| **GET**<br>&emsp;/readyz | <pre>200 OK<br>503 Service Unavailable</pre> | Kubernetesヘルスチェック<br>※ Backend専用 |
