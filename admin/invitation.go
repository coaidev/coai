package admin

import (
	"chat/globals"
	"chat/utils"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
)

func GetInvitationPagination(db *sql.DB, page int64) PaginationForm {
	var invitations []interface{}
	var total int64
	if err := globals.QueryRowDb(db, `
		SELECT COUNT(*) FROM invitation
	`).Scan(&total); err != nil {
		return PaginationForm{
			Status:  false,
			Message: err.Error(),
		}
	}

	// get used_user from auth table by `user_id`
	rows, err := globals.QueryDb(db, `
		SELECT invitation.code, invitation.quota, invitation.type, invitation.used, 
		       invitation.created_at, invitation.updated_at, 
		       COALESCE(auth.username, '-') as username
		FROM invitation
		LEFT JOIN auth ON auth.id = invitation.used_id
		ORDER BY invitation.id DESC LIMIT ? OFFSET ?
	`, pagination, page*pagination)
	if err != nil {
		return PaginationForm{
			Status:  false,
			Message: err.Error(),
		}
	}

	for rows.Next() {
		var invitation InvitationData
		var createdAt sql.NullTime
		var updatedAt sql.NullTime
		if err := rows.Scan(&invitation.Code, &invitation.Quota, &invitation.Type, &invitation.Used, &createdAt, &updatedAt, &invitation.Username); err != nil {
			return PaginationForm{
				Status:  false,
				Message: err.Error(),
			}
		}

		if createdAt.Valid {
			invitation.CreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
		} else {
			invitation.CreatedAt = ""
		}

		if updatedAt.Valid {
			invitation.UpdatedAt = updatedAt.Time.Format("2006-01-02 15:04:05")
		} else {
			invitation.UpdatedAt = ""
		}

		invitations = append(invitations, invitation)
	}

	return PaginationForm{
		Status: true,
		Total:  int(math.Ceil(float64(total) / float64(pagination))),
		Data:   invitations,
	}
}

func DeleteInvitationCode(db *sql.DB, code string) error {
	_, err := globals.ExecDb(db, `
		DELETE FROM invitation WHERE code = ?
	`, code)
	return err
}

func NewInvitationCode(db *sql.DB, code string, quota float32, t string) error {
	_, err := globals.ExecDb(db, `
		INSERT INTO invitation (code, quota, type)
		VALUES (?, ?, ?)
	`, code, quota, t)
	return err
}

func GenerateInvitations(db *sql.DB, num int, quota float32, t string) InvitationGenerateResponse {
	arr := make([]string, 0)
	idx := 0
	retry := 0
	for idx < num {
		code := fmt.Sprintf("%s-%s", t, utils.GenerateChar(24))
		if err := NewInvitationCode(db, code, quota, t); err != nil {
			// ignore duplicate code
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}

			if retry < 100 && strings.Contains(err.Error(), "Duplicate entry") {
				retry++
				continue
			}

			retry = 0
			return InvitationGenerateResponse{
				Status:  false,
				Message: err.Error(),
			}
		}
		arr = append(arr, code)
		idx++
	}

	return InvitationGenerateResponse{
		Status: true,
		Data:   arr,
	}
}
