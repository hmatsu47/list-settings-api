package api

import (
	"github.com/gin-gonic/gin"
)

type ListSettings struct {
	TagRepositoryUri	string
	SelectTags			[]string
}

func NewListSettings(tagRepositoryUri string, selectTags []string) *ListSettings {
	return &ListSettings{
		TagRepositoryUri: tagRepositoryUri,
		SelectTags:       selectTags,
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
func (s *ListSettings) GetUriSettings(c *gin.Context) {

	// c.JSON(http.StatusOK, result)
}


// リリース設定一覧の取得（タグ指定分）
func (s *ListSettings) GetTagSettings(c *gin.Context) {

	// c.JSON(http.StatusOK, result)
}
