# Authenticator

## Description
認証サービス  
- 認証

### セットアップ
```shell
# statikのインストール
$go get github.com/rakyll/statik

# configs/configs.yamlに関連情報記述

# statikによるシングルバイナリためのファイル生成
$ make statik

# image build 
$ make docker-build

# デプロイ
$ kubectl apply -f deployments/deployments.yaml
```


### エンドポイント
```markdown
- post /user
- post /login
```