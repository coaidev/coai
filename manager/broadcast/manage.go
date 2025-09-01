package broadcast

import (
	"chat/auth"
	"chat/globals"
	"chat/utils"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
		broadcastList, err = handleRow(rows, broadcastList)

		if err != nil {
			return nil, err
		}
	}

	return broadcastList, nil
}

func handleRow(rows *sql.Rows, broadcastList []Info) ([]Info, error) {
	if viper.GetString("database.driver") == "postgres" {
		return handleRowPostgres(rows, broadcastList)
	}
	return handleRowMysql(rows, broadcastList)
}

func handleRowMysql(rows *sql.Rows, broadcastList []Info) ([]Info, error) {
	var broadcast Info
	var createdAt []uint8
	if err := rows.Scan(&broadcast.Index, &broadcast.Content, &broadcast.Poster, &createdAt); err != nil {
		return nil, err
	}
	broadcast.CreatedAt = utils.ConvertTime(createdAt).Format("2006-01-02 15:04:05")
	broadcastList = append(broadcastList, broadcast)

	return broadcastList, nil
}

func handleRowPostgres(rows *sql.Rows, broadcastList []Info) ([]Info, error) {
	var broadcast Info
	var createdAt sql.NullTime // PostgreSQL specific

	if err := rows.Scan(&broadcast.Index, &broadcast.Content, &broadcast.Poster, &createdAt); err != nil {
		return nil, err
	}
	if createdAt.Valid {
		broadcast.CreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
	} else {
		broadcast.CreatedAt = ""
	}
	broadcastList = append(broadcastList, broadcast)

	return broadcastList, nil
}
