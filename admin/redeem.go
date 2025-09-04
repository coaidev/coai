package admin

import (
	"chat/globals"
	"chat/utils"
	"database/sql"
	"fmt"
	"math"
	"strings"
)

func GetRedeemData(db *sql.DB, page int64) PaginationForm {
	var data []interface{}
	var total int64
	if err := globals.QueryRowDb(db, `
		SELECT COUNT(*) FROM redeem
	`).Scan(&total); err != nil {
		return PaginationForm{
			Status:  false,
			Message: err.Error(),
		}
	}

	rows, err := globals.QueryDb(db, `
		SELECT code, quota, used, created_at, updated_at
		FROM redeem
		ORDER BY id DESC LIMIT ? OFFSET ?
	`, pagination, page*pagination)

	if err != nil {
		return PaginationForm{
			Status:  false,
			Message: err.Error(),
		}
	}

	for rows.Next() {
		var redeem RedeemData
		var createdAt sql.NullTime
		var updatedAt sql.NullTime
		if err := rows.Scan(&redeem.Code, &redeem.Quota, &redeem.Used, &createdAt, &updatedAt); err != nil {
			return PaginationForm{
				Status:  false,
				Message: err.Error(),
			}
		}

		if createdAt.Valid {
			redeem.CreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
		} else {
			redeem.CreatedAt = ""
		}

		if updatedAt.Valid {
			redeem.UpdatedAt = updatedAt.Time.Format("2006-01-02 15:04:05")
		} else {
			redeem.UpdatedAt = ""
		}

		data = append(data, redeem)
	}

	return PaginationForm{
		Status: true,
		Total:  int(math.Ceil(float64(total) / float64(pagination))),
		Data:   data,
	}
}

func DeleteRedeemCode(db *sql.DB, code string) error {
	_, err := globals.ExecDb(db, `
		DELETE FROM redeem WHERE code = ?
	`, code)

	return err
}

func GenerateRedeemCodes(db *sql.DB, num int, quota float32) RedeemGenerateResponse {
	arr := make([]string, 0)
	idx := 0
	for idx < num {
		code, err := CreateRedeemCode(db, quota)

		if err != nil {
			return RedeemGenerateResponse{
				Status:  false,
				Message: err.Error(),
			}
		}
		arr = append(arr, code)
		idx++
	}

	return RedeemGenerateResponse{
		Status: true,
		Data:   arr,
	}
}

func CreateRedeemCode(db *sql.DB, quota float32) (string, error) {
	code := fmt.Sprintf("nio-%s", utils.GenerateChar(32))
	_, err := globals.ExecDb(db, `
		INSERT INTO redeem (code, quota) VALUES (?, ?)
	`, code, quota)

	if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
		// code name is duplicate
		return CreateRedeemCode(db, quota)
	}

	return code, err
}
