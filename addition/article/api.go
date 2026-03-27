package article

import (
	"chat/auth"
	"chat/globals"
	"chat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

type WebsocketArticleForm struct {
	Token  string `json:"token" binding:"required"`
	Model  string `json:"model" binding:"required"`
	Prompt string `json:"prompt" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Web    bool   `json:"web"`
}

type WebsocketArticleResponse struct {
	Hash string                 `json:"hash"`
	Data StreamProgressResponse `json:"data"`
}

func ProjectTarDownloadAPI(c *gin.Context) {
	hash := strings.TrimSpace(c.Query("hash"))
	// Prevent path traversal by validating hash contains only safe characters
	if hash == "" || strings.ContainsAny(hash, "/\..") {
		c.JSON(400, gin.H{"error": "invalid hash"})
		return
	}
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=article.tar.gz")
	c.File(filepath.Join("storage", "article", hash+".tar.gz"))
}

func ProjectZipDownloadAPI(c *gin.Context) {
	hash := strings.TrimSpace(c.Query("hash"))
	// Prevent path traversal by validating hash contains only safe characters
	if hash == "" || strings.ContainsAny(hash, "/\..") {
		c.JSON(400, gin.H{"error": "invalid hash"})
		return
	}
	c.Writer.Header().Add("Content-Disposition", "attachment; filename=article.zip")
	c.File(filepath.Join("storage", "article", hash+".zip"))
}

func GenerateAPI(c *gin.Context) {
	var conn *utils.WebSocket
	if conn = utils.NewWebsocket(c, false); conn == nil {
		return
	}
	defer conn.DeferClose()

	form, err := utils.ReadForm[WebsocketArticleForm](conn)
	if err != nil {
		return
	}

	user := auth.ParseToken(c, form.Token)
	db := utils.GetDBFromContext(c)

	if !auth.HitGroups(db, user, globals.ArticlePermissionGroup) {
		return
	}

	if len(form.Title) == 0 {
		return
	}

	hash := CreateWorker(c, user, form.Model, form.Prompt, form.Title, form.Web, func(resp StreamProgressResponse) {
		conn.Send(WebsocketArticleResponse{
			Hash: "",
			Data: resp,
		})
	})
	conn.Send(WebsocketArticleResponse{
		Hash: hash,
	})
}
