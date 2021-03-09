# Authenticator

## Description
authenticatorはAIONのプラットフォーム上で動作するシステムに対して認証機能を提供するマイクロサービスです。  


## 事前準備
本マイクロサービスはDBにMySQLを利用します。  
また、misc/dump.sqlをMySQLにインポートする必要があります。  

デフォルトで使用するデータベース名とユーザーテーブルは misc/dump.sqlに用意されています,  
またconfigs/configs.yamlに対象テーブルを指定していただければ、指定されたユーザーテーブル内容を参照します。  
この場合 id(int),login_id(string),password(string) のカラムが必ず存在することが必要です。  

kubectlを事前にインストール必要
```shell
$ brew install kubectl
```

## セットアップ
```shell
$ git clone https://github.com/latonaio/authenticator.git
$ cd authenticator

# configs/configs.yamlの設定を編集
# HOST_NAME, PORT, DB_USER_NAME, DB_USER_PASSWORDを変更する
database:
  host_name: HOST_NAME
  port: PORT
  user_name: DB_USER_NAME
  user_password: DB_USER_PASSWORD
  name: Authenticator # database name
  table_name: Users # table name

# Docker Imageの生成
$ make docker-build
```

## 起動方法
```shell
#kubectlをinstall必要
## 起動方法　
$ kubectl apply -f deployments/deployments.yaml
```


## 利用方法
Authenticatorでは以下のAPIが利用できます。
```text
ユーザー認証
POST /login
```

```text
ユーザー作成
POST /user
```

## システム構成図
![img](docs/authenticator.png)


### テスト

### Unit テスト
authorizer のテストの際に秘密鍵、公開鍵のキーペアを使用します。
下記のコマンドでキーペアを生成し、それぞれを unit テストに直接入力してください。

`$ openssl rsa -pubout < private.key > public.key`


```
# pkg/authorizer/authorizer_test.go

func TestValidateJWTToken(t *testing.T) {
	const privateKey = "" <- ここに private.key の内容を入力
	const publicKey = "" <- ここに public.key の内容を入力
```

### 動作確認
`docker-compose` を使用して mysql と authenticator コンテナを立ててテストを行います。
ホスト側のポート 1323 を使用して authenticator と通信します。

#### セットアップ
`configs/configs.yaml` の database セクションの 設定を変更します。
```
database:
  host_name: mysql
  port: 3306
  # 他はそのまま
```

#### コンテナの起動
authenticator コンテナは起動後すぐに mysql とのコネクションを確立しようとします。
そのため、 mysql コンテナの起動が完了してから authenticator コンテナを起動してください。

```shell
# build
$ docker-compose build

# mysql コンテナを起動
$ docker-compose up mysql

# authenticator コンテナを起動
$ docker-compose up authenticator

```

#### user の登録・ログイン

user を登録します。
```
curl -X POST http://localhost:1323/users -d user_id=sample_user -d password=sample_user
```

ログインします。
```
curl -X POST http://localhost:1323/login -d user_name=sample_user -d user_password=sample_user
```