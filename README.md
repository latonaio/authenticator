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

