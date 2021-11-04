# authenticator

## 概要
authenticator は、主にエッジコンピューティング環境において、任意のアプリケーション・システム・マイクロサービス等に対して、認証機能を提供するマイクロサービスです。  
authenticator は、次のようなシステム環境において、特に重要な役割を果たします。  

* 多くのモバイル端末などがエッジNWクラスターに接続する、クラウドから遮断されたエッジコンピューティングスタック環境  

* 複数のモバイル端末などがエッジ端末に接続する、エッジアプリケーションシステム環境  

* あるエッジ端末がデータ連携先の別のエッジ端末に接続する、エッジコンピューティング環境    


## システム構成図  
![img](docs/authenticator.png)

## authenticator が含まれる アーキテクチャ（OMOTE-Bakoアプリケーションの例）  
authenticatorは、下記の黄色い枠の部分のリソースです。
![img](docs/omotebako_architecture_20211104_authenticator.png)


## 事前準備
本マイクロサービスは、主にエッジコンピューティング環境において認証情報を維持管理するためのDBとして、MySQLを利用します。
デフォルトでは Users テーブルを参照しますが、configs/configs.yamlに対象テーブルを指定していただければ、指定されたテーブル内容を参照します。  
この場合 id(int),login_id(string),password(string) のカラムが必ず存在することが必要です。 

kubectlの事前インストールが必要です
```shell
$ brew install kubectl
```

## セットアップ
1. リポジトリをクローンする。
```shell
$ git clone https://github.com/latonaio/authenticator.git
$ cd authenticator
```

2. `sh scripts/setup-configs.sh` を実行する。   
    ※ shellscript内での具体的な処理内容は以下の通り。
```
# 秘密鍵の生成 (JWT:JSON Web Token の署名に必要)
$ make generate-key-pair

# configs/configs.yamlの設定を編集
# HOST_NAME, PORT, DB_USER_NAME, DB_USER_PASSWORDを変更する
database:
  host_name: HOST_NAME
  port: PORT
  user_name: DB_USER_NAME
  user_password: DB_USER_PASSWORD
  name: Authenticator # database name
  table_name: Users # table name

# configs/configs.yamlの設定を編集
# private_key に生成した秘密鍵をコピペする
private_key: |-
  -----BEGIN RSA PRIVATE KEY-----
  XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/XXXXXXXX
  (中略)
  XXX/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/XXXXXXXXX=
  -----END RSA PRIVATE KEY-----
```
3. Docker Imageの生成
```
$ make docker-build
```

## 起動方法
authenticatorはaion-core上で稼働します。したがって、[aion-service-definitions](https://bitbucket.org/latonaio/aion-service-definitions/src/master/)にある、aion-coreのマニフェストファイルであるservices.ymlにauthenticatorが起動するように設定する必要があります。記入例は以下の通りです。

```
microservices:
  authenticator:         # サービス名
    startup: yes         # AION起動時に起動
    always: yes          # 常時起動
    scale: 1             # 複数同時に起動させない
    network: NodePort
    ports:
    - name: authenticator
      protocol: TCP
      port: 50500         # kube内のPod間で通信する際に使うポート番号
      nodePort: 30500     # kube外からの通信に必要なポート番号
```

## データベース管理
データベースの migration には goose を使用します。
下記コマンドで goose をインストールしてください
```
$ go get bitbucket.org/liamstask/goose/cmd/goose
```

goose のデータベース接続設定は `db/dbconf.yml` に記載されています。必要に応じて書き換えてください。
また、デフォルトで使用するデータベースのテーブルは `db/migrations/xxxxx_createUsersTable.sql (最も古いもの)` に定義されています。

migration に関するコマンドは以下です。

```
# 疎通確認
$ goose status

# マイグレーション
$ goose up

# ロールバック
$ goose down

# マイグレーションファイルの作成
$ goose create XXXXXX sql
```

データベースのスキーマを変更する場合は、マイグレーションファイルを作成し SQL を記述後マイグレーションを実行してください。

## モデル定義
### user
user は下記の項目を保持しています。より厳密な定義を知りたい場合は DB 定義を確認してください。

| name | description |
| --- | --- |
| login_id | ログイン ID |
| password | パスワード |
| qos | quality of service|

`login_id`, `password` には入力規則があります。詳しくは入力規則の項目を参照してください。

`qos` は現状 `default`、`raw` の二つの値のどちらかを登録できます。
入力しない場合や、これら二つ以外の値を登録しようとすると `default` が設定されます。
`qos` の値を `raw` にする事で、`authenticator` の入力規則を無視した `login_id`, `password` を登録することができます。

## 利用方法
Authenticatorでは以下のAPIが利用できます。

### POST /users
ユーザー登録を行います。

#### リクエスト
ユーザー登録には下記のパラメータを指定して、 POST リクエストを送信してください。

| name | description |
| --- | --- |
| login_id | ユーザを識別する id (必須) |
| password | ユーザ認証を行う password (必須) |
| qos | qos の値を指定します。デフォルトでは "default" が設定されます。|

login_id と password の入力規則に関しては 入力規則 のセクションを参照してください

```
POST /users
Origin: http://{{host_name}}
Content-Type: application/x-www-form-urlencoded

login_id=Sample_user&password=OK_password&qos=default
```

```
curl -X POST http://{{host_name}}/users -d login_id=Sample_user -d password=OK_password -d qos=default
```

```
curl -X POST http://{{host_name}}/users -d login_id=sampleuser -d password=okpassword -d qos=raw
```

#### レスポンス
authenticator はリクエストに対し下記のいずれかの応答をします。

| status code | description |
| --- | --- |
| 200 | ユーザーの登録に成功 |
| 400 | リクエストパラメータが不正 (入力規則を満たしているか確認してください)|
| 409 | login_id が既に登録済み or 削除されたユーザ |
| 500 | サーバー内エラー |

### GET /users/login_id/{{login_id}}
ユーザー情報を取得します。

#### リクエスト
ユーザー情報の取得に GET リクエストを送信してください。

```http_request
GET /users/login_id/{{logain_id}}
Origin: http://{{host_name}}
```

```
curl -X GET http://{{host_name}}/users/login_id/{{logain_id}}
```

#### レスポンス
authenticator はリクエストに対し下記のいずれかの応答をします。

| status code | description |
| --- | --- |
| 200 | user 情報の返却 |
| 404 | ユーザが未登録 |
| 409 | ユーザが削除済み |
| 500 | サーバー内エラー |

ユーザーが登録されている場合、ユーザー情報を返します。

```response-example
{"login_id":"xxxxxx"}
```

### PUT /users/login_id/{{login_id}}
ユーザー更新を行います。

#### リクエスト
ユーザー情報の更新には下記のパラメータを指定して、 PUT リクエストを送信してください。

| name | description |
| --- | --- |
| old_password | 変更前の password (必須) |
| login_id | ユーザを識別する id ( 指定しない場合は更新されません ) |
| password | ユーザ認証を行う password ( 指定しない場合は更新されません ) |
| qos | qos の値 ( 指定しない場合は更新されません ) |

login_id と password の入力規則に関しては 入力規則 のセクションを参照してください。
また、更新の際には old_password に更新前のパスワードを指定し、認可する必要があります。

```http_request
PUT /users/login_id/{{logain_id}}
Origin: http://{{host_name}}
Content-Type: application/x-www-form-urlencoded

old_password=OK_password&login_id=sampleuser&password=okpassword&qos=raw
```

```
curl -X PUT http://{{host_name}}/users/login_id/{{logain_id}} -d old_password=OK_password -d login_id=sampleuser -d password=okpassword -d qos=raw
```

#### レスポンス

authenticator はリクエストに対し下記のいずれかの応答をします。

| status code | description |
| --- | --- |
| 200 | user の更新に成功 |
| 400 | リクエストパラメータが不正 |
| 401 | ユーザーの認証に失敗 |
| 404 | ユーザが未登録 |
| 409 | 既に削除されたユーザー |
| 500 | サーバー内エラー |

### POST /login
ユーザー認証を行います。

#### リクエスト
ユーザー認証には下記のパラメータを指定して、 POST リクエストを送信してください。
認証に成功した場合 JWT が返却されます。

| name | description |
| --- | --- |
| login_id | ユーザを識別する id (必須) |
| password | ユーザ認証を行う password (必須) |

```http_request
POST /login
Origin: http://{{host_name}}
Content-Type: application/x-www-form-urlencoded

login_id=Sample_user&password=OK_password
```

```request_example
curl -X POST http://{{host_name}}/login -d login_id=Sample_user -d password=OK_password
```

#### レスポンス
authenticator はリクエストに対し下記のいずれかの応答をします。

| status code | description |
| --- | --- |
| 200 | ユーザーの認証に成功 (JWT を返却)|
| 400 | リクエストパラメータが不正 |
| 401 | ユーザーの認証に失敗 |
| 404 | ユーザが未登録 |
| 409 | 既に削除済みのユーザ |
| 500 | サーバー内エラー |

認証に成功すると JWT を返却します。

```response-example
{"jwt":"xxxxx.xxxxx.xxxxx"}
```

### POST /users/login_id/{{login_id}}
ユーザーの削除を行います。内部的にはデータの削除は行われずに、users table の deleted_at の項目に現在時刻が登録されます。

#### リクエスト
ユーザー情報の更新には下記のパラメータを指定して、 DELETE リクエストを送信してください。

| name | description |
| --- | --- |
| password | 削除対象のユーザの password |

password での認証に失敗した場合は削除されません。

```http_request
DELETE /users/login_id/{{logain_id}}
Origin: http://{{host_name}}
Content-Type: application/x-www-form-urlencoded

password=okpassword
```

```
curl -X DELETE http://{{host_name}}/users/login_id/{{logain_id}} -d password=password
```

#### レスポンス
authenticator はリクエストに対し下記のいずれかの応答をします。

| status code | description |
| --- | --- |
| 200 | ユーザーの削除に成功|
| 400 | リクエストパラメータが不正 |
| 401 | ユーザーの認証に失敗 |
| 404 | ユーザが未登録 |
| 409 | 既に削除済みのユーザ |
| 500 | サーバー内エラー |

### 認可
authenticator サーバーから取得した jwt を認可するには、下記のコマンドで出力される公開鍵を認可サーバーに配置する必要があります。
```
$ make generate-key-pair
```

認可サーバーでは authenticator サーバーの公開鍵を使用し jwt が改竄されていない事、有効期限切れでない事を確認する必要があります。 
`pkg/authrizer` でこれらを確認するライブラリを提供しています。

### JWT
JET のクレームには下記の項目が含まれています。

| name | description|
| --- | --- |
| user_id | ユーザーID |
| exp | 有効期限 |

## テスト

### Unit テスト
authorizer のテストの際に秘密鍵、公開鍵のキーペアを使用します。
下記のコマンドでキーペアを生成し、それぞれを unit テストに直接入力してください。

`$ openssl genrsa 4096 > private.key`
`$ openssl rsa -pubout < private.key > public.key`

```
# pkg/authorizer/authorizer_test.go

func TestValidateJWTToken(t *testing.T) {
	const privateKey = "" <- ここに private.key の内容を入力
	const publicKey = "" <- ここに public.key の内容を入力
```


#### 秘密鍵の配置
authenticator は JWT の生成に秘密鍵を使用します。
下記コマンドで生成される秘密鍵を `configs/configs.yaml` の `private_key` 項目にセットしてください

```
$ make generate-key-pair
```

example: configs/configs.yaml

```
server:
  port: 1323
  shutdown_wait_time: 1
database:
  host_name: mysql
  port: 3306
  user_name: xxxxx
  user_password: xxxxx
  name: Authenticator # database name
  table_name: Users # table name 強制指定のため
  max_open_connection: 10
  max_idle_connection: 5
  max_life_time: 24
jwt:
  exp: 24 # hour
private_key: |-
  -----BEGIN RSA PRIVATE KEY-----
  XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
  XXXX/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/XXXXXXXXXXXXXXXXXXXX
  (中略)
  XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX
  XXX/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/XXXXXXXXXX/XXXXXXXXXXXXXX
  -----END RSA PRIVATE KEY-----
```

#### ローカル環境で実行 (おすすめ)
mysql のみ docker-compose で起動し、authenticator をローカル環境で起動することも可能です。
この場合下記コマンドで authenticator を立ち上げることで、秘密鍵の生成と環境変数へのセットを自動で行います。

`$ make local-run` 

`configs/configs.yaml` は下記のように設定してください
```
database:
  host_name: localhost
  port: 3306
  # 他はそのまま
```

#### 登録
```
curl -X POST http://localhost:1323/users -d login_id=Sample_user -d password=OK_password
```

#### 取得
```
curl -X GET http://localhost:1323/users/login_id/Sample_user
```

#### ログイン
```
curl -X POST http://localhost:1323/login -d login_id=Sample_user -d password=OK_password
```

#### 更新
```
curl -X PUT http://localhost:1323/users/login_id/Sample_user -d old_password=OK_password -d login_id=sampleuser -d password=okpassword -d qos=raw
```

#### 削除
```
curl -X POST http://localhost:1323/users/login_id/Sample_user -d password=OK_password
```


## 入力規則
user を新規登録する際は下記の入力規則にしたがって登録してください。

### login_id

#### 使用可能文字
アルファベット（a～z, A〜Z 大文字小文字を区別する）、数字（0～9）、記号(ダッシュ（-）、アンダースコア（_）、アポストロフィ（'）、ピリオド（.）)

#### 文字数制限
6〜30 文字

### password

#### 使用可能文字
アルファベット（a～z, A〜Z 大文字小文字を区別する）、数字（0～9）、記号(ダッシュ（-）、アンダースコア（_）、アポストロフィ（'）、ピリオド（.）)

#### パスワード長制限
8〜30 文字

#### その他の条件
- ユーザ名の文字列がそのままパスワードの文字列に含まれていないこと
- アルファベットの大文字、小文字がそれぞれ 1 文字以上含まれていること

### 入力規則を無視する
qos (quality of service) を "raw" とすることで入力規則を無視してユーザーの登録、更新が可能です。
