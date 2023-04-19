package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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
	flag.Parse()
	// リポジトリ URI・付与するタグはコマンドラインパラメータで取得
	tagRepositoryUri := flag.Arg(0)
	if tagRepositoryUri == "" {
		panic("リポジトリの指定がありません")
	}
	var selectTags []string
	// Server Instance 生成
	listSettings := api.NewListSettings(tagRepositoryUri, selectTags)
	s := NewGinListSettingsServer(listSettings, *port)
	// 停止まで HTTP Request を処理
	log.Fatal(s.ListenAndServe())
}
