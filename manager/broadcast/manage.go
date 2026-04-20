package broadcast

import (
	"chat/auth"
	"chat/globals"
	"chat/utils"
	"context"
	"github.com/gin-gonic/gin"
)

func createBroadcast(c *gin.Context, user *auth.User, content string) error {
	db := utils.GetDBFromContext(c)
	cache := utils.GetCacheFromContext(c)

	if _, err := globals.ExecDb(db, `INSERT INTO broadcast (poster_id, content) VALUES (?, ?)`, user.GetID(db), content); err != nil {
		return err
	}

	cache.Del(context.Background(), ":broadcast")

	return nil
}

func getBroadcastList(c *gin.Context) ([]Info, error) {
	db := utils.GetDBFromContext(c)

	var broadcastList []Info
	rows, err := globals.QueryDb(db, `
		SELECT broadcast.id, broadcast.content, auth.username, broadcast.created_at
		FROM broadcast
		INNER JOIN auth ON broadcast.poster_id = auth.id
		ORDER BY broadcast.id DESC
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var broadcast Info
		var createdAt []uint8
		if err := rows.Scan(&broadcast.Index, &broadcast.Content, &broadcast.Poster, &createdAt); err != nil {
			return nil, err
		}
		broadcast.CreatedAt = utils.ConvertTime(createdAt).Format("2006-01-02 15:04:05")
		broadcastList = append(broadcastList, broadcast)
	}

	return broadcastList, nil
}

func removeBroadcast(c *gin.Context, index int) error {
	db := utils.GetDBFromContext(c)
	cache := utils.GetCacheFromContext(c)

	if _, err := globals.ExecDb(db, `DELETE FROM broadcast WHERE id = ?`, index); err != nil {
		return err
	}

	cache.Del(context.Background(), ":broadcast")

	return nil
}

func updateBroadcast(c *gin.Context, index int, content string) error {
	db := utils.GetDBFromContext(c)
	cache := utils.GetCacheFromContext(c)

	if _, err := globals.ExecDb(db, `UPDATE broadcast SET content = ? WHERE id = ?`, content, index); err != nil {
		return err
	}

	cache.Del(context.Background(), ":broadcast")

	return nil
}
