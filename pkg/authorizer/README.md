# Authorizer
authenticator で提供される JWT の検証を行います。

## 使用方法
`pkg/authorizer` をインポートしてください。

### 事前準備
JWT のデジタル署名の検証に authenticator の公開鍵を使用します。
生成された公開鍵をコピーして、authorizer を使用する pod に配置してください。
また、使用する pod の環境変数 `CREDENTIAL_FILE_PATH` に配置した公開鍵の file_path 設定する必要があります。 

### JWT の検証
`pkg/authorizer/authorizer.go` で提供される関数 `VerifyJWTToken()` を使用します。
JWT の改竄の検証、および有効期限を確認が可能です。

下記のサンプルコードは JWT を検証し、有効の場合に クレームの内容を出力します。
```example
func printJwtClaims() {
	const Jwt = "xxxxx.xxxxx.xxxxx"
	jwtToken, err := VerifyJWTToken(Jwt)
	if err != nil {
		log.Fatalf("Invalid JWT: %v", err)
	}
	log.Printf("user_id = %v", jwtToken.Claims)
}
```
