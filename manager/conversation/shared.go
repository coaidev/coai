package conversation

import (
	"chat/auth"
	"chat/globals"
	"chat/utils"
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type SharedPreviewForm struct {
	Name           string    `json:"name"`
	ConversationId int64     `json:"conversation_id"`
	Time           time.Time `json:"time"`
	Hash           string    `json:"hash"`
}

type SharedForm struct {
	Username string            `json:"username"`
	Name     string            `json:"name"`
	Messages []globals.Message `json:"messages"`
	Model    string            `json:"model"`
	Time     time.Time         `json:"time"`
}

type SharedHashForm struct {
	Id             int64 `json:"id"`
	ConversationId int64 `json:"conversation_id"`
	Refs           []int `json:"refs"`
}

func GetRef(refs []int) (result string) {
	for _, v := range refs {
		result += strconv.Itoa(v) + ","
	}
	return strings.TrimSuffix(result, ",")
}

func ShareConversation(db *sql.DB, user *auth.User, id int64, refs []int) (string, error) {
	if id < 0 || user == nil {
		return "", nil
	}

	ref := GetRef(refs)
	hash := utils.Md5EncryptForm(SharedHashForm{
		Id:             user.GetID(db),
		ConversationId: id,
		Refs:           refs,
	})

	if _, err := globals.ExecDb(db, `
		INSERT INTO sharing (hash, user_id, conversation_id, refs) VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE refs = ?
	`, hash, user.GetID(db), id, ref, ref); err != nil {
		return "", err
	}

	return hash, nil
}

func GetSharedMessages(db *sql.DB, userId int64, conversationId int64, refs []string) []globals.Message {
	conversation := LoadConversation(db, userId, conversationId)
	if conversation == nil {
		return nil
	}

	messages := make([]globals.Message, 0)
	for _, v := range refs {
		if v == "-1" {
			return conversation.GetMessage()
		} else {
			id, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			messages = append(messages, conversation.GetMessageById(id))
		}
	}
	return messages
}

func ListSharedConversation(db *sql.DB, user *auth.User) []SharedPreviewForm {
	if user == nil {
		return nil
	}

	id := user.GetID(db)
	rows, err := globals.QueryDb(db, `
		SELECT conversation.conversation_name, conversation.conversation_id, sharing.updated_at, sharing.hash
		FROM sharing
		INNER JOIN conversation 
		    ON conversation.conversation_id = sharing.conversation_id 
		    AND conversation.user_id = sharing.user_id
		WHERE sharing.user_id = ?
		ORDER BY sharing.updated_at DESC
		LIMIT 100
	`, id)
	if err != nil {
		return nil
	}

	result := make([]SharedPreviewForm, 0)
	for rows.Next() {
		var updated sql.NullTime
		var form SharedPreviewForm
		if err := rows.Scan(&form.Name, &form.ConversationId, &updated, &form.Hash); err != nil {
			continue
		}

		if updated.Valid {
			form.Time = updated.Time
		} else {
			form.Time = time.Unix(0, 0)
		}

		result = append(result, form)
	}
	return result
}

func DeleteSharedConversation(db *sql.DB, user *auth.User, hash string) error {
	if user == nil {
		return nil
	}

	id := user.GetID(db)
	if _, err := globals.ExecDb(db, `
		DELETE FROM sharing WHERE user_id = ? AND hash = ?
	`, id, hash); err != nil {
		return err
	}
	return nil
}

func GetSharedConversation(db *sql.DB, hash string) (*SharedForm, error) {
	var shared SharedForm
	var (
		uid     int64
		cid     int64
		ref     string
		updated sql.NullTime
	)
	if err := globals.QueryRowDb(db, `
		SELECT auth.username, sharing.refs, sharing.updated_at, conversation.conversation_name,
		       sharing.user_id, sharing.conversation_id, conversation.model
		FROM sharing
		INNER JOIN auth ON auth.id = sharing.user_id
		INNER JOIN conversation ON conversation.conversation_id = sharing.conversation_id AND conversation.user_id = sharing.user_id
		WHERE sharing.hash = ?
	`, hash).Scan(&shared.Username, &ref, &updated, &shared.Name, &uid, &cid, &shared.Model); err != nil {
		return nil, err
	}

	if updated.Valid {
		shared.Time = updated.Time
	} else {
		shared.Time = time.Unix(0, 0)
	}

	refs := strings.Split(ref, ",")
	shared.Messages = GetSharedMessages(db, uid, cid, refs)

	return &shared, nil
}

func UseSharedConversation(db *sql.DB, user *auth.User, hash string) *Conversation {
	shared, err := GetSharedConversation(db, hash)
	if err != nil {
		return nil
	}

	if user == nil {
		// anonymous
		return &Conversation{
			Auth:    false,
			UserID:  -1,
			Id:      -1,
			Name:    shared.Name,
			Message: shared.Messages,
			Model:   globals.GPT3Turbo,
		}
	}

	// create new conversation
	id := user.GetID(db)
	return &Conversation{
		Auth:    true,
		Id:      GetConversationLengthByUserID(db, id) + 1,
		UserID:  id,
		Name:    shared.Name,
		Model:   globals.GPT3Turbo,
		Message: shared.Messages,
	}
}

func (c *Conversation) LoadSharing(db *sql.DB, hash string) {
	if strings.TrimSpace(hash) == "" || c.Shared == true {
		return
	}

	shared, err := GetSharedConversation(db, hash)
	if err != nil {
		return
	}

	c.InsertMessages(shared.Messages, 0)
	c.SetName(db, shared.Name)
	c.Shared = true
}
