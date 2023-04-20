# list-settings-api

Go で内部管理用 API を作るテスト (3)

## `.yaml`ファイルから API コードの枠組みを生成

- こちらを踏襲
  - https://github.com/hmatsu47/select-repository-api

```sh:install
go mod init github.com/hmatsu47/select-repository-api
mkdir internal
cd internal
（作成した`.yaml`ファイルを`internal`内にコピー）
cd ..
mkdir api
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4
oapi-codegen -output-config -old-config-style -package=api -generate=types -alias-types internal/list-settings-api.yaml > api/config-types.yaml
oapi-codegen -output-config -old-config-style -package=api -generate=gin,spec -alias-types internal/list-settings-api.yaml > api/config-server.yaml
oapi-codegen -config api/config-types.yaml internal/list-settings-api.yaml > api/types.gen.go
oapi-codegen -config api/config-server.yaml internal/list-settings-api.yaml > api/server.gen.go
go mod tidy
```

## 起動方法

`go run main.go [-port=待機ポート番号（TCP）] CORS許可URL（カンマ区切り複数指定可） URI形式一覧の設定ファイル保存先パスプレフィックス タグ形式一覧の対象ECRリポジトリURI タグ形式一覧の付与タグ1:タグ形式一覧の環境名1 [タグ形式一覧の付与タグ2:タグ形式一覧の環境名2 [タグ形式一覧の付与タグ3:タグ形式一覧の環境名3 ...]]`
