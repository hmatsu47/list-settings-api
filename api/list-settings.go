package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TagKey struct {
	TagName				string
	EnvironmentName		string
}
type ListSettings struct {
	ConfigPathPrefix	string
	TagRepositoryUri	string
	TagKeys				*[]TagKey
}

func NewListSettings(configPathPrefix string, tagRepositoryUri string, tagKeys *[]TagKey) *ListSettings {
	return &ListSettings{
		ConfigPathPrefix: configPathPrefix,
		TagRepositoryUri: tagRepositoryUri,
		TagKeys:          tagKeys,
	}
}

// エラーメッセージ返却用
func sendError(c *gin.Context, code int, message string) {
	selectErr := Error{
		Message: message,
	}
	c.JSON(code, selectErr)
}

// リリース設定一覧の取得（URI指定分）
func (l *ListSettings) GetUriSettings(c *gin.Context) {
	result, err := ReadSettings(l.ConfigPathPrefix)
	if err != nil {
		sendError(c, http.StatusInternalServerError, fmt.Sprintf("リリース設定の取得でエラーが発生しました : %s", err))
	}
	c.JSON(http.StatusOK, result)
}


// リリース設定一覧の取得（タグ指定分）
func (l *ListSettings) GetTagSettings(c *gin.Context) {
	var result []TagSetting
	region := strings.Split(l.TagRepositoryUri, ".")[3]
	ecrClient, err := EcrClient(region)
	if err != nil {
		sendError(c, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	result, err = TagSettingList(context.TODO(), ecrClient, l.TagRepositoryUri, *l.TagKeys)
	if err != nil {
		sendError(c, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	c.JSON(http.StatusOK, result)
}
