package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/hmatsu47/list-settings-api/api"
)

func NewGinListSettingsServer(listSettings *api.ListSettings, port int) *http.Server {
	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Swagger specの読み取りに失敗しました\n: %s", err)
		os.Exit(1)
	}

	// Swagger Document 非公開
	swagger.Servers = nil

	// Gin Router 設定
	r := gin.Default()

	// HTTP Request の Validation 設定
	r.Use(middleware.OapiRequestValidator(swagger))

	// Handler 実装
	r = api.RegisterHandlers(r, listSettings)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
	}
	return s
}

func main() {
	port := flag.Int("port", 28080, "Port for API server")
	configPathPrefix := ""
	tagRepositoryUri := ""
	flag.Parse()
	if flag.NArg() == 0 {
		panic("タグの情報が指定されていません")
	}
	// タグと対応する環境名はコマンドラインパラメータで取得
	var tagKeys []api.TagKey
	for i, v := range flag.Args() {
		if i == 0 {
			configPathPrefix = v
		} else if i == 1 {
			tagRepositoryUri = v
		} else {
			tag := api.TagKey{
				TagName:         strings.Split(v, ":")[0],
				EnvironmentName: strings.Split(v, ":")[1],
			}
			tagKeys = append(tagKeys, tag)
		}
	}
	if tagRepositoryUri == "" {
		panic("タグ形式一覧用のECRリポジトリURIが指定されていません")
	}
	// Server Instance 生成
	listSettings := api.NewListSettings(configPathPrefix, tagRepositoryUri, &tagKeys)
	s := NewGinListSettingsServer(listSettings, *port)
	// 停止まで HTTP Request を処理
	log.Fatal(s.ListenAndServe())
}
